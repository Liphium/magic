package surrealdb

import (
	"context"
	"embed"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/models"
)

//go:embed schemas/*
var schemaFS embed.FS

var Surreal *surrealdb.DB

// ConnectSurreal connects to SurrealDB.
func ConnectSurreal() {
	ctx := context.Background()
	var err error
	Surreal, err = surrealdb.FromEndpointURLString(ctx, os.Getenv("SURREAL_HOST"))
	if err != nil {
		util.Fatal(log, "Failed to connect to SurrealDB", "err", err)
	}
	defer Surreal.Close(context.Background())

	if _, err := Surreal.SignIn(ctx, surrealdb.Auth{
		Username: os.Getenv("SURREAL_USER"),
		Password: os.Getenv("SURREAL_PASSWORD"),
	}); err != nil {
		log.Fatal("Failed to sign in to SurrealDB", err)
	}

	if err := Surreal.Use(ctx, os.Getenv("SURREAL_DB"), os.Getenv("SURREAL_NS")); err != nil {
		log.Fatal("Failed to select SurrealDB database and namespace", err)
	}

	log.Info("Connected to SurrealDB", "host", os.Getenv("SURREAL_HOST"), "db", os.Getenv("SURREAL_DB"), "ns", os.Getenv("SURREAL_NS"))

	runMigrations()
}

type MigrationData struct {
	Version int `json:"version"`
}

var migrationId = models.NewRecordID("migration", "state")

func runMigrations() {
	ctx := context.Background()

	// Get current migration version
	current, err := surrealdb.Select[MigrationData](ctx, Surreal, migrationId)
	if err != nil {
		if _, ok := errors.AsType[*surrealdb.ServerError](err); !ok {
			util.Fatal(log, "Failed to get migration state", "err", err)
		}

		_, err := surrealdb.Create[any](ctx, Surreal, migrationId, MigrationData{
			Version: 0, // 1 has not completed yet
		})
		if err != nil {
			util.Fatal(log, "Failed to create migration state", "err", err)
		}
		current = &MigrationData{
			Version: 0,
		}
	}

	files, err := schemaFS.ReadDir("schemas")
	if err != nil {
		util.Fatal(log, "Failed to read schema directory", "err", err)
	}
	var biggest int = 0
	var migrations = map[int]string{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		version, _, found := strings.Cut(file.Name(), "_")
		if !found {
			util.Fatal(log, "Invalid migration file name", "file", file.Name())
		}
		number, err := strconv.Atoi(version)
		if err != nil {
			util.Fatal(log, "Invalid migration file name", "file", file.Name(), "err", err)
		}

		if number > biggest {
			biggest = number
		}
		migrations[number] = filepath.Join("schemas", file.Name())
	}

	for version := 1; version <= biggest; version++ {
		if version <= current.Version {
			continue
		}
		migration, ok := migrations[version]
		if !ok {
			util.Fatal(log, "Missing migration file", "version", version)
		}

		schema, err := schemaFS.ReadFile(migration)
		if err != nil {
			util.Fatal(log, "Failed to read migration file", "err", err)
		}

		if _, err := surrealdb.Query[any](ctx, Surreal, string(schema), nil); err != nil {
			util.Fatal(log, "Failed to run migration", "err", err)
		}

		if _, err := surrealdb.Update[any](ctx, Surreal, migrationId, MigrationData{
			Version: version,
		}); err != nil {
			util.Fatal(log, "Failed to update migration state", "err", err)
		}

		log.Info("SurrealDB migration completed", "version", version, "file", migration)
	}

	log.Info("Migrations for SurrealDB successfully completed", "version", current.Version)
}
