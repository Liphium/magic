package main

import (
	"github.com/Liphium/magic/mcli"
	"github.com/joho/godotenv"
)

func main() {

	// Load environment file
	godotenv.Load()

	// Run the cli
	mcli.RunCli()
}
