package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Arguments: []cli.Argument{
			&cli.IntArg{
				Name: "someint",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			fmt.Printf("We got %d", cmd.IntArg("someint"))
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
