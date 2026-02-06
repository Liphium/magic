package mservices

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/netip"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
)

// Make sure the driver complies
var _ ServiceDriver = &PostgresDriver{}

const (
	DefaultPostgresUsername = "postgres"
	DefaultPostgresPassword = "postgres"
)

var pgLog *log.Logger = log.New(os.Stdout, "pg-manager ", log.Default().Flags())

type PostgresDriver struct {
	image string

	// Database credentials
	username string
	password string

	databases []string
}

// Create a new PostgreSQL service driver.
//
// It currently supports version 14-17.
//
// This driver will eventually be renamed into the legacy driver for people who still want to use PostgreSQL v17 or lower. Eventually it will be deprecated and fully removed (only once the v18 driver is available).
func NewPostgresDriver(image string) *PostgresDriver {
	imageVersion := strings.Split(image, ":")[1]

	// Supported (confirmed and tested) major versions for this Postgres driver
	var supportedPostgresVersions = []string{"14", "15", "16", "17"}

	// Do a quick check to make sure the image version is actually supported
	supported := false
	for _, version := range supportedPostgresVersions {
		if strings.HasPrefix(imageVersion, fmt.Sprintf("%s.", version)) {
			supported = true
		}
	}
	if !supported {
		pgLog.Fatalln("ERROR: Version", imageVersion, "is currently not supported.")
	}

	return &PostgresDriver{
		image: image,
	}
}

func (pd *PostgresDriver) WithUsername(name string) *PostgresDriver {
	pd.username = name
	return pd
}

func (pd *PostgresDriver) WithPassword(password string) *PostgresDriver {
	pd.password = password
	return pd
}

func (pd *PostgresDriver) NewDatabase(name string) *PostgresDriver {
	pd.databases = append(pd.databases, name)
	return pd
}

// A unique identifier for the database container
func (pd *PostgresDriver) GetUniqueId() string {
	return "postgres1417"
}

func (pd *PostgresDriver) GetRequiredPortAmount() int {
	return 1
}

func (pd *PostgresDriver) GetImage() string {
	return pd.image
}

// Should create a new container for the database or use the existing one (returns container id + error in case one happened)
func (pd *PostgresDriver) CreateContainer(ctx context.Context, c *client.Client, a ContainerAllocation) (string, error) {

	// Set to default username and password when not set
	if pd.username == "" {
		pd.username = DefaultPostgresUsername
	}
	if pd.password == "" {
		pd.password = DefaultPostgresPassword
	}
	if pd.image == "" {
		pd.image = "postgres:latest"
	}

	// Check if the container already exists
	f := make(client.Filters)
	f.Add("name", a.Name)
	summary, err := c.ContainerList(ctx, client.ContainerListOptions{
		Filters: f,
		All:     true,
	})
	if err != nil {
		return "", fmt.Errorf("couldn't list containers: %s", err)
	}
	containerId := ""
	var mounts []mount.Mount = nil
	for _, container := range summary.Items {
		for _, n := range container.Names {
			if n == a.Name {
				pgLog.Println("Found existing container...")
				containerId = container.ID

				// Inspect the container to get the mounts
				resp, err := c.ContainerInspect(ctx, container.ID, client.ContainerInspectOptions{})
				if err != nil {
					return "", fmt.Errorf("couldn't inspect container: %s", err)
				}
				mounts = resp.Container.HostConfig.Mounts
			}
		}
	}

	// Delete the container if it exists
	if containerId != "" {
		if _, err := c.ContainerRemove(ctx, containerId, client.ContainerRemoveOptions{
			RemoveVolumes: false,
			Force:         true,
		}); err != nil {
			return "", fmt.Errorf("couldn't delete database container: %s", err)
		}
	}

	// Create the port on the postgres container (this is not the port for outside)
	port, err := network.ParsePort("5432/tcp")
	if err != nil {
		return "", fmt.Errorf("couldn't create port for postgres container: %s", err)
	}
	exposedPorts := network.PortSet{port: struct{}{}}

	// If no existing mounts, create a new volume for PostgreSQL data
	if mounts == nil {
		volumeName := fmt.Sprintf("%s-data", a.Name)
		mounts = []mount.Mount{
			{
				Type:   mount.TypeVolume,
				Source: volumeName,
				Target: "/var/lib/postgresql/data",
			},
		}
	}

	// Create the network config for the container (exposes the container to the host)
	networkConf := &container.HostConfig{
		PortBindings: network.PortMap{
			port: []network.PortBinding{{HostIP: netip.MustParseAddr("127.0.0.1"), HostPort: fmt.Sprintf("%d", a.Ports[0])}},
		},
		Mounts: mounts,
	}

	// Pull the image
	pgLog.Println("Pulling image", pd.image, "...")
	pullResponse, err := c.ImagePull(ctx, pd.image, client.ImagePullOptions{})
	if err != nil {
		return "", fmt.Errorf("couldn't pull image %s: %v", pd.image, err)
	}
	pullResponse.Wait(ctx)

	// Create the container
	resp, err := c.ContainerCreate(ctx, client.ContainerCreateOptions{
		Config: &container.Config{
			Image: pd.image,
			Env: []string{
				fmt.Sprintf("POSTGRES_PASSWORD=%s", pd.password),
				fmt.Sprintf("POSTGRES_USER=%s", pd.username),
				"POSTGRES_DATABASE=postgres",
			},
			ExposedPorts: exposedPorts,
		},
		HostConfig: networkConf,
		Name:       a.Name,
	})
	if err != nil {
		return "", fmt.Errorf("couldn't create postgres container: %s", err)
	}

	// Start the container
	pgLog.Println("Trying to start container...")
	if _, err := c.ContainerStart(ctx, containerId, client.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("couldn't start postgres container: %s", err)
	}

	pgLog.Println("Database container started.")
	return resp.ID, nil
}

