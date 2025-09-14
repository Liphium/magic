package starter

import (
	"fmt"

	"github.com/Liphium/magic"
	"github.com/Liphium/magic/mconfig"
	"github.com/Liphium/magic/scripting"
)

func BuildMagicConfig() magic.Config {
	return magic.Config{
		AppName: "magic-example-real-project",
		PlanDeployment: func(ctx *mconfig.Context) {
			// Create a PostgreSQL database for the posts service
			postsDB := mconfig.NewPostgresDatabase("posts")
			ctx.AddDatabase(postsDB)

			// Allocate a new port for the service. This makes it possible to run multiple instances of this app
			// locally, without weird configuration hell. Magic will pick a port in case the preferred one is taken.
			port := ctx.ValuePort(8080)

			// Set up environment variables for the application
			ctx.WithEnvironment(mconfig.Environment{
				// Database connection environment variables
				"DB_HOST":     postsDB.Host(ctx),
				"DB_PORT":     postsDB.Port(ctx),
				"DB_USER":     postsDB.Username(),
				"DB_PASSWORD": postsDB.Password(),
				"DB_DATABASE": postsDB.DatabaseName(ctx),

				// Make the server listen on localhost using the port allocated by Magic
				"LISTEN": mconfig.ValueWithBase([]mconfig.EnvironmentValue{port}, func(s []string) string {
					return fmt.Sprintf("127.0.0.1:%s", s[0])
				}),
			})

			// Load any additional secrets from a .env file if it exists, you could use this to load additional credentials
			// for services Magic might not support (yet c:).
			// _ = ctx.LoadSecretsToEnvironment(".env")
		},
		StartFunction: Start,
		Scripts: []scripting.Script{
			// Scripts to deal with the database, can always come in handy
			scripting.CreateScript("db-reset", "Reset the database by dropping and recreating all tables", resetDatabase),
			scripting.CreateScript("db-seed", "Seed the database with sample posts", seedDatabase),

			// Scripts to call endpoints, really useful for tests and development
			scripting.CreateScript("create-post", "Create a post using the endpoint", createPost),
			scripting.CreateScript("list-posts", "List posts using the endpoint", printPosts),
		},
	}
}
