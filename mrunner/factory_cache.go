package mrunner

import (
	"fmt"
	"os"
	"path/filepath"
)

// Get the cache directory
func (f Factory) CacheDirectory() string {
	return filepath.Join(f.mDir, "cache")
}

// Get the directory of a script in the cache directory.
func (f Factory) ScriptCacheDirectory(script string) string {
	return filepath.Join(f.CacheDirectory(), "script_"+script)
}

// Get the directory of a test in the cache directory.
func (f Factory) TestCacheDirectory(script string) string {
	return filepath.Join(f.CacheDirectory(), "test_"+script)
}

// Get the directory of a config in the cache directory.
func (f Factory) ConfigCacheDirectory(config string, profile string) string {
	return filepath.Join(f.CacheDirectory(), config+"_"+profile)
}

const generatedWorkFile = `go %s
`

// Create work file or update in case already there.
func (f Factory) CreateCacheWorkFile(version string) error {

	// Create work file if it doesn't exist
	workPath := filepath.Join(f.CacheDirectory(), "go.work")
	contentBytes, err := os.ReadFile(workPath)
	if err != nil {

		// Create the work file with the version from the go mod
		if err := os.WriteFile("go.work", []byte(fmt.Sprintf(generatedWorkFile, version)), 0755); err != nil {
			return fmt.Errorf("couldn't write go.work: %s", err)
		}
		return nil
	}

	// Replace the go version in case it already exists
	content := ReplaceLinesSanitized(string(contentBytes), &GoVersionReplacer{
		Version: version,
	}, &OneLineCommentReplacer{})
	if err := os.WriteFile(workPath, []byte(content), 0755); err != nil {
		return fmt.Errorf("couldn't replace go.work: %s", err)
	}

	return nil
}
