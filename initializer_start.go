package magic

import (
	"fmt"
	"os"

	"github.com/Liphium/magic/v4/util"
	"github.com/spf13/pflag"
)

// Start your application with Magic. Make sure to provide all the arguments in the config that are required.
func Start(config Config) {

	// Parse flags before preparing
	pflag.Parse()

	factory, runner := prepare(config, "")
	if factory == nil || runner == nil {
		fmt.Println()
		return
	}
	util.Log.Println("Successfully prepared everything!")
	fmt.Println()

	// Start the app
	config.StartFunction()

	util.Log.Println("Shutting down...")
	os.Exit(0)
}
