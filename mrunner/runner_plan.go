package mrunner

import (
	"fmt"
	"maps"

	"github.com/Liphium/magic/v2/integration"
	"github.com/Liphium/magic/v2/mconfig"
	"github.com/Liphium/magic/v2/util"
)

// Get the current plan (might not be set yet, call GeneratePlan first)
func (r *Runner) Plan() *mconfig.Plan {
	return r.plan
}

// Deploy the containers for the magic context
func (r *Runner) GeneratePlan() *mconfig.Plan {
	if r.ctx == nil {
		util.Log.Fatalln("no context set")
	}

	// Prepare database containers
	types, err := r.prepareDatabases()
	if err != nil {
		util.Log.Fatalln("couldn't start databases:", err)
	}

	// Prepare all of the ports
	allocatedPorts := map[uint]uint{}
	if r.ctx.Ports() != nil {
		for _, port := range r.ctx.Ports() {
			// Generate a new port in case the current one is taken
			toAllocate := port
			if integration.ScanPort(port) {
				toAllocate, err = scanForOpenPort()
				if err != nil {
					util.Log.Fatalln("Couldn't find open port for", port, ":", err)
				}
			}

			// Add the port to the plan
			allocatedPorts[port] = toAllocate
		}
	}

	// Load into plan
	r.plan = &mconfig.Plan{
		AppName:        r.ctx.AppName(),
		Profile:        r.ctx.Profile(),
		DatabaseTypes:  types,
		AllocatedPorts: allocatedPorts,
	}
	r.ctx.ApplyPlan(r.plan)

	// Generate the environment variables and add to plan
	environment := map[string]string{}
	if r.Environment() != nil {
		environment = r.Environment().Generate()
	}
	r.plan.Environment = environment
	return r.plan
}

func (r *Runner) prepareDatabases() ([]mconfig.PlannedDatabaseType, error) {

	// Scan for open ports per type
	types := map[mconfig.DatabaseType]mconfig.PlannedDatabaseType{}
	for _, db := range r.ctx.Databases() {
		if _, ok := types[db.Type()]; !ok {
			openPort, err := scanForOpenPort()
			if err != nil {
				return nil, err
			}

			types[db.Type()] = mconfig.PlannedDatabaseType{
				Type:      db.Type(),
				Port:      openPort,
				Databases: []mconfig.PlannedDatabase{},
			}
		}
	}

	// Add all of the databases
	for _, db := range r.ctx.Databases() {
		dbType := types[db.Type()]
		dbType.Databases = append(dbType.Databases, mconfig.PlannedDatabase{
			Name:     db.Name(),
			Username: db.DefaultUsername(),
			Password: db.DefaultPassword(),
			Hostname: "127.0.0.1",
			Type:     dbType.Type,
			Port:     dbType.Port,
		})
		types[db.Type()] = dbType
	}

	// Convert to list
	plannedTypes := make([]mconfig.PlannedDatabaseType, len(types))
	i := 0
	for value := range maps.Values(types) {
		plannedTypes[i] = value
		i++
	}
	return plannedTypes, nil
}

// Scan for an open port in the default range
func scanForOpenPort() (uint, error) {
	openPort, err := integration.ScanForOpenPort(DefaultStartPort, DefaultEndPort)
	if err != nil {
		return 0, fmt.Errorf("couldn't find open port: %e", err)
	}
	return openPort, err
}
