package magic

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofrs/flock"
)

type Factory struct {
	projectDir string
	lock       *flock.Flock
}

// Create a new factory, will search for the current project directory
func createFactory() (Factory, error) {
	dir, err := os.Getwd()
	if err != nil {
		return Factory{}, err
	}

	for i := 0; i < 3; i++ {
		modPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(modPath); err == nil {
			return Factory{projectDir: dir}, nil
		}
		dir = filepath.Dir(dir)
	}
	return Factory{}, errors.New("could not find project directory")
}

// Print a warning in case .magic is not in the current .gitignore
func (f Factory) WarnIfNotIgnored() {
	gitignorePath := filepath.Join(f.projectDir, ".gitignore")

	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		return
	}

	for line := range strings.Lines(string(content)) {
		if strings.Contains(line, ".magic") {
			return
		}

	}
	Log.Println("WARNING: .magic is not in your .gitignore file.")
}

// Get the current magic cache directory
func (f Factory) MagicDirectory() string {
	return filepath.Join(f.projectDir, ".magic")
}

// Get the location of the lock file for a profile
func (f Factory) LockFile(profile string) string {
	return filepath.Join(f.MagicDirectory(), fmt.Sprintf("%s.lock", profile))
}

// Get the location of the plan file for a profile
func (f Factory) PlanFile(profile string) string {
	return filepath.Join(f.MagicDirectory(), fmt.Sprintf("%s.mplan", profile))
}

// Trys to lock the lock file for the profile. If no error is returned, the profile was locked.
func (f Factory) TryLockProfile(profile string) error {
	fileLock := flock.New(f.LockFile(profile))

	locked, err := fileLock.TryLock()
	if err != nil {
		return err
	}
	if !locked {
		return errors.New("couldn't lock profile lock file")
	}

	f.lock = fileLock
	return nil
}

// Unlock the lock of the current factory in case one is there
func (f Factory) Unlock() error {
	if f.lock != nil {
		return f.lock.Unlock()
	}
	return nil
}
