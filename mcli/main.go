package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
)

func main() {

	// Load environment file
	godotenv.Load()

	var startPath = ""
	var initScriptPath = ""
	var initTestPath = ""
	var startProfile = ""
	cmd := &cli.Command{
		Description: "Testing and debugging like Magic.",
		Commands: []*cli.Command{
			{
				Name:        "init",
				Description: "Magically initialize a new project.",
				Action:      initCommand,
				Commands: []*cli.Command{
					{
						Name:        "script",
						Description: "Magically generates a script template.",
						Action:      func(ctx context.Context, c *cli.Command) error { return initScriptCommand(initScriptPath) },
						Arguments: []cli.Argument{
							&cli.StringArg{
								Name:        "path",
								Destination: &initScriptPath,
							},
						},
					},
					{
						Name:        "test",
						Description: "Magically generates a test template.",
						Action:      func(ctx context.Context, c *cli.Command) error { return initTestCommand(initTestPath) },
						Arguments: []cli.Argument{
							&cli.StringArg{
								Name:        "path",
								Destination: &initTestPath,
							},
						},
					},
				},
			},
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
