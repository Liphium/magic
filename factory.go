package magic

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Liphium/magic/v2/util"
	"github.com/gofrs/flock"
)

// The maximum amount of folders Magic tries to go back
const maxRecursiveTries = 20

var errProfileLocked = errors.New("profile is already locked by a different instance")

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

	for range maxRecursiveTries {
		modPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(modPath); err == nil {
			return Factory{projectDir: dir}, nil
		}
		dir = filepath.Dir(dir)
	}
	return Factory{}, errors.New("could not find project directory")
}

// Print a warning in case .magic is not in the current .gitignore
func (f *Factory) WarnIfNotIgnored() {
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
	util.Log.Println("WARNING: .magic is not in your .gitignore file.")
}

// Get the current magic cache directory
func (f *Factory) MagicDirectory() string {
	return filepath.Join(f.projectDir, ".magic")
}

// Get the location of the lock file for a profile
func (f *Factory) LockFile(profile string) string {
	return filepath.Join(f.MagicDirectory(), fmt.Sprintf("%s.lock", profile))
}

// Get the location of the plan file for a profile
func (f *Factory) PlanFile(profile string) string {
	return filepath.Join(f.MagicDirectory(), fmt.Sprintf("%s.json", profile))
}

// Check if a profile is locked (a magic instance is running)
func (f *Factory) IsProfileLocked(profile string) bool {
	fileLock := flock.New(f.LockFile(profile))
	locked, err := fileLock.TryLock()
	if err != nil {
		return true
	}
	if locked {
		fileLock.Unlock()
		return false
	}
	return true
}

// Trys to lock the lock file for the profile. If no error is returned, the profile was locked.
func (f *Factory) TryLockProfile(profile string) error {

	// Create the lock file in case it doesn't exist
	if _, err := os.Stat(f.LockFile(profile)); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(f.LockFile(profile)), 0755); err != nil {
			return fmt.Errorf("couldn't create cache directory: %s", err)
		}
		_, err := os.Create(f.LockFile(profile))
		if err != nil {
			return fmt.Errorf("couldn't create lock file: %s", err)
		}
	}

	fileLock := flock.New(f.LockFile(profile))

	locked, err := fileLock.TryLock()
	if err != nil {
		return err
	}
	if !locked {
		return errProfileLocked
	}

	f.lock = fileLock
	return nil
}

// Unlock the lock of the current factory in case one is there
func (f *Factory) Unlock() error {
	if f.lock != nil {
		return f.lock.Unlock()
	}
	return nil
}
