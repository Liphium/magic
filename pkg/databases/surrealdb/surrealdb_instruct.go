package surrealdb

import (
	"context"
	"fmt"

	"github.com/Liphium/magic/v4/mconfig"
	"github.com/moby/moby/client"
	sdb "github.com/surrealdb/surrealdb.go"
)

// Handles the instructions for SurrealDB. Supports the following instructions currently:
// - Clear tables
// - Drop tables
func (sd *SurrealDriver) HandleInstruction(ctx context.Context, c *client.Client, container mconfig.ContainerInformation, instruction mconfig.Instruction) error {
	switch instruction {
	case mconfig.InstructionClearTables:
		return sd.ClearTables(ctx, container)
	case mconfig.InstructionDropTables:
		return sd.DropTables(ctx, container)
	}
	return nil
}

// Clear all tables in all namespaces/databases (keeps the table definitions intact, just removes all of the records)
func (sd *SurrealDriver) ClearTables(ctx context.Context, container mconfig.ContainerInformation) error {
	return sd.iterateTables(ctx, container, func(tableName string, conn *sdb.DB) error {
		if _, err := sdb.Query[[]map[string]any](ctx, conn, fmt.Sprintf("DELETE `%s`;", tableName), nil); err != nil {
			return fmt.Errorf("couldn't clear table %s: %s", tableName, err)
		}
		return nil
	})
}

// Drop all tables in all namespaces/databases (actually deletes all of your tables)
func (sd *SurrealDriver) DropTables(ctx context.Context, container mconfig.ContainerInformation) error {
	return sd.iterateTables(ctx, container, func(tableName string, conn *sdb.DB) error {
		if _, err := sdb.Query[map[string]any](ctx, conn, fmt.Sprintf("REMOVE TABLE IF EXISTS `%s`;", tableName), nil); err != nil {
			return fmt.Errorf("couldn't drop table %s: %s", tableName, err)
		}
		return nil
	})
}

// iterateTablesFn is a function that processes each table in a database
type iterateTablesFn func(tableName string, conn *sdb.DB) error

// iterateTables iterates through all tables in all namespaces/databases and applies the given function
func (sd *SurrealDriver) iterateTables(ctx context.Context, container mconfig.ContainerInformation, fn iterateTablesFn) error {
	// For all databases, connect and iterate tables
	for _, database := range sd.Databases {
		if err := func() error {
			// Connect to the database
			conn, err := connect(ctx, container)
			if err != nil {
				return fmt.Errorf("couldn't connect to surrealdb: %s", err)
			}
			defer conn.Close(ctx)

			if err := conn.Use(ctx, database.Namespace, database.Database); err != nil {
				return fmt.Errorf("couldn't select surrealdb namespace/database %s/%s: %s", database.Namespace, database.Database, err)
			}

			// Get all of the tables in the current database
			result, err := sdb.Query[map[string]any](ctx, conn, "INFO FOR DB;", nil)
			if err != nil {
				return fmt.Errorf("couldn't get database tables: %s", err)
			}

			tables, ok := (*result)[0].Result["tables"].(map[string]any)
			if !ok {
				return fmt.Errorf("couldn't parse database tables from INFO FOR DB response")
			}
			for name := range tables {
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
