package main

import (
	"context"
	"errors"
	"fmt"
	"log"
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

import (
	"fmt"

	"github.com/Liphium/magic/mconfig"
)

// This is the function called once you run the project
func Run(ctx *mconfig.Context) {
	fmt.Println("Hello magic!")
}

func Start() {
	// TODO: Run your application here
}
`

// Command: magic init
func initCommand(ctx context.Context, c *cli.Command) error {

	// See if the magic directory already exists
	_, err := integration.GetMagicDirectory(1)
	if err == nil {
		return errors.New("magic project already exists")
	}

	// Create the .magic directory
	log.Println("Creating magic folder..")
	if err := os.Mkdir("magic", 0755); err != nil {
		return err
	}
	if err := os.Chdir("magic"); err != nil {
		return err
	}

	// Create all files needed
	log.Println("Creating files..")
	if err := integration.CreateFileWithContent(".gitignore", magicGitIgnore); err != nil {
		return err
	}
	if err := integration.CreateFileWithContent("config.go", magicConfig); err != nil {
		return err
	}

	// Create all directories needed
	log.Println("Creating directories..")
	if err := os.Mkdir("scripts", 0755); err != nil {
		log.Fatalln("Failed to create scripts folder: ", err)
	}
	if err := os.Mkdir("tests", 0755); err != nil {
		log.Fatalln("Failed to create tests folder: ", err)
	}

	// Run go mod tidy
	log.Println("Importing packages..")
	if err := os.Chdir(".."); err != nil {
		return err
	}
	dir, err := os.Getwd()
	if err != nil{
		log.Fatalln("Failed to get cwd: ", err)
	}
	fmt.Println("currently in ", dir)
	err = integration.ExecCmdWithFunc(func(s string) {
		fmt.Println(s)
	}, false, "go", "mod", "tidy")
	if err != nil {
		log.Fatalln("Failed to tidy: ", err)
	}
	err = integration.ExecCmdWithFunc(func(s string) {
		fmt.Println(s)
	}, false, "go", "get", "github.com/Liphium/magic/mconfig")
	if err != nil {
		log.Fatalln("Failed to go get: ", err)
	}

	// Print success message
	log.Println()
	log.Println("Successfully initialized project.")
	log.Println("Use magic init script/test <name> to create new tests/scripts.")
	log.Println("Let's hope you become a great wizard!")

	return nil
}


