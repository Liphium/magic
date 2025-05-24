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
func newDB(db *mconfig.Database) *Database {
	return &Database{
		dbType:      db.Type(),
		name:        db.Name(),
		initialized: false,
	}
}
