package scripting

import (
	"log"
)

type ScriptFunctionGeneric[T any] = func(T) error

type Script struct {
	Name      string
	Collector func() interface{}
	Handler   func(interface{}) error
}

func CreateScript[T any](name string, f ScriptFunctionGeneric[T]) Script {
	collector, err := CreateCollector[T]()
	if err != nil {
		log.Fatalln("couldn't create collector:", err)
	}

	return Script{
		Name:      name,
		Collector: collector,
		Handler: func(data interface{}) error {
			return f(data.(T))
		},
	}
}
