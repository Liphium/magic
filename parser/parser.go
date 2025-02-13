package parser

import (
	"fmt"

	"github.com/pelletier/go-toml/v2"
)

type MagicImage struct {
	Version int
	Name    string
	Tags    []string
}

func Parse(image []byte) error {
	var mi MagicImage
	if err := toml.Unmarshal(image, &mi); err != nil {
		return err
	}

	fmt.Println(mi)
	return nil
}
