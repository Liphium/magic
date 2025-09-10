package scripting

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"slices"

	"github.com/charmbracelet/huh"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type CollectionField struct {
	Index int
	Field huh.Field
	Set   func(reflect.Value) error
}

var SupportedKinds = append([]reflect.Kind{
	reflect.String,
}, numberKindsSupported...)

// Create a function that can take a struct
func CreateCollector[T any]() (func() interface{}, error) {

	// Create the validator in case it's not there yet
	if validate == nil {
		validate = createValidator()
	}

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
		fields := []huh.Field{}
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
			baseField := huh.NewInput().Title(prompt)
			if slices.Contains(numberKindsSupported, field.Type.Kind()) {
				collectionField = CreateNumberField(field, baseField)
			} else if field.Type.Kind() == reflect.String {
				collectionField = CreateStringField(field, baseField)
			}
			collectionField.Index = i
			fields = append(fields, collectionField.Field)
			collectionFields = append(collectionFields, collectionField)
		}

		// Run the form
		form := huh.NewForm(huh.NewGroup(fields...))
		if err := form.Run(); err != nil {
			log.Fatalln("couldn't get data for script:", err)
		}

		// Create a new object and return it
		value := reflect.New(genType).Elem()
		for _, field := range collectionFields {
			if err := field.Set(value.Field(field.Index)); err != nil {
				log.Fatalf("couldn't set field %d of %s: %s", field.Index, value.Type().Name(), err)
			}
		}

		return value.Interface()
	}
	return collector, nil
}

func createValidator() *validator.Validate {
	v := validator.New()
	return v
}

// (Sort of) translate the error and valdiate using the validator package
func runValidationWithTranslation(field string, value any, rules string) error {
	err := validate.Var(value, rules)
	if vErr, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range vErr {
			return fmt.Errorf("%s is not valid: %s", field, fieldError.Tag())
		}
	}
	return fmt.Errorf("validation failed: %s", err)
}
