package mdatabase

import (
	"context"
	"fmt"
	"log"
	"net/netip"
	"os"
	"strings"
	"time"

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

type PostgresDriver struct{}

// A unique identifier for the database container
func (pd *PostgresDriver) GetUniqueId() string {
	return "postgres"
}

// Should create a new container for the database or use the existing one (returns container id + error in case one happened)
func (pd *PostgresDriver) StartContainer(ctx context.Context, a DatabaseContainerAllocation, c *client.Client) (string, error) {

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

	// Check if an environment variable is set for the postgres image
	postgresImage := os.Getenv("MAGIC_POSTGRES_IMAGE")
	if postgresImage == "" {
		postgresImage = "postgres:latest"
	}

	// Create the container
	resp, err := c.ContainerCreate(ctx, client.ContainerCreateOptions{
		Config: &container.Config{
			Image: postgresImage,
			Env: []string{
				fmt.Sprintf("POSTGRES_PASSWORD=%s", DefaultPostgresPassword),
				fmt.Sprintf("POSTGRES_USER=%s", DefaultPostgresUsername),
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

	// Wait for the container to start (with pg_isready)
	pgLog.Println("Waiting for PostgreSQL to be ready...")
	readyCommand := "pg_isready -d postgres"
	cmd := strings.Split(readyCommand, " ")
	execConfig := client.ExecCreateOptions{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}
	for {
		execIDResp, err := c.ExecCreate(ctx, containerId, execConfig)
		if err != nil {
			return "", fmt.Errorf("couldn't create command for readiness of container: %s", err)
		}
		execStartCheck := client.ExecStartOptions{Detach: false, TTY: false}
		if _, err := c.ExecStart(ctx, execIDResp.ID, execStartCheck); err != nil {
			return "", fmt.Errorf("couldn't start command for readiness of container: %s", err)
		}
		respInspect, err := c.ExecInspect(ctx, execIDResp.ID, client.ExecInspectOptions{})
		if err != nil {
			return "", fmt.Errorf("couldn't inspect command for readiness of container: %s", err)
		}
		if respInspect.ExitCode == 0 {
			break
		}

		time.Sleep(200 * time.Millisecond)
	}
	time.Sleep(200 * time.Millisecond) // Some additional time, sometimes takes longer

	pgLog.Println("Database container started.")
	return resp.ID, nil
}
