package main

import (
	"context"
	"log"
	"os"

	init_command "github.com/Liphium/magic/mcli/init"
	start_command "github.com/Liphium/magic/mcli/start"
	test_command "github.com/Liphium/magic/mcli/test"
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
			test_command.BuildCommand(),
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatalln("ERROR:", err)
	}
}
