package init_command

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Liphium/magic/integration"
)

const defaultTestBase = `package magic_tests

import (
	"fmt"

	"github.com/Liphium/magic/mconfig"
)

func Run%s(p *mconfig.Plan) {
	fmt.Println("Hello, I'm the greatest wizzard of all time!")
}
`

// Command: magic init test <name>
func initTestCommand(fp string) error {

	// Get magic dir
	mDir, err := integration.GetMagicDirectory(5)
	if err != nil {
		return err
	}

	// Create magic/tests if it doesn't exist
	if sE, err := integration.DoesDirExist(filepath.Join(mDir, "tests")); err != nil {
		return err
	} else if sE {
		log.Println("Creating tests folder..")
		if err = os.Mkdir(filepath.Join(mDir, "tests"), 0755); err != nil {
			log.Fatalln("Failed to create tests folder: %w", err)
		}
	}

	// Evaluate the filepath
	_, filename, path, err := integration.EvaluateNewPath(filepath.Join(mDir, "tests", fp))
	if err != nil {
		return fmt.Errorf("bad path %q: %w", fp, err)
	}

	// Generate test base
	testName := integration.SnakeToCamelCase(strings.TrimRight(filename, ".go"), true)
	testBase := fmt.Sprintf(defaultTestBase, testName)

	// Create all files needed
	log.Println("Creating test..")
	if err := integration.CreateFileWithContent(path, testBase); err != nil {
		return err
	}

	// Print success message
	log.Println()
	log.Println("Successfully created test template.")

	return nil
}
