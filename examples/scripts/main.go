package main

import (
	"fmt"
	magic_scripts "scripts-example/scripts"
	"time"

	"github.com/Liphium/magic/v2"
	"github.com/Liphium/magic/v2/mconfig"
	"github.com/Liphium/magic/v2/scripting"
)

// Our main function is wrapped using Magic, in a real app you should use go build tags to have two main functions, one for
// development with Magic and one without. You can look at our real-project example.
func main() {
	magic.Start(magic.Config{
		AppName: "scripts_example",
		PlanDeployment: func(ctx *mconfig.Context) {
			// We don't need to deploy anything here
		},
		StartFunction: Start,

		// All scripts have to be registered here.
		Scripts: []scripting.Script{
			scripting.CreateScript("some_script", "A random script", magic_scripts.SomeScript),
		},
	})
}

func Start() {
	fmt.Println("Thanks for using Magic.")
	fmt.Println()
	fmt.Println("To run your script use:")
	fmt.Println("go run . -r <script-name> [arguments]")
	fmt.Println("If you don't provide any arguments, you'll be prompted to enter them in a CLI form!")
	fmt.Println()
	fmt.Println("Become a great wizard!")
	time.Sleep(10 * time.Minute)
}
