package mrunner

import "github.com/Liphium/magic/mconfig"

type Database struct {
	name string
}

const DefaultDatabaseUser = ""

// Create a new database from one in the config
func (r *Runner) CreateDatabaseFrom(db *mconfig.Database) {

	r.databases = append(r.databases, &Database{
		name: db.Name(),
	})
}
