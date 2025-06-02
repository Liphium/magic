package mrunner

import (
	"fmt"
	"log"
	"maps"

	"github.com/Liphium/magic/integration"
	"github.com/Liphium/magic/mconfig"
)

// Deploy the containers for the magic context
func (r *Runner) GeneratePlan() string {

	// Prepare database containers
	types, err := r.prepareDatabases()
	if err != nil {
		log.Fatalln("couldn't start databases:", err)
	}

	// Load into plan
	r.plan = &mconfig.Plan{
		DatabaseTypes: types,
	}
	r.ctx.ApplyPlan(r.plan)

	// Generate the environment variables and add to plan
	environment := map[string]string{}
	if r.Environment() != nil {
		environment = r.Environment().Generate()
	}
	r.plan.Environment = environment

	// Convert plan to printable string
	printable, err := r.plan.ToPrintable()
	if err != nil {
		log.Fatalln("couldn't generate plan:", err)
	}
	return printable
}

func (r *Runner) prepareDatabases() ([]mconfig.PlannedDatabaseType, error) {

	// Scan for open ports per type
	types := map[mconfig.DatabaseType]mconfig.PlannedDatabaseType{}
	for _, db := range r.ctx.Databases() {
		if _, ok := types[db.Type()]; !ok {
			openPort, err := integration.ScanForOpenPort(DefaultStartPort, DefaultEndPort)
			if err != nil {
				return nil, fmt.Errorf("couldn't find open port: %e", err)
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
			ConfigName: db.Name(),
			Name:       mconfig.DefaultDatabaseName(r.config, r.profile, db.Name()),
			Username:   db.DefaultUsername(),
			Password:   db.DefaultPassword(),
			Hostname:   "127.0.0.1",
			Type:       dbType.Type,
			Port:       dbType.Port,
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