// Check for postgres health
func (pd *PostgresDriver) IsHealthy(ctx context.Context, c *client.Client, container ContainerInformation) (bool, error) {
	readyCommand := "pg_isready -d postgres"
	cmd := strings.Split(readyCommand, " ")
	execConfig := client.ExecCreateOptions{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}

	// Try to execute the command
	execIDResp, err := c.ExecCreate(ctx, container.ID, execConfig)
	if err != nil {
		return false, fmt.Errorf("couldn't create command for readiness of container: %s", err)
	}
	execStartCheck := client.ExecStartOptions{Detach: false, TTY: false}
	if _, err := c.ExecStart(ctx, execIDResp.ID, execStartCheck); err != nil {
		return false, fmt.Errorf("couldn't start command for readiness of container: %s", err)
	}
	respInspect, err := c.ExecInspect(ctx, execIDResp.ID, client.ExecInspectOptions{})
	if err != nil {
		return false, fmt.Errorf("couldn't inspect command for readiness of container: %s", err)
	}

	return respInspect.ExitCode == 0, nil
}

// Initialize the internal container with a script (for example)
func (pd *PostgresDriver) Initialize(ctx context.Context, c *client.Client, container ContainerInformation) error {
	connStr := fmt.Sprintf("host=127.0.0.1 port=%d user=postgres password=postgres dbname=postgres sslmode=disable", container.Ports[0])

	// Connect to the database
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("couldn't connect to postgres: %s", err)
	}
	defer conn.Close()

	for _, db := range pd.databases {
		pgLog.Println("Creating database", db+"...")
		_, err := conn.Exec(fmt.Sprintf("CREATE DATABASE %s", db))
		if err != nil && !strings.Contains(err.Error(), "already exists") {
			return fmt.Errorf("couldn't create postgres database: %s", err)
		}
	}

	return nil
}

// Handles the instructions for PostgreSQL.
//
// Supports the following instructions currently:
// - Clear tables
// - Drop tables
func (pd *PostgresDriver) HandleInstruction(ctx context.Context, c *client.Client, container ContainerInformation, instruction Instruction) error {
	switch instruction {
	case InstructionClearTables:
		return pd.ClearTables(container)
	case InstructionDropTables:
		return pd.DropTables(container)
	}
	return nil
}

// iterateTablesFn is a function that processes each table in the database
type iterateTablesFn func(tableName string, conn *sql.DB) error

// iterateTables iterates through all tables in all databases and applies the given function
func (pd *PostgresDriver) iterateTables(container ContainerInformation, fn iterateTablesFn) error {
	// For all databases, connect and iterate tables
	for _, db := range pd.databases {
		connStr := fmt.Sprintf("host=127.0.0.1 port=%d user=postgres password=postgres dbname=%s sslmode=disable", container.Ports[0], db)

		// Connect to the database
		conn, err := sql.Open("postgres", connStr)
		if err != nil {
			return fmt.Errorf("couldn't connect to postgres: %v", err)
		}
		defer conn.Close()

		// Get all of the tables
		res, err := conn.Query("SELECT table_name FROM information_schema.tables WHERE table_schema NOT IN ('pg_catalog', 'information_schema')")
		if err != nil {
			return fmt.Errorf("couldn't get database tables: %v", err)
		}
		for res.Next() {
			var name string
			if err := res.Scan(&name); err != nil {
				return fmt.Errorf("couldn't get database table name: %v", err)
			}
			if err := fn(name, conn); err != nil {
				return err
			}
		}
	}

	return nil
}

// Clear all tables in all databases (keeps table schema alive, just removes the content of all tables)
func (pd *PostgresDriver) ClearTables(container ContainerInformation) error {
	return pd.iterateTables(container, func(tableName string, conn *sql.DB) error {
		if _, err := conn.Exec(fmt.Sprintf("truncate %s CASCADE", tableName)); err != nil {
			return fmt.Errorf("couldn't truncate table %s: %v", tableName, err)
		}
		return nil
	})
}

// Drop all tables in all databases (actually deletes all of your tables)
func (pd *PostgresDriver) DropTables(container ContainerInformation) error {
	return pd.iterateTables(container, func(tableName string, conn *sql.DB) error {
		if _, err := conn.Exec(fmt.Sprintf("DROP TABLE %s CASCADE", tableName)); err != nil {
			return fmt.Errorf("couldn't drop table table %s: %v", tableName, err)
		}
		return nil
	})
}
