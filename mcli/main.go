package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Description: "Testing and debugging like Magic.",
		Commands: []*cli.Command{
			{
				Name:        "init",
				Description: "Magically initialize a new project.",
				Action:      initCommand,
			},
			{
				Name:        "start",
				Description: "Magically start your project.",
			},
			{
				Name:        "test",
				Description: "Magically test your project.",
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
