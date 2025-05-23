package integration

import "os"

// Get the magic directory
func GetMagicDirectory(recursive bool) (os.DirEntry, error) {
	files, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	// Find the magic folder
	for _, entry := range files {
		if entry.IsDir() && entry.Name() == ".magic" {
			return entry, nil
		}
	}
	return nil, nil
}
