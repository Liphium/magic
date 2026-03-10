package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Liphium/magic/v3/mconfig"
	mservices "github.com/Liphium/magic/v3/mrunner/services"
	"github.com/moby/moby/client"
)

// Should create a new container for the database or use the existing one (returns container id + error in case one happened)
func (pd *PostgresDriver) CreateContainer(ctx context.Context, c *client.Client, a mconfig.ContainerAllocation) (string, error) {
	if pd.Image == "" {
		return "", fmt.Errorf("please specify a proper image")
	}

	return mservices.CreateContainer(ctx, pgLog, c, a, mservices.ManagedContainerOptions{
		Image: pd.Image,
		Env: []string{
			fmt.Sprintf("POSTGRES_PASSWORD=%s", PostgresPassword),
			fmt.Sprintf("POSTGRES_USER=%s", PostgresUsername),
			"POSTGRES_DB=postgres",
		},
		Ports: []string{
			"5432/tcp",
		},
		Volumes: []mservices.ContainerVolume{
			{NameSuffix: "data", Target: "/var/lib/postgresql"},
		},
	})
}

// Check for postgres health
func (pd *PostgresDriver) IsHealthy(ctx context.Context, c *client.Client, container mconfig.ContainerInformation) (bool, error) {
	readyCommand := "pg_isready -d postgres -U postgres -t 0"
	cmd := strings.Split(readyCommand, " ")

	// Try to execute the command
	respInspect, err := mservices.ExecuteCommand(ctx, c, container.ID, cmd)
	if err != nil {
		return false, fmt.Errorf("couldn't execute command for readiness of container: %s", err)
	}

	if mconfig.VerboseLogging {
		pgLog.Println("Database health check response code:", respInspect.ExitCode)
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

	for _, db := range pd.Databases {
		pgLog.Println("Creating database", db+"...")
		_, err := conn.Exec(fmt.Sprintf("CREATE DATABASE %s", db))
		if err != nil && !strings.Contains(err.Error(), "already exists") {
			return fmt.Errorf("couldn't create postgres database: %s", err)
		}
	}

	return nil
}
