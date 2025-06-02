package mrunner

import (
	"github.com/Liphium/magic/mconfig"
	"github.com/docker/docker/client"
)

const DefaultStartPort uint = 10000
const DefaultEndPort uint = 60000

type Runner struct {
	module  string
	config  string
	profile string
	client  *client.Client
	ctx     *mconfig.Context
	plan    *mconfig.Plan
}

func (r *Runner) Environment() *mconfig.Environment {
	return r.ctx.Environment()
}

// Create a new runner
func NewRunner(ctx *mconfig.Context) (*Runner, error) {

	// Create a new client for the docker sdk
	dc, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	// Create the runner
	return &Runner{
		module:  ctx.Module(),
		config:  ctx.Config(),
		profile: ctx.Profile(),
		client:  dc,
		ctx:     ctx,
		plan:    ctx.Plan(),
	}, nil
}
