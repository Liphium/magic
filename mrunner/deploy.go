package mrunner

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Liphium/magic/mconfig"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

// Deploy all the containers nessecary for the application
func (r *Runner) Deploy() {

	// Deploy the database containers
	for _, dbType := range r.plan.DatabaseTypes {
		ctx := context.Background()
		name := fmt.Sprintf("mgc-%s-%s-%d", r.config, r.profile, dbType.Type)

		// Reserve the port for the container
		port, err := nat.NewPort("tcp", "5432")
		if err != nil {
			log.Fatalln("couldn't create port for postgres container:", err)
		}
		exposedPorts := nat.PortSet{port: struct{}{}}

		// Create the network config for the container
		networkConf := &container.HostConfig{
			NetworkMode: "host",
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
			},
			ExposedPorts: exposedPorts,
		}, networkConf, &network.NetworkingConfig{}, nil, name)
		if err != nil {
			log.Fatalln("couldn't create postgres container:", err)
		}

		// Start the container
		if err := r.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
			log.Fatalln("couldn't start postgres container:", err)
		}
	}

	fmt.Println("Postgres container started. Sleeping 5s before connection...")
	time.Sleep(5 * time.Second) // wait for db to be ready

	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=testdb sslmode=disable"
	fmt.Println("Connect string:", connStr)

}
