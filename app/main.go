package main

import (
	"fmt"
	"log"

	"github.com/Liphium/magic/scripting"
)

type SomeScriptThingy struct {
	Name  string `prompt:"Name for the test account." validate:"required"`
	Email string `prompt:"E-Mail address for the test account." validate:"required,email"`
}

func someScriptFunc(data SomeScriptThingy) error {
	log.Println("chosen name:", data.Name)
	return nil
}

func main() {
	script := scripting.CreateScript("hello", "", someScriptFunc)
	fmt.Println(script.Collector())
}
