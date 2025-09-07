package scripting

import (
	"errors"
	"fmt"
	"reflect"
	"slices"

	"github.com/charmbracelet/huh"
)

type CollectionField struct {
	Field huh.Field
	Set   func(reflect.Value)
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

		// Build the fields for huh
		collectionFields := []CollectionField{}
		for i := 0; i < genType.NumField(); i++ {
			field := genType.Field(i)

			// Get the prompt for the field
			prompt := field.Tag.Get("prompt")
			if prompt == "" {
				prompt = "Enter value for " + field.Name + ":"
			}

			// Generate the huh field
			var collectionField CollectionField
			if slices.Contains(numberKindsSupported, field.Type.Kind()) {
				baseField := huh.NewText().Title(prompt)
				collectionField = createNumberField(i, field.Type.Kind(), baseField)
			}

			collectionFields = append(collectionFields, collectionField)
		}

		// Run the form

		// Create a new object and return it
		value := reflect.New(genType).Elem()
		for _, field := range collectionFields {
			field.Set(value)
		}

		return value.Interface()
	}
	return collector, nil
}
