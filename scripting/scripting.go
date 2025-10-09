package scripting

import (
	"log"
	"reflect"

	"github.com/Liphium/magic/v2/mrunner"
)

type Script struct {
	Name        string
	Description string
	Collector   func([]string) interface{}
	Handler     func(*mrunner.Runner, interface{}) error
}

// f must be of type func(*mrunner.Runner, T, ...) (..., error).
//
// You can have any return types for f and the second parameter, T, can be a struct that will be dynamically generated
// based on CLI arguments or a CLI form. Otherwise it will be ignored. Read more about it in the documentation.
func CreateScript(name string, description string, f interface{}) Script {

	// Make sure f is a function
	scriptType := reflect.TypeOf(f)
	if scriptType.Kind() != reflect.Func {
		log.Fatalf("No function is provided for script %s.", name)
	}

	// Find the runner parameter position (can be anywhere, or not exist)
	runnerPos := -1
	runnerType := reflect.TypeFor[*mrunner.Runner]()
	for i := 0; i < scriptType.NumIn(); i++ {
		if scriptType.In(i) == runnerType {
			runnerPos = i
			break
		}
	}

	// Find the first parameter that's not the runner (this will be collected)
	collectionPos := -1
	collectionType := reflect.TypeFor[any]()
	for i := 0; i < scriptType.NumIn(); i++ {
		if i != runnerPos {
			collectionPos = i
			collectionType = scriptType.In(i)
			break
		}
	}
	collector, err := CreateCollector(collectionType)
	if err != nil {
		log.Fatalln("Something went wrong with internally: Couldn't create collector:", err)
	}

	// Enforce last return value being an error
	if scriptType.NumOut() < 1 || scriptType.Out(scriptType.NumOut()-1).Name() != reflect.TypeFor[error]().Name() {
		log.Fatalf("Last return type of script %s isn't an error.", name)
	}

	return Script{
		Name:        name,
		Description: description,
		Collector:   collector,
		Handler: func(runner *mrunner.Runner, data interface{}) error {
			// Use reflection to call the function f
			scriptValue := reflect.ValueOf(f)

			// Prepare arguments for all parameters
			args := make([]reflect.Value, scriptType.NumIn())

			// Set the runner parameter if it exists
			if runnerPos != -1 {
				args[runnerPos] = reflect.ValueOf(runner)
			}

			// Set the collected parameter if it exists
			if collectionPos != -1 {
				dataValue := reflect.ValueOf(data)
				if dataValue.Type().ConvertibleTo(collectionType) {
					args[collectionPos] = dataValue.Convert(collectionType)
				} else {
					args[collectionPos] = dataValue
				}
			}

			// Fill remaining parameters with zero values
			for i := 0; i < scriptType.NumIn(); i++ {
				if i != runnerPos && i != collectionPos {
					args[i] = reflect.Zero(scriptType.In(i))
				}
			}

			// Call the function
			results := scriptValue.Call(args)

			// Check if the function returns an error (last return value should be error)
			if len(results) > 0 {
				lastResult := results[len(results)-1]
				if lastResult.Type().Implements(reflect.TypeFor[error]()) {
					if !lastResult.IsNil() {
						return lastResult.Interface().(error)
					}
				}
			}

			return nil
		},
	}
}
