//go:build !release

package main

import (
	"real-project/starter"

	"github.com/Liphium/magic/v3"
)

func main() {
	magic.Start(starter.BuildMagicConfig())
}
