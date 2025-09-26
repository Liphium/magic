package scripting

import (
	"errors"
	"reflect"

	"github.com/charmbracelet/huh"
)

func CreateStringField(field reflect.StructField, baseField *huh.Input) CollectionField {

	// If there is a validate struct tag, use it for basic string validation (e.g., min/max length)
	baseField.Validate(func(s string) error {
		return ValidateString(field, s)
	})

	return CollectionField{
		Field: baseField,
		Set: func(v reflect.Value) error {
			obj := baseField.GetValue()
			value, ok := obj.(string)
			if !ok {
				return errors.New("value of text field is not string")
			}
			return ValidateAndSetString(field, v, value)
		},
	}
}

// Validate and set a string in a struct field
func ValidateAndSetString(field reflect.StructField, structValue reflect.Value, value string) error {
	if err := ValidateString(field, value); err != nil {
		return err
	}
	structValue.SetString(value)
	return nil
}

// Validate a string based on the struct tag
func ValidateString(field reflect.StructField, value string) error {
	tag := field.Tag.Get("validate")
	if tag != "" {
		if err := runValidationWithTranslation(field.Name, value, tag); err != nil {
			return err
		}
	}
	return nil
}
