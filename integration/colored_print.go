package integration

import (
	"fmt"

	"github.com/mgutz/ansi"
)

func ColorizeScript(text string, name string, color string) string {
	pColor := ansi.ColorFunc(fmt.Sprint(color, "+h"))

	return pColor(fmt.Sprintf("Script %s: ", name)) + text
	//return pColor(text)
}

func ColorizeTest(text string, name string, color string) string {
	pColor := ansi.ColorFunc(fmt.Sprint(color, "+h"))

	return pColor(fmt.Sprintf("Test %s: ", name)) + text
	//return pColor(text)
}
