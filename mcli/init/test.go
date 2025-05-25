package init_command

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/Liphium/magic/integration"
)

const defaultTestBase = `package magic_tests

import (
	"fmt"

	"github.com/Liphium/magic/mtest"
)

func Run%s(p *mtest.Plan) {
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

	// Evaluate the filepath
	_, filename, path, err := integration.EvaluateNewPath(filepath.Join(mDir, "tests", fp))
	if err != nil {
		return fmt.Errorf("bad path "+fp+": %w", err)
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
