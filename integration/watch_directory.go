package integration

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"

	"github.com/Liphium/magic/mconfig"
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

	if mconfig.VerboseLogging {
		log.Printf("Watching %s...", dir)
	}

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

type WatchContext[J any, C any] struct {
	Print            func(string)                                                      // Called for prints.
	Error            func(error)                                                       // Called when an error happens.
	Start            func(currentContext C, lastJob *J, retrievalChannel chan J) error // Gets called to start the process.
	Stop             func(J) error                                                     // Gets called to stop a job.
	RetrievalChannel chan J                                                            // The channel a new job gets passed through (once ready)
}

// Helper function for handling watching properly. Returns a function that can be called by multiple goroutines. None of the functions in context can be nil.
func HandleWatching[J any, C any](context WatchContext[J, C], startContext C) func(C, string) {
	var debounceTimer *time.Timer

	// Create a waiting boolean for making sure we're not waiting for rebuilding twice
	waitMutex := &sync.Mutex{}
	waiting := false

	listener := func(ctx C, message string) {
		if debounceTimer != nil {
			debounceTimer.Stop()
		}

		// Create new timer for 500ms
		debounceTimer = time.AfterFunc(500*time.Millisecond, func() {
			waitMutex.Lock()
			defer waitMutex.Unlock()
			if waiting {
				if mconfig.VerboseLogging {
					context.Print("Changes detected, but already trying rebuild.")
				}
				return
			}
			waiting = true

			// Print what the user wants us to say when the change is the one being accepted
			context.Print(message)

			// Wait for the previous job to be cancellable and cancel it
			job, ok := <-context.RetrievalChannel
			if !ok {
				context.Error(fmt.Errorf("couldn't get previous process"))
				return
			}
			if err := context.Stop(job); err != nil && !errors.Is(err, os.ErrProcessDone) {
				context.Error(fmt.Errorf("couldn't kill previous process: %w", err))
				return
			}

			// Start the new job
			if err := context.Start(ctx, &job, context.RetrievalChannel); err != nil {
				context.Error(err)
			}
		})
	}

	// Start the first job
	context.Start(startContext, nil, context.RetrievalChannel)
	return listener
}
