package main

import (
	"file-sharing/app"

	"github.com/Liphium/magic/v3"
)

func main() {
	magic.Start(app.GetConfig())
}
