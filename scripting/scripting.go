package scripting

import (
	"log"
)

type ScriptFunction[T any] = func(T) error

type Script struct {
	Name        string
	Description string
	Collector   func() interface{}
	Handler     func(interface{}) error
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
		Handler: func(data interface{}) error {
			return f(data.(T))
		},
	}
}
