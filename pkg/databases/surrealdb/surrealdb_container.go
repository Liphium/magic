package surrealdb

import (
	"context"
	"fmt"

	"github.com/Liphium/magic/v4/mconfig"
	mservices "github.com/Liphium/magic/v4/mrunner/services"
	"github.com/moby/moby/client"
	sdb "github.com/surrealdb/surrealdb.go"
)

// Should create a new container for the database or use the existing one (returns container id + error in case one happened)
func (sd *SurrealDriver) CreateContainer(ctx context.Context, c *client.Client, a mconfig.ContainerAllocation) (string, error) {
	if sd.Image == "" {
		return "", fmt.Errorf("please specify a proper image")
	}

	return mservices.CreateContainer(ctx, surrealLog, c, a, mservices.ManagedContainerOptions{
		Image: sd.Image,
		Ports: RequiredPorts,
		Volumes: []mservices.ContainerVolume{
			// The SurrealDB container image runs as a non-root user, so the volume has to be mounted at the user's home directory for it to be writable.
			{NameSuffix: "data", Target: "/home/nonroot"},
		},
		Cmd: []string{
			"start",
			"--user", SurrealUsername,
			"--pass", SurrealPassword,
			"rocksdb:/home/nonroot/data/database.db",
		},
	})
}

// Check for SurrealDB health using the built-in is-ready command
func (sd *SurrealDriver) IsHealthy(ctx context.Context, c *client.Client, container mconfig.ContainerInformation) (bool, error) {
	cmd := []string{"/surreal", "is-ready"}

	// Try to execute the command
	respInspect, err := mservices.ExecuteCommand(ctx, c, container.ID, cmd)
	if err != nil {
		return false, fmt.Errorf("couldn't execute command for readiness of container: %s", err)
	}

	if mconfig.VerboseLogging {
		surrealLog.Println("Database health check response code:", respInspect.ExitCode)
	}

	return respInspect.ExitCode == 0, nil
}

// Initialize the container by creating all of the namespaces and databases the driver should manage
func (sd *SurrealDriver) Initialize(ctx context.Context, c *client.Client, container mconfig.ContainerInformation) error {
	// Connect to the database
	conn, err := connect(ctx, container)
	if err != nil {
		return fmt.Errorf("couldn't connect to surrealdb: %s", err)
	}
	defer conn.Close(ctx)

	for _, database := range sd.Databases {
		surrealLog.Println("Creating namespace", database.Namespace, "and database", database.Database+"...")

		// Define the namespace, switch over to it and then define the database itself
		query := fmt.Sprintf("DEFINE NAMESPACE IF NOT EXISTS `%s`; USE NS `%s`; DEFINE DATABASE IF NOT EXISTS `%s`;",
			database.Namespace, database.Namespace, database.Database)
		if _, err := sdb.Query[map[string]any](ctx, conn, query, nil); err != nil {
			return fmt.Errorf("couldn't create surrealdb namespace/database %s/%s: %s", database.Namespace, database.Database, err)
		}

		// Select the namespace and database so instructions know where to operate
		if err := conn.Use(ctx, database.Namespace, database.Database); err != nil {
			return fmt.Errorf("couldn't select surrealdb namespace/database %s/%s: %s", database.Namespace, database.Database, err)
		}
	}

	return nil
}

// connect opens a new connection to the SurrealDB container using the root credentials.
func connect(ctx context.Context, container mconfig.ContainerInformation) (*sdb.DB, error) {
	conn, err := sdb.FromEndpointURLString(ctx, fmt.Sprintf("ws://127.0.0.1:%d", container.Ports[0]))
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to surrealdb: %s", err)
	}

	if _, err := conn.SignIn(ctx, sdb.Auth{
		Username: SurrealUsername,
		Password: SurrealPassword,
	}); err != nil {
		return nil, fmt.Errorf("couldn't sign into surrealdb: %s", err)
	}

	return conn, nil
}
