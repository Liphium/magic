package scripting

import (
	"fmt"
	"log"
)

type ScriptFunctionGeneric[T any] = func(T) error

type Script struct {
	Name      string
	Collector func() interface{}
	Handler   func(interface{}) error
}

func CreateScript[T any](name string, f ScriptFunctionGeneric[T]) Script {
	collector, err := createCollector[T]()
	if err != nil {
		log.Fatalln(err)
	}

	return Script{
		Name:      name,
		Collector: collector,
		Handler: func(data interface{}) error {
			return f(data.(T))
		},
	}
}

type SomeScriptThingy struct {
	Name string `prompt:"Enter a name for the test account." validate:"required"`
}

func someScriptFunc(data SomeScriptThingy) error {
	log.Println("chosen name:", data.Name)
	return nil
}

func main() {
	script := CreateScript("hello", someScriptFunc)
	fmt.Println(script.Collector())
}
