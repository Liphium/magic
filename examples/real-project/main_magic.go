//go:build !release
// +build !release

package main

import (
	"real-project/starter"

	"github.com/Liphium/magic"
)

func main() {
	magic.Start(starter.BuildMagicConfig())
}
