package scripting

import (
	"errors"
	"fmt"
	"reflect"
	"slices"

	"github.com/charmbracelet/huh"
)

type CollectionField[T any] struct {
	Field huh.Field
	Get   func() T
}

var SupportedKinds = append([]reflect.Kind{
	reflect.String,
}, numberKindsSupported...)

// Create a function that can take a struct
func createCollector[T any]() (func() interface{}, error) {
	genType := reflect.TypeFor[T]()
	if genType.Kind() != reflect.Struct {
		return nil, errors.New("generic type isn't a struct")
	}

	// Check if all fields are supported
	for i := 0; i < genType.NumField(); i++ {
		field := genType.Field(i)
		if !slices.Contains(SupportedKinds, field.Type.Kind()) {
			return nil, fmt.Errorf("collecting a %s is currently not supported", field.Type.Kind().String())
		}
	}

	// Create the actual collection function
	collector := func() interface{} {

		// Create a new struct of the type and generate the huh form for it
		value := reflect.New(genType).Elem()
		for i := 0; i < genType.NumField(); i++ {
			field := genType.Field(i)

			// Get the prompt for the field
			prompt := field.Tag.Get("prompt")
			if prompt == "" {
				prompt = "Enter value for " + field.Name + ":"
			}

			// Generate the huh field
			var baseField huh.Field
			if slices.Contains(numberKindsSupported, field.Type.Kind()) {
				baseField = huh.NewText().Title(prompt)
				// TODO: Create number field
			}

			fmt.Print(prompt + " ")
			var input string
			fmt.Scanln(&input)
			value.Field(i).SetString(input)
		}
		return value.Interface()
	}
	return collector, nil
}
