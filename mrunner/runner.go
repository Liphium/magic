package mrunner

import (
	"github.com/Liphium/magic/mconfig"
	"github.com/docker/docker/client"
)

const DefaultStartPort uint = 10000
const DefaultEndPort uint = 60000

type Runner struct {
	config      string
	profile     string
	client      *client.Client
	environment *mconfig.Environment
	databases   []*Database
}

// Create a new runner
func NewRunner(ctx *mconfig.Context) (*Runner, error) {

	// Create a new client for the docker sdk
	dc, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	// Import databases
	databases := []*Database{}
	for _, db := range ctx.Databases() {
		databases = append(databases, newDB(db))
	}

	// Create the runner
	return &Runner{
		config:      ctx.Config(),
		profile:     ctx.Profile(),
		client:      dc,
		environment: ctx.Environment(),
		databases:   databases,
	}, nil
}
