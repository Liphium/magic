package main

import (
	"fmt"

	"github.com/Liphium/magic/pkg/databases/surrealdb"

	"github.com/Liphium/magic/v4"
	"github.com/Liphium/magic/v4/mconfig"
	"github.com/Liphium/magic/v4/scripting"
)

func BuildConfig() magic.Config {
	return magic.Config{
		AppName: "magic-example-surrealdb",
		PlanDeployment: func(ctx *mconfig.Context) {

			// Create the SurrealDB driver
			driver := surrealdb.NewDriver("surrealdb/surrealdb:v3").
				NewDatabase("liphium", "posts")
			ctx.Register(driver)

			// Allocate a new port for the service. This makes it possible to run multiple instances of this app locally, without weird configuration hell. Magic will pick a port in case the preferred one is taken.
			port := ctx.ValuePort(8080)

			// Register the environment variables required for connecting in surreal.go
			ctx.WithEnvironment(mconfig.Environment{
				"SURREAL_HOST": mconfig.ValueWithBase([]mconfig.EnvironmentValue{driver.Host(ctx), driver.Port(ctx)}, func(s []string) string {
					return fmt.Sprintf("ws://%s:%s", s[0], s[1])
				}),
				"SURREAL_USER":     driver.Username(),
				"SURREAL_PASSWORD": driver.Password(),
				"SURREAL_NS":       mconfig.ValueStatic("liphium"),
				"SURREAL_DB":       mconfig.ValueStatic("posts"),

				// Make the server listen on localhost using the port allocated by Magic
				"LISTEN": mconfig.ValueWithBase([]mconfig.EnvironmentValue{port}, func(s []string) string {
					return fmt.Sprintf("127.0.0.1:%s", s[0])
				}),
			})
		},
		StartFunction: Start,
		Scripts: []scripting.Script{
			// Scripts to deal with the database, can always come in handy
			scripting.CreateScript("db-reset", "Reset the database by dropping all tables", ResetDatabase),
			scripting.CreateScript("db-clear", "Clear the database by truncating all tables", ClearDatabases),
			scripting.CreateScript("db-seed", "Seed the database with sample posts", SeedDatabase),

			// Scripts to call endpoints, really useful for tests and development
			scripting.CreateScript("create-post", "Create a post using the endpoint", CreatePost),
			scripting.CreateScript("list-posts", "List posts using the endpoint", PrintPosts),
		},
	}
}
