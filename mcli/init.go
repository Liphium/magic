package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

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
	log.Println("Creating folder..")
	if err := os.Mkdir(".magic", 0755); err != nil {
		return err
	}
	if err := os.Chdir(".magic"); err != nil {
		return err
	}

	// Create all files needed
	log.Println("Creating files..")
	if err := createFile(".gitignore", magicGitIgnore); err != nil {
		return err
	}
	if err := createFile("config.go", magicConfig); err != nil {
		return err
	}

	// Run go mod tidy
	log.Println("Importing packages..")
	cmd := exec.Command("go", "mod", "tidy")

	// Set up the output pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalln("couldn't open pipe: ", err.Error())
		return nil
	}
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	// Start the command
	if err := cmd.Run(); err != nil {
		log.Fatalln("couldn't import: ", err)
		return err
	}

	// Print success message
	log.Println()
	log.Println("Successfully initialized project.")
	log.Println("Use magic init script/test <name> to create new tests/scripts.")
	log.Println("Let's hope you become a great wizard!")

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
