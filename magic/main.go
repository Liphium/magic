package main

import (
	"context"
	"log"
	"os"

	init_command "github.com/Liphium/magic/magic/init"
	start_command "github.com/Liphium/magic/magic/start"
	test_command "github.com/Liphium/magic/magic/test"
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
