//go:build !release

package main

import (
	"real-project/starter"

	"github.com/Liphium/magic/v2"
)

func main() {
	magic.Start(starter.BuildMagicConfig())
}
