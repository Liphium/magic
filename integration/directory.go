package integration

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Get the magic directory
func GetMagicDirectory(amount int) (string, error) {
	if amount <= 0 {
		return "", errors.New("amount can't be 0 or less")
	}

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
			if entry.IsDir() && entry.Name() == "magic" {
				return filepath.Join(wd, "magic"), nil
			}
		}
		wd = filepath.Dir(wd)
	}
	return "", errors.New("can't find magic directory")
}

func DoesDirExist(dirPath string) (bool, error){
	_, err := os.Stat(filepath.Dir(dirPath));
	if err != nil{
		return false, fmt.Errorf("path to dir does not exist: %w", err)
	} else {
		s, err := os.Stat(dirPath);
		if err != nil{
			return true, nil
		} else if !s.IsDir(){
			return false, errors.New("path leads to an existing file not a dir")
		} else{
			return false, nil
		}
	}
}

func CreateDirIfNotExist(path string, dir string) error {
	if _, err := os.ReadDir(path); err != nil {
		return errors.New("path doesn't exist")
	}
	_, err := os.ReadDir(filepath.Join(path, dir))
	if err == nil {
		if err = os.Chdir(path); err != nil {
			return err
		}
		if err := os.Mkdir(dir, 0755); err != nil {
			return err
		}
		return nil
	}
	return errors.New("directory already exists")

}

func PrintCurrentDirAll() {
	wd, _ := os.Getwd()
	fmt.Println(wd)
	files, _ := os.ReadDir(".")

	// Find the magic folder
	for _, entry := range files {
		fmt.Println(entry.Name())
	}

}
