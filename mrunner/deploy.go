package mrunner

import (
	"fmt"
	"log"

	"github.com/Liphium/magic/integration"
	"github.com/Liphium/magic/mconfig"
)

// Deploy the containers for the magic context
func (r *Runner) Deploy(ctx *mconfig.Context) {

	// Prepare database containers
	if err := r.prepareDatabases(); err != nil {
		log.Fatalln("couldn't start databases:", err)
	}

	// Add all of the environment variables
	if err := ctx.Environment().Apply(); err != nil {
		log.Fatalln("couldn't set environment variables:", err)
	}

	// TODO: Run the containers
}

func (r *Runner) prepareDatabases() error {

	// Scan for open ports per type
	ports := map[mconfig.DatabaseType]uint{}
	for _, db := range r.databases {
		if ports[db.dbType] == 0 {
			var err error
			ports[db.dbType], err = integration.ScanForOpenPort(DefaultStartPort, DefaultEndPort)
			if err != nil {
				return fmt.Errorf("couldn't find open port: %e", err)
			}
		}
	}

	// TODO: Add databases credentials and stuff

	return nil
}
