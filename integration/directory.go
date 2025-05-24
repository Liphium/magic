package integration

import (
	"os"
	"path/filepath"
)

// Get the magic directory
func GetMagicDirectory(amount int) (os.DirEntry, error) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for i := 0; i < amount; i++ {

		files, err := os.ReadDir(wd)
		if err != nil {
			return nil, err
		}

		// Find the magic folder
		for _, entry := range files {
			if entry.IsDir() && entry.Name() == ".magic" {
				return entry, nil
			}
		}
		wd = filepath.Dir(wd)
	}
	return nil, nil
}
