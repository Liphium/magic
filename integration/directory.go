package integration

import (
	"errors"
	"os"
	"path/filepath"
)

// Get the magic directory
func GetMagicDirectory(amount int) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}


	for i := 0; i < amount; i++ {

		files, err := os.ReadDir(wd)
		if err != nil {
			return "", err
		}

		// Find the magic folder
		for _, entry := range files {
			if entry.IsDir() && entry.Name() == ".magic" {
				return filepath.Join(wd, "magic"), nil
			}
		}
		wd = filepath.Dir(wd)
	}
	return "", errors.New("cant find .magic directory")
}

func CreateDirIfNotExist(path string, dir string) error {
	if _,err := os.ReadDir(path); err != nil{
		return errors.New("path doesn't exist")
	}
	_, err := os.ReadDir(filepath.Join(path, dir))
	if err == nil {
		if err = os.Chdir(path); err != nil{
			return err
		}
		if err := os.Mkdir(dir, 0755); err != nil {
			return err
		}
		return nil
	}
	return errors.New("directory already exists")
	
}
