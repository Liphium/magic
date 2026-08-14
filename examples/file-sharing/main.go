package main

import (
	"file-sharing/app"

	"github.com/Liphium/magic/v4"
)

func main() {
	magic.Start(app.GetConfig())
}
