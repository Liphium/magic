package mrunner

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Liphium/magic/mconfig"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
)

// Deploy all the containers nessecary for the application
func (r *Runner) Deploy() {

	// Deploy the database containers
	for _, dbType := range r.plan.DatabaseTypes {
		ctx := context.Background()
		name := fmt.Sprintf("mgc-%s-%s-%d", r.config, r.profile, dbType.Type)

		// Check if the container already exists
		f := filters.NewArgs()
		f.Add("name", name)
		summary, err := r.client.ContainerList(ctx, container.ListOptions{
			Filters: f,
		})
		if err != nil {
			log.Fatalln("couldn't list containers:", err)
		}
		containerId := ""
		for _, c := range summary {
			for _, n := range c.Names {
				if n == name {
					containerId = c.ID
				}
			}
		}

		// Create container if it doesn't exist
		if containerId == "" {
			containerId = r.createDatabaseContainer(ctx, dbType, name)
		}

		// Start the container
		if err := r.client.ContainerStart(ctx, containerId, container.StartOptions{}); err != nil {
			log.Fatalln("couldn't start postgres container:", err)
		}

		fmt.Println("hello wrold")

		// Wait for the container to start (with pg_isready)
		execConfig := container.ExecOptions{
			Cmd:          []string{"pg_isready"},
			AttachStdout: true,
			AttachStderr: true,
		}
		for {
			execIDResp, err := r.client.ContainerExecCreate(ctx, containerId, execConfig)
			if err != nil {
				log.Fatalln("couldn't create command for readiness of container:", err)
			}
			respInspect, err := r.client.ContainerExecInspect(ctx, execIDResp.ID)
			if err != nil {
				log.Fatalln("couldn't inspect command for readiness of container:", err)
			}
			if err := r.client.ContainerExecStart(ctx, execIDResp.ID, container.ExecStartOptions{}); err != nil {
				log.Fatalln("couldn't start command for readiness of container:", err)
			}
			if respInspect.ExitCode == 0 {
				break
			}
			time.Sleep(200 * time.Millisecond)
		}

		fmt.Println("creating..")

		// Create all of the databases
		r.createDatabases(dbType)
	}
}

// Create a new container for a postgres database. Returns container id.
func (r *Runner) createDatabaseContainer(ctx context.Context, dbType mconfig.PlannedDatabaseType, name string) string {

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
	fmt.Println("Connect string:", connStr)

	// Connect to the database
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for _, db := range dbType.Databases {
		_, err := conn.Exec(fmt.Sprintf("CREATE DATABASE %s", db.Name))
		if err != nil {
			log.Fatalln("couldn't create postgres database:", err)
		}
	}
}
