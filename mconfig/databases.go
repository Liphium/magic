package mconfig

import (
	"fmt"
	"log"
)

type DatabaseType = uint

const (
	DatabasePostgres DatabaseType = 1
)

func IsValidDatabaseType(t DatabaseType) bool {
	return t == DatabasePostgres
}

type Database struct {
	dbType DatabaseType // Type of the database
	name   string
}

func (db *Database) Type() DatabaseType {
	return db.dbType
}

// Get the name of the database (as in the config)
func (db *Database) Name() string {
	return db.name
}

// Get the host of the database for environment variables
func (db *Database) Host(ctx *Context) EnvironmentValue {
	return EnvironmentValue{
		get: func() string {
			planned, err := ctx.Plan().Database(db.name)
			if err != nil {
				log.Fatalln("error while getting database with name", db.name+":", err)
			}
			return planned.Hostname
		},
	}
}

// Get the name of the database for environment variables
func (db *Database) DatabaseName(ctx *Context) EnvironmentValue {
	return EnvironmentValue{
		get: func() string {
			planned, err := ctx.Plan().Database(db.name)
			if err != nil {
				log.Fatalln("error while getting database with name", db.name+":", err)
			}
			return planned.Name
		},
	}
}

// Get the port of the database for environment variables
func (db *Database) Port(ctx *Context) EnvironmentValue {
	return EnvironmentValue{
		get: func() string {
			planned, err := ctx.Plan().Database(db.name)
			if err != nil {
				log.Fatalln("error while getting database with name", db.name+":", err)
			}
			return fmt.Sprintf("%d", planned.Port)
		},
	}
}

// Get the password of the database for environment variables
func (db *Database) Password() EnvironmentValue {
	return ValueStatic(db.DefaultPassword())
}

// Get the username of the database for environment variables
func (db *Database) Username() EnvironmentValue {
	return ValueStatic(db.DefaultUsername())
}

// Get the default password for the database type
func (db *Database) DefaultPassword() string {
	return DefaultPassword(db.dbType)
}

// Get the default username for the database type
func (db *Database) DefaultUsername() string {
	return DefaultUsername(db.dbType)
}

// Get the default name for the database using the runner
func (db *Database) DefaultDatabaseName(ctx *Context) string {
	return DefaultDatabaseName(ctx.config, ctx.profile, db.name)
}

// Get the default password for a database by type.
func DefaultPassword(dbType DatabaseType) string {
	switch dbType {
	case DatabasePostgres:
		return "postgres"
	default:
		return "admin"
	}
}

// Get the default username for a database by type.
func DefaultUsername(dbType DatabaseType) string {
	switch dbType {
	case DatabasePostgres:
		return "postgres"
	default:
		return "admin"
	}
}

// Get the default database name for a database.
func DefaultDatabaseName(config string, profile string, databaseName string) string {
	return "mgc-" + config + "-" + profile + "-" + databaseName
}

// Create a new Postgres database.
func NewPostgresDatabase(name string) *Database {
	return &Database{
		dbType: DatabasePostgres,
		name:   name,
	}
}
