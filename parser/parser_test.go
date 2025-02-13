package parser_test

import (
	"testing"

	"github.com/Liphium/magic/parser"
)

func TestBasicFunc(t *testing.T) {
	doc := `
	version = 2
	name = "hi hi"
	tags = ["some", "tags"]
	`

	if err := parser.Parse([]byte(doc)); err != nil {
		t.Fatal(err)
	}
}
