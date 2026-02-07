package mrunner

import (
	"github.com/Liphium/magic/v2/mconfig"
	"github.com/moby/moby/client"
)

const DefaultStartPort uint = 10000
const DefaultEndPort uint = 60000

type Runner struct {
	appName  string
	profile  string
	client   *client.Client
	ctx      *mconfig.Context
	plan     *mconfig.Plan
	services []mconfig.ServiceDriver
}

func (r *Runner) Environment() *mconfig.Environment {
	return r.ctx.Environment()
}

// Create a new runner
func NewRunner(ctx *mconfig.Context) (*Runner, error) {

	// Create a new client for the docker sdk
	dc, err := client.New(client.FromEnv)
	if err != nil {
		return nil, err
	}

	// Create the runner
	return &Runner{
		appName:  ctx.AppName(),
		profile:  ctx.Profile(),
		client:   dc,
		ctx:      ctx,
		plan:     ctx.Plan(),
		services: ctx.Services(),
	}, nil
}

// Create a new runner
func NewRunnerFromPlan(plan *mconfig.Plan) (*Runner, error) {

	// Create a new client for the docker sdk
	dc, err := client.New(client.FromEnv)
	if err != nil {
		return nil, err
	}

	// Create the runner
	return &Runner{
		client: dc,
		ctx:    nil,
		plan:   plan,
	}, nil
}
