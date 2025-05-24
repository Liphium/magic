package mconfig

type DatabaseType = uint

const (
	DatabasePostgres DatabaseType = iota
)

type Database struct {
	dbType DatabaseType // Type of the database
	name   string
}

func (db *Database) Type() DatabaseType {
	return db.dbType
}

// Get the name of the database (for magic s)
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

// Get the name of the database for environment variables
func (db *Database) DatabaseName(ctx *Context) EnvironmentValue {
	return ValueStatic(db.DefaultDatabaseName(ctx))
}

// Get the port of the database for environment variables
func (db *Database) Port() EnvironmentValue {
	return ValueStatic("hi")
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
	switch db.dbType {
	case DatabasePostgres:
		return "postgres"
	default:
		return "admin"
	}
}

// Get the default username for the database type
func (db *Database) DefaultUsername() string {
	switch db.dbType {
	case DatabasePostgres:
		return "postgres"
	default:
		return "admin"
	}
}

// Get the default name for the database using the runner
func (db *Database) DefaultDatabaseName(ctx *Context) string {
	return "mgc-" + ctx.config + "-" + ctx.profile + "-" + db.name
}

// Create a new Postgres database.
func NewPostgresDatabase(name string) *Database {
	return &Database{
		dbType: DatabasePostgres,
		name:   name,
	}
}
