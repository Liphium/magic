package mconfig

type DatabaseType = uint

const (
	DatabasePostgres DatabaseType = iota
)

type Database struct {
	dbType DatabaseType // Type of the database
	name   string
}

// Get the name of the database
func (db *Database) Name() string {
	return db.name
}

// Get the host of the database for environment variables
func (db *Database) Host() EnvironmentValue {
	return EnvironmentValue{
		get: func() string {
			return "hi"
		},
	}
}

// Get the password of the database for environment variables
func (db *Database) Password() EnvironmentValue {
	return EnvironmentValue{
		get: func() string {
			return "hi"
		},
	}
}

// Get the port of the database for environment variables
func (db *Database) Port() EnvironmentValue {
	return EnvironmentValue{
		get: func() string {
			return "hi"
		},
	}
}

// Get the username of the database for environment variables
func (db *Database) Username() EnvironmentValue {
	return EnvironmentValue{
		get: func() string {
			return "hi"
		},
	}
}

// Create a new Postgres database.
func NewPostgresDatabase(name string) *Database {
	return &Database{
		dbType: DatabasePostgres,
		name:   name,
	}
}
