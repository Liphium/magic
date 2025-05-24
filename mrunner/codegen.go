package mrunner

import "fmt"

const runFile = `package main

import (
	"os"
	"log"
	
	"github.com/Liphium/magic/mconfig"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalln("Please specify config and profile (as first and second argument)!")
	}
	config := os.Args[1]
	profile := os.Args[2]

	// Create context
	context := mconfig.DefaultContext(config, profile)
	run(context)

	// Create the runner from context
	runner, err := mrunner.NewRunner(context)
	if err != nil {
		log.Fatalln("Couldn't create runner:", err)
	}

	fmt.Println("mgc:", runner.GeneratePlan())
%s
}
`

const runFileDeployer = `
	// Deploy the containers and start
	runner.Deploy()

	// Start the app
	start(runner)
`

// Generate the run file calling the runner
func GenerateRunFile(deployContainers bool) string {
	if deployContainers {
		return fmt.Sprintf(runFile, "")
	}
	return fmt.Sprintf(runFile, runFileDeployer)
}
