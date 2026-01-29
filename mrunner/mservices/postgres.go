package mservices

import (
	"context"
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

func NewPostgresDriver(image string) *PostgresDriver {
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
	return "postgres"
}

// Should create a new container for the database or use the existing one (returns container id + error in case one happened)
func (pd *PostgresDriver) StartContainer(ctx context.Context, c *client.Client, a ContainerAllocation) (string, error) {

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
			port: []network.PortBinding{{HostIP: netip.MustParseAddr("127.0.0.1"), HostPort: fmt.Sprintf("%d", a.Port)}},
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
func (pd *PostgresDriver) IsHealthy(ctx context.Context, c *client.Client, id string) (bool, error) {
	readyCommand := "pg_isready -d postgres"
	cmd := strings.Split(readyCommand, " ")
	execConfig := client.ExecCreateOptions{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}

	// Try to execute the command
	execIDResp, err := c.ExecCreate(ctx, id, execConfig)
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
func (pd *PostgresDriver) Initialize(ctx context.Context, c *client.Client, id string) error {
	// TODO: Create
	return nil
}
