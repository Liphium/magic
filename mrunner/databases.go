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

// Create a new database from one in the config
func (r *Runner) CreateDatabaseFrom(db *mconfig.Database) {
	r.databases = append(r.databases, &Database{
		dbType:      db.Type(),
		name:        db.Name(),
		initialized: false,
	})
}
