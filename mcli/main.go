package main

import (
	"context"
	"log"
	"os"

	init_command "github.com/Liphium/magic/mcli/init"
	start_command "github.com/Liphium/magic/mcli/start"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
)

func main() {

	// Load environment file
	godotenv.Load()

	cmd := &cli.Command{
		Description: "Testing and debugging like Magic.",
		Commands: []*cli.Command{
			init_command.BuildCommand(),
			start_command.BuildCommand(),
			{
				Name:        "test",
				Description: "Magically test your project.",
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatalln("ERROR:", err)
	}
}
