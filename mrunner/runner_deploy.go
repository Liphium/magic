package mrunner

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Liphium/magic/mconfig"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
)

// Deploy all the containers nessecary for the application
func (r *Runner) Deploy(deleteContainers bool) {

	// Clear all state in case wanted
	if deleteContainers {
		log.Println("Clearing all state...")
		r.Clear()
	}

	// Deploy the database containers
	for _, dbType := range r.plan.DatabaseTypes {
		ctx := context.Background()
		name := dbType.ContainerName(r.module, r.config, r.profile)
		log.Println("Creating database container", name+"...")

		// Check if the container already exists
		f := filters.NewArgs()
		f.Add("name", name)
		summary, err := r.client.ContainerList(ctx, container.ListOptions{
			Filters: f,
			All:     true,
		})
		if err != nil {
			log.Fatalln("couldn't list containers:", err)
		}
		containerId := ""
		var containerMounts []mount.Mount = nil
		for _, c := range summary {
			for _, n := range c.Names {
				if strings.Contains(n, name) {
					log.Println("Found existing container...")
					containerId = c.ID

					// Inspect the cotainer to get the mounts
					resp, err := r.client.ContainerInspect(ctx, c.ID)
					if err != nil {
						log.Fatalln("Couldn't inspect container:", err)
					}
					containerMounts = resp.HostConfig.Mounts
				}
			}
		}

		// Delete the container if it exists
		if containerId != "" {
			if err := r.client.ContainerRemove(ctx, containerId, container.RemoveOptions{
				Force: true,
			}); err != nil {
				log.Fatalln("Couldn't delete database container:", err)
			}
		}

		// Create the new container with the volumes
		log.Println("Creating new container...")
		containerId = r.createDatabaseContainer(ctx, dbType, name, containerMounts)

		// Start the container
		log.Println("Trying to start container...")
		if err := r.client.ContainerStart(ctx, containerId, container.StartOptions{}); err != nil {
			log.Fatalln("couldn't start postgres container:", err)
		}

		// Wait for the container to start (with pg_isready)
		log.Println("Waiting for PostgreSQL to be ready...")
		readyCommand := "pg_isready -d postgres"
		cmd := strings.Split(readyCommand, " ")
		execConfig := container.ExecOptions{
			Cmd:          cmd,
			AttachStdout: true,
			AttachStderr: true,
		}
		for {
			execIDResp, err := r.client.ContainerExecCreate(ctx, containerId, execConfig)
			if err != nil {
				log.Fatalln("couldn't create command for readiness of container:", err)
			}
			execStartCheck := container.ExecStartOptions{Detach: false, Tty: false}
			if err := r.client.ContainerExecStart(ctx, execIDResp.ID, execStartCheck); err != nil {
				log.Fatalln("couldn't start command for readiness of container:", err)
			}
			respInspect, err := r.client.ContainerExecInspect(ctx, execIDResp.ID)
			if err != nil {
				log.Fatalln("couldn't inspect command for readiness of container:", err)
			}
			if respInspect.ExitCode == 0 {
				break
			}

			time.Sleep(200 * time.Millisecond)
		}
		time.Sleep(200 * time.Millisecond) // Some additional time, sometimes takes longer

		// Create all of the databases
		log.Println("Connecting to PostgreSQL...")
		r.createDatabases(dbType)
	}

	// Load environment variables into current application
	log.Println("Loading environment...")
	for key, value := range r.plan.Environment {
		if err := os.Setenv(key, value); err != nil {
			log.Fatalln("couldn't set environment variable", key+":", err)
		}
	}

	log.Println("Deployment finished.")
	log.Println(" ")
}

