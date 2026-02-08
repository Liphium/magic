package databases

import (
	"context"
	"database/sql"
	"fmt"
	"net/netip"
	"strings"

	"github.com/Liphium/magic/v2/mconfig"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
)

// Should create a new container for the database or use the existing one (returns container id + error in case one happened)
func (pd *PostgresDriver) CreateContainer(ctx context.Context, c *client.Client, a mconfig.ContainerAllocation) (string, error) {

	// Set to default username and password when not set
	if pd.image == "" {
		return "", fmt.Errorf("please specify a proper image")
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
			if strings.HasSuffix(n, a.Name) {
				pgLegacyLog.Println("Found existing container...")
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
		pgLegacyLog.Println("Deleting old container...")
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

	// Create the container
	resp, err := c.ContainerCreate(ctx, client.ContainerCreateOptions{
		Config: &container.Config{
			Image: pd.image,
			Env: []string{
				fmt.Sprintf("POSTGRES_PASSWORD=%s", PostgresPassword),
				fmt.Sprintf("POSTGRES_USER=%s", PostgresUsername),
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

	return resp.ID, nil
}

// Check for postgres health
func (pd *PostgresDriver) IsHealthy(ctx context.Context, c *client.Client, container mconfig.ContainerInformation) (bool, error) {
	readyCommand := "pg_isready -d postgres -U postgres -t 0"
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

	if mconfig.VerboseLogging {
		pgLegacyLog.Println("Database health check response code:", respInspect.ExitCode)
	}

	return respInspect.ExitCode == 0, nil
}

// Initialize the internal container with a script (for example)
func (pd *PostgresDriver) Initialize(ctx context.Context, c *client.Client, container mconfig.ContainerInformation) error {
	connStr := fmt.Sprintf("host=127.0.0.1 port=%d user=postgres password=postgres dbname=postgres sslmode=disable", container.Ports[0])

	// Connect to the database
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("couldn't connect to postgres: %s", err)
	}
	defer conn.Close()

	for _, db := range pd.databases {
		pgLegacyLog.Println("Creating database", db+"...")
		_, err := conn.Exec(fmt.Sprintf("CREATE DATABASE %s", db))
		if err != nil && !strings.Contains(err.Error(), "already exists") {
			return fmt.Errorf("couldn't create postgres database: %s", err)
		}
	}

	return nil
}
