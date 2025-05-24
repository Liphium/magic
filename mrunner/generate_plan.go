package mrunner

import (
	"fmt"
	"log"
	"maps"

	"github.com/Liphium/magic/integration"
	"github.com/Liphium/magic/mconfig"
	"github.com/Liphium/magic/mtest"
)

// Deploy the containers for the magic context
func (r *Runner) GeneratePlan() string {

	// Prepare database containers
	types, err := r.prepareDatabases()
	if err != nil {
		log.Fatalln("couldn't start databases:", err)
	}

	// Generate the environment variables
	environment := map[string]string{}
	if r.Environment() != nil {
		environment = r.Environment().Generate()
	}

	// Load into plan
	plan := mtest.Plan{
		Environment:   environment,
		DatabaseTypes: types,
	}

	// Convert plan to printable string
	printable, err := plan.ToPrintable()
	if err != nil {
		log.Fatalln("couldn't generate plan:", err)
	}
	return printable
}

func (r *Runner) prepareDatabases() ([]mtest.PlannedDatabaseType, error) {

	// Scan for open ports per type
	types := map[mconfig.DatabaseType]mtest.PlannedDatabaseType{}
	for _, db := range r.databases {
		if _, ok := types[db.Type()]; !ok {
			openPort, err := integration.ScanForOpenPort(DefaultStartPort, DefaultEndPort)
			if err != nil {
				return nil, fmt.Errorf("couldn't find open port: %e", err)
			}

			types[db.Type()] = mtest.PlannedDatabaseType{
				Type:      db.Type(),
				Port:      openPort,
				Databases: []mtest.PlannedDatabase{},
			}
		}
	}

	// Add all of the databases
	for _, db := range r.databases {
		dbType := types[db.Type()]
		dbType.Databases = append(dbType.Databases, mtest.PlannedDatabase{})
		types[db.Type()] = dbType
	}

	// Convert to list
	plannedTypes := make([]mtest.PlannedDatabaseType, len(types))
	for value := range maps.Values(types) {
		plannedTypes = append(plannedTypes, value)
	}
	return plannedTypes, nil
}
