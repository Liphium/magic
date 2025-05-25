package main

import (
	"context"
	"log"
	"os"

	init_command "github.com/Liphium/magic/mcli/init"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
)

func main() {

	// Load environment file
	godotenv.Load()

	var startPath = ""
	var startProfile = ""
	cmd := &cli.Command{
		Description: "Testing and debugging like Magic.",
		Commands: []*cli.Command{
			init_command.BuildCommand(),
			{
				Name:        "start",
				Description: "Magically start your project.",
				Action:      func(ctx context.Context, c *cli.Command) error { return startCommand(startPath, startProfile) },
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "profile",
						Aliases:     []string{"p"},
						Value:       "",
						Destination: &startProfile,
						Usage:       "To run multiple instances of the same magic config",
					},
				},
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:        "path",
						Destination: &startPath,
					},
				},
			},
			{
				Name:        "test",
				Description: "Magically test your project.",
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatalln("ERROR: ", err)
	}
}
