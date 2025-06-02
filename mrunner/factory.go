package mrunner

import (
	"fmt"
	"os"
	"path/filepath"
)

type Factory struct {
	mDir string // Magic directory as a base directory
}

// Create a new factory from the magic directory.
func NewFactory(mDir string) Factory {
	return Factory{
		mDir: mDir,
	}
}

// Get the magic directory
func (f Factory) MagicDirectory() string {
	return f.mDir
}

// Get the directory of the module
func (f Factory) ModuleDirectory() (string, error) {
	return filepath.Abs(filepath.Join(f.mDir, "../"))
}

// Get the module name and version
func (f Factory) ModuleNameAndVersion() (module string, version string, err error) {

	// Get the module directory
	modDir, err := f.ModuleDirectory()
	if err != nil {
		return "", "", fmt.Errorf("couldn't get mod dir: %s", err)
	}

	// Search the file contents
	content, err := os.ReadFile(filepath.Join(modDir, "go.mod"))
	if err != nil {
		return "", "", fmt.Errorf("couldn't read go.mod: %s", err)
	}
	filters := []Filter{FilterModFileGoVersion, FilterModFileModuleName}
	results := ScanLinesSanitize(string(content), filters, &OneLineCommentReplacer{})

	// Make sure the searching actually returned good results
	ver := results[FilterModFileGoVersion]
	mod := results[FilterModFileModuleName]
	if len(ver) != 1 || len(mod) != 1 {
		return "", "", fmt.Errorf("couldn't find module name or version")
	}
	return mod[0], ver[0], nil
}
