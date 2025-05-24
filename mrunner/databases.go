package mrunner

import "github.com/Liphium/magic/mconfig"

type Database struct {
	dbType mconfig.DatabaseType `json:"t"`
	name   string               `json:"n"`

	initialized bool   `json:"-"`
	database    string `json:"db"`
	username    string `json:"un"`
	password    string `json:"pw"`
	host        string `json:"hn"`
	port        uint   `json:"po"`
}

// Get the default password for the database type
func (db *Database) DefaultPassword() string {
	switch db.dbType {
	case mconfig.DatabasePostgres:
		return "postgres"
	default:
		return "admin"
	}
}

// Get the default username for the database type
func (db *Database) DefaultUsername() string {
	switch db.dbType {
	case mconfig.DatabasePostgres:
		return "postgres"
	default:
		return "admin"
	}
}

// Get the default name for the database using the runner
func (db *Database) DefaultDatabaseName(runner *Runner) string {
	return "mgc-" + runner.config + "-" + runner.profile + "-" + db.name
}

// Create a new database from one in the config
func (r *Runner) CreateDatabaseFrom(db *mconfig.Database) {
	r.databases = append(r.databases, &Database{
		dbType:      db.Type(),
		name:        db.Name(),
		initialized: false,
	})
}
