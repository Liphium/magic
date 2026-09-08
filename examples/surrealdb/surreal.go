package main

import (
	"context"
	"embed"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/models"
)

//go:embed schemas/*
var schemaFS embed.FS

var Surreal *surrealdb.DB

// ConnectSurreal connects to SurrealDB.
func ConnectSurreal() {
	if Surreal != nil {
		return
	}

	ctx := context.Background()
	var err error
	Surreal, err = surrealdb.FromEndpointURLString(ctx, os.Getenv("SURREAL_HOST"))
	if err != nil {
		slog.Error("Failed to connect to SurrealDB", "err", err)
		os.Exit(1)
	}

	if _, err := Surreal.SignIn(ctx, surrealdb.Auth{
		Username: os.Getenv("SURREAL_USER"),
		Password: os.Getenv("SURREAL_PASSWORD"),
	}); err != nil {
		slog.Error("Failed to sign in to SurrealDB", "err", err)
		os.Exit(1)
	}

	if err := Surreal.Use(ctx, os.Getenv("SURREAL_NS"), os.Getenv("SURREAL_DB")); err != nil {
		slog.Error("Failed to select SurrealDB namespace and database", "err", err)
		os.Exit(1)
	}

	slog.Info("Connected to SurrealDB", "host", os.Getenv("SURREAL_HOST"), "db", os.Getenv("SURREAL_DB"), "ns", os.Getenv("SURREAL_NS"))

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
	if err != nil || current == nil {
		if err != nil && !surrealdb.IsNotFound(err) {
			slog.Error("Failed to get migration state", "err", err)
			os.Exit(1)
		}

		if _, err := surrealdb.Create[any](ctx, Surreal, migrationId, MigrationData{
			Version: 0, // 1 has not completed yet
		}); err != nil {
			slog.Error("Failed to create migration state", "err", err)
			os.Exit(1)
		}
		current = &MigrationData{
			Version: 0,
		}
	}

	files, err := schemaFS.ReadDir("schemas")
	if err != nil {
		slog.Error("Failed to read schema directory", "err", err)
		os.Exit(1)
	}
	var biggest int = 0
	var migrations = map[int]string{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		version, _, found := strings.Cut(file.Name(), "_")
		if !found {
			slog.Error("Invalid migration file name", "file", file.Name())
			os.Exit(1)
		}
		number, err := strconv.Atoi(version)
		if err != nil {
			slog.Error("Invalid migration file name", "file", file.Name(), "err", err)
			os.Exit(1)
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
			slog.Error("Missing migration file", "version", version)
			os.Exit(1)
		}

		schema, err := schemaFS.ReadFile(migration)
		if err != nil {
			slog.Error("Failed to read migration file", "err", err)
			os.Exit(1)
		}

		if _, err := surrealdb.Query[any](ctx, Surreal, string(schema), nil); err != nil {
			slog.Error("Failed to run migration", "err", err)
			os.Exit(1)
		}

		if _, err := surrealdb.Update[any](ctx, Surreal, migrationId, MigrationData{
			Version: version,
		}); err != nil {
			slog.Error("Failed to update migration state", "err", err)
			os.Exit(1)
		}

		slog.Info("SurrealDB migration completed", "version", version, "file", migration)
		current.Version++
	}

	slog.Info("Migrations for SurrealDB successfully completed", "version", current.Version)
}
