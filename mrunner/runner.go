package mrunner

import (
	"github.com/docker/docker/client"
)

const DefaultStartPort uint = 10000
const DefaultEndPort uint = 60000

type Runner struct {
	client    *client.Client
	databases []*Database
}

// Create a new runner
func NewRunner() (*Runner, error) {

	// Create a new client for the docker sdk
	dc, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &Runner{
		client: dc,
	}, nil
}
