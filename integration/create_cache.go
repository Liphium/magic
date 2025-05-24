package integration

import (
	"errors"
	"os"
)

func CreateCache() error {
	mDir, err := GetMagicDirectory(3)
	if err != nil {
		return err
	}
	if mDir == nil {
		return errors.New("couldn't find .magic folder")
	}

	if err := os.Mkdir(".cache", 0755); err != nil {
		return err
	}

	return nil
}
