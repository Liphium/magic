//go:build !release

package main

import (
	"real-project/starter"

	"github.com/Liphium/magic/v4"
)

func main() {
	magic.Start(starter.BuildMagicConfig())
}
