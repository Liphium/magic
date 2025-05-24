package integration

import (
	"os"
)

func CreateCache() error {
	mDir, err := GetMagicDirectory(3)
	if err != nil {
		return err
	}
	files, err := os.ReadDir(mDir)
	if err != nil {
		return err
	}

	err = os.Chdir(mDir)
	if err != nil {
		return err
	}

	// Find the magic folder
	for _, entry := range files {
		if entry.IsDir() && entry.Name() == "cache" {
			return nil
		}
	}

	if err := os.Mkdir("cache", 0755); err != nil {
		return err
	}

	return nil
}
