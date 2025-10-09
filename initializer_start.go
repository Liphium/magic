package magic

import (
	"fmt"
	"os"

	"github.com/Liphium/magic/v2/util"
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

	// Make sure to unlock the lock filecin case the app crashes
	defer func() {
		recover() // Make sure the file is always unlocked, even if the function below panics
		factory.Unlock()
	}()

	// Start the app
	config.StartFunction()

	util.Log.Println("Shutting down...")
	os.Exit(0)
}
