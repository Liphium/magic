package integration

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"

	"github.com/fsnotify/fsnotify"
)

// exclusions has to be a list of non-relative paths
func WatchDirectory(dir string, listener func(), exclusions ...string) error {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	// Start listening for events.
	go func() {
		for {
			select {
			case _, ok := <-watcher.Events:
				if !ok {
					return
				}
				listener()
			case err, ok := <-watcher.Errors:
				if !ok {
					log.Println("error not okay")
					return
				}
				log.Println("watch error:", err)
			}
		}
	}()

	// Start watching all of the directories recursively
	for i, exclusion := range exclusions {
		exclusions[i] = filepath.Clean(exclusion)
	}
	return startWatchingRecursive(watcher, dir, exclusions)
}

// Helper function that calls itself recursively adding all directories to the watcher
func startWatchingRecursive(watcher *fsnotify.Watcher, dir string, cleanedExclusions []string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("couldn't read directory %s: %s", dir, err)
	}

	log.Println("Watching", dir)
	if err := watcher.Add(dir); err != nil {
		return fmt.Errorf("couldn't watch directory %s: %s", dir, err)
	}

	for _, entry := range entries {
		entryPath := filepath.Clean(filepath.Join(dir, entry.Name()))

		// If it's a directory and not an exclusion, watch it as well
		if entry.IsDir() && !slices.ContainsFunc(cleanedExclusions, func(path string) bool {
			return filepath.Clean(path) == entryPath
		}) {
			if err := startWatchingRecursive(watcher, entryPath, cleanedExclusions); err != nil {
				return err
			}
		}
	}
	return nil
}
