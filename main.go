package main

import (
	"os"

	"github.com/Liphium/magic/mcli"
	"github.com/Liphium/magic/mconfig"
	"github.com/joho/godotenv"
)

func main() {

	// Load environment file
	godotenv.Load()

	// Set verbose logging
	mconfig.VerboseLogging = os.Getenv("MAGIC_VERBOSE") == "true"

	// Run the cli
	mcli.RunCli()
}