// Create a new container for a postgres database. Returns container id.
func (r *Runner) createDatabaseContainer(ctx context.Context, dbType mconfig.PlannedDatabaseType, name string, mounts []mount.Mount) string {

	// Reserve the port for the container
	port, err := nat.NewPort("tcp", "5432")
	if err != nil {
		log.Fatalln("couldn't create port for postgres container:", err)
	}
	exposedPorts := nat.PortSet{port: struct{}{}}

	// Create the network config for the container
	networkConf := &container.HostConfig{
		PortBindings: nat.PortMap{
			port: []nat.PortBinding{{HostIP: "127.0.0.1", HostPort: fmt.Sprintf("%d", dbType.Port)}},
		},
		Mounts: mounts,
	}

	// Create the container
	resp, err := r.client.ContainerCreate(ctx, &container.Config{
		Image: "postgres:latest",
		Env: []string{
			fmt.Sprintf("POSTGRES_PASSWORD=%s", mconfig.DefaultPassword(dbType.Type)),
			fmt.Sprintf("POSTGRES_USERNAME=%s", mconfig.DefaultUsername(dbType.Type)),
			"POSTGRES_DATABASE=postgres",
		},
		ExposedPorts: exposedPorts,
	}, networkConf, &network.NetworkingConfig{}, nil, name)
	if err != nil {
		log.Fatalln("couldn't create postgres container:", err)
	}
	return resp.ID
}

// Connect to postgres and create all the needed databases.
func (r *Runner) createDatabases(dbType mconfig.PlannedDatabaseType) {
	connStr := fmt.Sprintf("host=127.0.0.1 port=%d user=postgres password=postgres dbname=postgres sslmode=disable", dbType.Port)

	// Connect to the database
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln("couldn't connect to postgres:", err)
	}
	defer conn.Close()

	for _, db := range dbType.Databases {
		log.Println("Creating database", db.Name+"...")
		_, err := conn.Exec(fmt.Sprintf("CREATE DATABASE %s", db.Name))
		if err != nil && !strings.Contains(err.Error(), "already exists") {
			log.Fatalln("couldn't create postgres database:", err)
		}
	}
}

// Delete all containers and reset all state
func (r *Runner) Clear() {
	ctx := context.Background()
	for _, dbType := range r.plan.DatabaseTypes {

		// Try to find the container for the type
		f := filters.NewArgs()
		name := dbType.ContainerName(r.module, r.config, r.profile)
		f.Add("name", name)
		summary, err := r.client.ContainerList(ctx, container.ListOptions{
			Filters: f,
		})
		if err != nil {
			log.Fatalln("Couldn't list containers:", err)
		}
		containerId := ""
		for _, c := range summary {
			for _, n := range c.Names {
				fmt.Println("found", n)
				if strings.Contains(n, name) {
					containerId = c.ID
				}
			}
		}

		// If there is no container, nothing to delete
		if containerId == "" {
			continue
		}

		// Delete the container
		log.Println("Deleting container", name+"...")
		if err := r.client.ContainerRemove(ctx, containerId, container.RemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		}); err != nil {
			log.Fatalln("Couldn't delete database container:", err)
		}
	}
}

// Clear the databases, at runtime
func (r *Runner) ClearDatabases() {

	// Delete all the databases of every type
	for _, dbType := range r.plan.DatabaseTypes {
		connStr := fmt.Sprintf("host=127.0.0.1 port=%d user=postgres password=postgres dbname=postgres sslmode=disable", dbType.Port)

		// Connect to the database
		conn, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatalln("couldn't connect to postgres:", err)
		}
		defer conn.Close()

		for _, db := range dbType.Databases {
			if mconfig.VerboseLogging {
				log.Println("Re-creating database", db.Name+"...")
			}

			// Drop the database
			_, err := conn.Exec(fmt.Sprintf("DROP DATABASE %s", db.Name))
			if err != nil {
				log.Fatalln("couldn't drop postgres database:", err)
			}

			// Create it again
			_, err = conn.Exec(fmt.Sprintf("CREATE DATABASE %s", db.Name))
			if err != nil {
				log.Fatalln("couldn't create postgres database:", err)
			}
		}
	}
}
