package scripting

import (
	"fmt"
	"log"
	"reflect"
	"slices"
	"strings"

	"github.com/charmbracelet/huh"
)

// Collect all fields needed for a struct from the arguments array passed in
func useArgumentCollector(genType reflect.Type, arguments []string) (interface{}, error) {
	numFields := genType.NumField()

	// Check if the number of arguments matches the number of struct fields
	if len(arguments) != numFields {

		// Build error message with field information
		var fieldInfo []string
		for i := 0; i < numFields; i++ {
			field := genType.Field(i)
			fieldInfo = append(fieldInfo, fmt.Sprintf("%d. %s (%s) \n", i+1, field.Name, field.Type.Kind().String()))
		}

		if len(arguments) < numFields {
			return nil, fmt.Errorf("insufficient arguments: expected %d, got %d. \n \nRequired arguments: \n%s",
				numFields, len(arguments), strings.Join(fieldInfo, ""))
		} else {
			return nil, fmt.Errorf("too many arguments: expected %d, got %d. \n \nRequired arguments: \n%s",
				numFields, len(arguments), strings.Join(fieldInfo, ""))
		}
	}

	// Create a new instance of the struct
	value := reflect.New(genType).Elem()

	// Parse and set each argument to the corresponding field
	for i, arg := range arguments {
		field := genType.Field(i)
		fieldValue := value.Field(i)

		// Handle different field types
		if field.Type.Kind() == reflect.String {
			if err := ValidateAndSetString(field, fieldValue, arg); err != nil {
				return nil, fmt.Errorf("invalid argument for %s: %w", field.Name, err)
			}
		} else if slices.Contains(numberKindsSupported, field.Type.Kind()) {
			if err := ValidateAndSetNumber(field, fieldValue, arg); err != nil {
				return nil, fmt.Errorf("invalid argument for %s: %w", field.Name, err)
			}
		} else {
			return nil, fmt.Errorf("unsupported field type %s for field %s", field.Type.Kind().String(), field.Name)
		}
	}

	return value.Interface(), nil
}

// Collect all fields needed for a struct from the command line (using a huh form)
func useCommandLineCollector(genType reflect.Type) interface{} {

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
