package scripting

import (
	"log"

	"github.com/Liphium/magic/mrunner"
)

type ScriptFunction[T any] = func(*mrunner.Runner, T) error

type Script struct {
	Name        string
	Description string
	Collector   func([]string) interface{}
	Handler     func(*mrunner.Runner, interface{}) error
}

func CreateScript[T any](name string, description string, f ScriptFunction[T]) Script {
	collector, err := CreateCollector[T]()
	if err != nil {
		log.Fatalln("couldn't create collector:", err)
	}

	return Script{
		Name:        name,
		Description: description,
		Collector:   collector,
		Handler: func(runner *mrunner.Runner, data interface{}) error {
			return f(runner, data.(T))
		},
	}
}
