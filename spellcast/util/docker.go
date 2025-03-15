package util

import (
	"log"

	"github.com/docker/docker/client"
)

var DockerClient *client.Client

// Initialize the docker API client
func InitDocker() {
	var err error
	DockerClient, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalln("couldn't init docker:", err)
	}
}
