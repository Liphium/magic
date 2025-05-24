package mrunner

import (
	"fmt"

	"github.com/Liphium/magic/integration"
	"github.com/Liphium/magic/mconfig"
)

func (r *Runner) Deploy() {

	// Prepare database containers

}

func (r *Runner) PrepareDatabases() error {

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
