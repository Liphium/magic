package mcli

import (
	"context"
	"log"
	"os"

	init_command "github.com/Liphium/magic/mcli/init"
	start_command "github.com/Liphium/magic/mcli/start"
	test_command "github.com/Liphium/magic/mcli/test"
	"github.com/urfave/cli/v3"
)

func RunCli() {
	cmd := &cli.Command{
		Name:        "magic",
		Usage:       "a tool for testing by Liphium.",
		Description: "Testing and debugging like Magic.",
		Commands: []*cli.Command{
			init_command.BuildCommand(),
			start_command.BuildCommand(),
			test_command.BuildCommand(),
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Println("ERROR:", err)
	}
}
