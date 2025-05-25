package init_command

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Liphium/magic/integration"
)

const defaultScriptBase = `package magic_scripts

import (
	"fmt"

	"github.com/Liphium/magic/mconfig"
)

func Run%s(p *mconfig.Plan) {
	fmt.Println("I'm a wizzard!")
}
`

// Command: magic init script <name>
func initScriptCommand(fp string) error {

	// Get magic dir
	mDir, err := integration.GetMagicDirectory(5)
	if err != nil {
		return err
	}

	// Create magic/scripts if it doesn't exist
	if sE, err := integration.DoesDirExist(filepath.Join(mDir, "scripts")); err != nil {
		return err
	} else if sE {
		log.Println("Creating scripts folder..")
		if err = os.Mkdir(filepath.Join(mDir, "scripts"), 0755); err != nil {
			log.Fatalln("Failed to create scripts folder: %w", err)
		}
	}

	// Evaluate the filepath
	_, filename, path, err := integration.EvaluateNewPath(filepath.Join(mDir, "scripts", fp))
	if err != nil {
		return fmt.Errorf("bad path %s: %w", fp, err)
	}

	// Generate Script base
	scriptName := integration.SnakeToCamelCase(strings.TrimRight(filename, ".go"), true)
	scriptBase := fmt.Sprintf(defaultScriptBase, scriptName)

	// Create all files needed
	log.Println("Creating script..")
	if err := integration.CreateFileWithContent(path, scriptBase); err != nil {
		return err
	}

	// Print success message
	log.Println()
	log.Println("Successfully created script template.")

	return nil
}
