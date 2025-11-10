package mrunner

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Liphium/magic/v2/mconfig"
	"github.com/Liphium/magic/v2/util"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
)

// Deploy all the containers nessecary for the application
func (r *Runner) Deploy(deleteContainers bool) error {

	// Make sure the Docker connection is working
	_, err := r.client.Info(context.Background())
	if client.IsErrConnectionFailed(err) || client.IsErrNotFound(err) {
		return fmt.Errorf("please make sure Docker is running, and that Magic (or the Go toolchain) has access to it. (%s)", err)
	}

	// Clear all state in case wanted
	if deleteContainers {
		util.Log.Println("Clearing all state...")
		r.Clear()
	}

	// Deploy the database containers
	for _, dbType := range r.plan.DatabaseTypes {
		ctx := context.Background()
		name := dbType.ContainerName(r.appName, r.profile)
		util.Log.Println("Creating database container", name+"...")

		// Check if the container already exists
		f := filters.NewArgs()
		f.Add("name", name)
		summary, err := r.client.ContainerList(ctx, container.ListOptions{
			Filters: f,
			All:     true,
		})
		if err != nil {
			return fmt.Errorf("couldn't list containers: %s", err)
		}
		containerId := ""
		var containerMounts []mount.Mount = nil
		for _, c := range summary {
			for _, n := range c.Names {
				if strings.Contains(n, name) {
					util.Log.Println("Found existing container...")
					containerId = c.ID

					// Inspect the cotainer to get the mounts
					resp, err := r.client.ContainerInspect(ctx, c.ID)
					if err != nil {
						return fmt.Errorf("couldn't inspect container: %s", err)
					}
					containerMounts = resp.HostConfig.Mounts
				}
			}
		}

		// Delete the container if it exists
		if containerId != "" {
			if err := r.client.ContainerRemove(ctx, containerId, container.RemoveOptions{
				RemoveVolumes: false,
				Force:         true,
			}); err != nil {
				return fmt.Errorf("couldn't delete database container: %s", err)
			}
		}

		// Create the new container with the volumes
		util.Log.Println("Creating new container...")
		containerId, err = r.createDatabaseContainer(ctx, dbType, name, containerMounts)
		if err != nil {
			return fmt.Errorf("couldn't create database container: %s", err)
		}

		// Start the container
		util.Log.Println("Trying to start container...")
		if err := r.client.ContainerStart(ctx, containerId, container.StartOptions{}); err != nil {
			return fmt.Errorf("couldn't start postgres container: %s", err)
		}

		// Wait for the container to start (with pg_isready)
		util.Log.Println("Waiting for PostgreSQL to be ready...")
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
				return fmt.Errorf("couldn't create command for readiness of container: %s", err)
			}
			execStartCheck := container.ExecStartOptions{Detach: false, Tty: false}
			if err := r.client.ContainerExecStart(ctx, execIDResp.ID, execStartCheck); err != nil {
				return fmt.Errorf("couldn't start command for readiness of container: %s", err)
			}
			respInspect, err := r.client.ContainerExecInspect(ctx, execIDResp.ID)
			if err != nil {
				return fmt.Errorf("couldn't inspect command for readiness of container: %s", err)
			}
			if respInspect.ExitCode == 0 {
				break
			}

			time.Sleep(200 * time.Millisecond)
		}
		time.Sleep(200 * time.Millisecond) // Some additional time, sometimes takes longer

		// Create all of the databases
		util.Log.Println("Connecting to PostgreSQL...")
		if err := r.createDatabases(dbType); err != nil {
			return err
		}
	}

	// Load environment variables into current application
	util.Log.Println("Loading environment...")
	for key, value := range r.plan.Environment {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("couldn't set environment variable %s: %s", key, err)
		}
	}

	util.Log.Println("Deployment finished.")
	return nil
}

// Create a new container for a postgres database. Returns container id.
func (r *Runner) createDatabaseContainer(ctx context.Context, dbType mconfig.PlannedDatabaseType, name string, mounts []mount.Mount) (string, error) {

	// Reserve the port for the container
	port, err := nat.NewPort("tcp", "5432")
	if err != nil {
		return "", fmt.Errorf("couldn't create port for postgres container: %s", err)
	}
	exposedPorts := nat.PortSet{port: struct{}{}}

	// If no existing mounts, create a new volume for PostgreSQL data
	if mounts == nil {
		volumeName := fmt.Sprintf("%s-postgres-data", name)
		mounts = []mount.Mount{
			{
				Type:   mount.TypeVolume,
				Source: volumeName,
				Target: "/var/lib/postgresql/data",
			},
		}
	}

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
		return "", fmt.Errorf("couldn't create postgres container: %s", err)
	}
	return resp.ID, nil
}

