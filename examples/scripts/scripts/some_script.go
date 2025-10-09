package magic_scripts

import (
	"log"

	"github.com/Liphium/magic/v2/mrunner"
)

type SomeScriptOptions struct {
	Name  string `prompt:"Name for the test account." validate:"required"`
	Email string `prompt:"E-Mail address for the test account." validate:"required,email"`
}

// This is the main entrypoint for your script.
// Run it with go run . -r some_script
//
// You can get the runner here if you need it. But you could also just delete the parameter and the code
// would just work the same.
func SomeScript(runner *mrunner.Runner, data SomeScriptOptions) error {
	log.Println("chosen name:", data.Name)
	log.Println("chosen email:", data.Email)
	return nil
}
