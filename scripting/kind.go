package scripting

import (
	"fmt"
	"log"
	"reflect"
	"slices"

	"github.com/Liphium/magic/v2/mconfig"
	"github.com/Liphium/magic/v2/util"
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
func CreateCollector(genType reflect.Type) (func([]string) interface{}, error) {

	// Create the validator in case it's not there yet
	if validate == nil {
		validate = createValidator()
	}

	if genType.Kind() != reflect.Struct {
		if mconfig.VerboseLogging {
			util.Log.Println("Ignoring script argument due to not being a struct...")
		}
		return func(s []string) interface{} {
			return "hi magic"
		}, nil
	}

	// Check if all fields are supported
	for i := 0; i < genType.NumField(); i++ {
		field := genType.Field(i)
		if !slices.Contains(SupportedKinds, field.Type.Kind()) && field.Tag.Get("magic") != "ignore" {
			return nil, fmt.Errorf("collecting a %s is currently not supported", field.Type.Kind().String())
		}
	}

	collector := func(arguments []string) interface{} {

		// Collect from command line arguments in case there are some
		if len(arguments) > 0 {
			result, err := useArgumentCollector(genType, arguments)
			if err != nil {
				log.Fatalln("argument collection failed:", err)
			}
			return result
		}

		// Collect from command line form in case there are no arguments
		return useCommandLineCollector(genType)
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
	return err
}