// Connect to postgres and create all the needed databases.
func (r *Runner) createDatabases(dbType mconfig.PlannedDatabaseType) error {
	connStr := fmt.Sprintf("host=127.0.0.1 port=%d user=postgres password=postgres dbname=postgres sslmode=disable", dbType.Port)

	// Connect to the database
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("couldn't connect to postgres: %s", err)
	}
	defer conn.Close()

	for _, db := range dbType.Databases {
		util.Log.Println("Creating database", db.Name+"...")
		_, err := conn.Exec(fmt.Sprintf("CREATE DATABASE %s", db.Name))
		if err != nil && !strings.Contains(err.Error(), "already exists") {
			return fmt.Errorf("couldn't create postgres database: %s", err)
		}
	}
	return nil
}

// Delete all containers and reset all state
func (r *Runner) Clear() {
	ctx := context.Background()
	for _, dbType := range r.plan.DatabaseTypes {

		// Try to find the container for the type
		f := filters.NewArgs()
		name := dbType.ContainerName(r.appName, r.profile)
		f.Add("name", name)
		summary, err := r.client.ContainerList(ctx, container.ListOptions{
			Filters: f,
			All:     true,
		})
		if err != nil {
			log.Fatalln("Couldn't list containers:", err)
		}
		containerId := ""
		for _, c := range summary {
			for _, n := range c.Names {
				if strings.Contains(n, name) {
					containerId = c.ID
				}
			}
		}

		// If there is no container, nothing to delete
		if containerId == "" {
			continue
		}

		// Get all the attached volumes to delete them manually
		containerInfo, err := r.client.ContainerInspect(ctx, containerId)
		if err != nil {
			util.Log.Println("Warning: Couldn't inspect container:", err)
		}
		var volumeNames []string
		if containerInfo.Mounts != nil {
			for _, mnt := range containerInfo.Mounts {
				if mnt.Type == mount.TypeVolume && mnt.Name != "" {
					volumeNames = append(volumeNames, mnt.Name)
				}
			}
		}

		// Delete the container
		util.Log.Println("Deleting container", name+"...")
		if err := r.client.ContainerRemove(ctx, containerId, container.RemoveOptions{
			RemoveVolumes: false,
			Force:         true,
		}); err != nil {
			util.Log.Fatalln("Couldn't delete database container:", err)
		}

		// Delete all the attached volumes
		for _, volumeName := range volumeNames {
			util.Log.Println("Deleting volume", volumeName+"...")
			if err := r.client.VolumeRemove(ctx, volumeName, true); err != nil {
				util.Log.Println("Warning: Couldn't delete volume", volumeName+":", err)
			}
		}
	}
}

// Stop all containers
func (r *Runner) StopContainers() {
	ctx := context.Background()
	for _, dbType := range r.plan.DatabaseTypes {

		// Try to find the container for the type
		f := filters.NewArgs()
		name := dbType.ContainerName(r.appName, r.profile)
		f.Add("name", name)
		summary, err := r.client.ContainerList(ctx, container.ListOptions{
			Filters: f,
		})
		if err != nil {
			util.Log.Fatalln("Couldn't list containers:", err)
		}
		containerId := ""
		for _, c := range summary {
			for _, n := range c.Names {
				if strings.Contains(n, name) {
					containerId = c.ID
				}
			}
		}

		// If there is no container, nothing to stop
		if containerId == "" {
			continue
		}

		// Stop the container
		util.Log.Println("Stopping container", name+"...")
		if err := r.client.ContainerStop(ctx, containerId, container.StopOptions{}); err != nil {
			util.Log.Fatalln("Couldn't stop database container:", err)
		}
	}
}

// Clear the databases, at runtime
func (r *Runner) ClearDatabases() {

	// Delete all the databases of every type
	for _, dbType := range r.plan.DatabaseTypes {

		for _, db := range dbType.Databases {
			connStr := fmt.Sprintf("host=127.0.0.1 port=%d user=postgres password=postgres dbname=%s sslmode=disable", dbType.Port, db.Name)

			// Connect to the database
			conn, err := sql.Open("postgres", connStr)
			if err != nil {
				log.Fatalln("couldn't connect to postgres:", err)
			}
			defer conn.Close()

			// Clear all of the tables
			res, err := conn.Query("SELECT table_name FROM information_schema.tables WHERE table_schema NOT IN ('pg_catalog', 'information_schema')")
			if err != nil {
				log.Fatalln("couldn't get database tables:", err)
			}
			for res.Next() {
				var name string
				if err := res.Scan(&name); err != nil {
					util.Log.Fatalln("couldn't get database table name:", err)
				}
				if _, err := conn.Exec(fmt.Sprintf("truncate %s CASCADE", name)); err != nil {
					util.Log.Fatalln("couldn't delete from table", name+":", err)
				}
			}
		}
	}
}
