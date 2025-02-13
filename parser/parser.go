package parser

import "github.com/pelletier/go-toml/v2"

type MagicImage struct {
	Version int
	Name    string
	Tags    []string
}

func Parse(image []byte) {
	var mi MagicImage
	toml.Unmarshal(image, &mi)
}
