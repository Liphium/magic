package main

import (
	"context"
	"errors"
	"os"

	"github.com/Liphium/magic/integration"
	"github.com/urfave/cli/v3"
)

// Default gitignore file for magic
const magicGitIgnore string = `
*.DS_Store
databases/
`

// Default config.go file for magic
const magicConfig string = `package config

func run(ctx *mconfig.Context) {
	fmt.Println("Hello magic!")
}
`

// Command: magic init
func initCommand(ctx context.Context, c *cli.Command) error {

	// See if the magic directory already exists
	dir, err := integration.GetMagicDirectory(false)
	if err == nil && dir != nil {
		return errors.New("magic project already exists")
	}

	// Create the .magic directory
	if err := os.Mkdir(".magic", 0755); err != nil {
		return err
	}
	if err := os.Chdir(".magic"); err != nil {
		return err
	}

	// Create all files needed
	if err := createFile(".gitignore", magicGitIgnore); err != nil {
		return err
	}
	if err := createFile("config.go", magicConfig); err != nil {
		return err
	}

	return nil
}

// Create a new file with content.
func createFile(name string, content string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write([]byte(content))
	return err
}
