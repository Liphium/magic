package databases

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Liphium/magic/v2/mconfig"
	"github.com/moby/moby/client"
)

// Handles the instructions for PostgreSQL.
// Supports the following instructions currently:
// - Clear tables
// - Drop tables
func (pd *PostgresDriver) HandleInstruction(ctx context.Context, c *client.Client, container mconfig.ContainerInformation, instruction mconfig.Instruction) error {
	switch instruction {
	case mconfig.InstructionClearTables:
		return pd.ClearTables(container)
	case mconfig.InstructionDropTables:
		return pd.DropTables(container)
	}
	return nil
}

// iterateTablesFn is a function that processes each table in the database
type iterateTablesFn func(tableName string, conn *sql.DB) error

// iterateTables iterates through all tables in all databases and applies the given function
func (pd *PostgresDriver) iterateTables(container mconfig.ContainerInformation, fn iterateTablesFn) error {
	// For all databases, connect and iterate tables
	for _, db := range pd.Databases {
		if err := func() error {
			connStr := fmt.Sprintf("host=127.0.0.1 port=%d user=postgres password=postgres dbname=%s sslmode=disable", container.Ports[0], db)

			// Connect to the database
			conn, err := sql.Open("postgres", connStr)
			if err != nil {
				return fmt.Errorf("couldn't connect to postgres: %v", err)
			}
			defer conn.Close()

			// Get all of the tables
			res, err := conn.Query("SELECT table_name FROM information_schema.tables WHERE table_schema NOT IN ('pg_catalog', 'information_schema')")
			if err != nil {
				return fmt.Errorf("couldn't get database tables: %v", err)
			}
			for res.Next() {
				var name string
				if err := res.Scan(&name); err != nil {
					return fmt.Errorf("couldn't get database table name: %v", err)
				}
				if err := fn(name, conn); err != nil {
					return err
				}
			}

			return nil
		}(); err != nil {
			return err
		}
	}

	return nil
}

// Clear all tables in all databases (keeps table schema alive, just removes the content of all tables)
func (pd *PostgresDriver) ClearTables(container mconfig.ContainerInformation) error {
	return pd.iterateTables(container, func(tableName string, conn *sql.DB) error {
		if _, err := conn.Exec(fmt.Sprintf("truncate %s CASCADE", tableName)); err != nil {
			return fmt.Errorf("couldn't truncate table %s: %v", tableName, err)
		}
		return nil
	})
}

// Drop all tables in all databases (actually deletes all of your tables)
func (pd *PostgresDriver) DropTables(container mconfig.ContainerInformation) error {
	return pd.iterateTables(container, func(tableName string, conn *sql.DB) error {
		if _, err := conn.Exec(fmt.Sprintf("DROP TABLE %s CASCADE", tableName)); err != nil {
			return fmt.Errorf("couldn't drop table table %s: %v", tableName, err)
		}
		return nil
	})
}
