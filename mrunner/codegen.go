package mrunner

import (
	"fmt"
)

const runFile = `package main

import (
	"os"
	"log"
	"fmt"
	
	"github.com/Liphium/magic/mconfig"
	"github.com/Liphium/magic/mrunner"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatalln("Please specify config, profile and magic directory!")
	}
	config := os.Args[1]
	profile := os.Args[2]
	magicDir := os.Args[3]

	// Create context
	context := mconfig.DefaultContext(config, profile, magicDir)
	Run(context)

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
	Start()
`

// Generate the run file calling the runner
func GenerateRunFile(deployContainers bool) string {
	if deployContainers {
		return fmt.Sprintf(runFile, runFileDeployer)
	}
	return fmt.Sprintf(runFile, "")
}
