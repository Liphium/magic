package scripting

import (
	"errors"
	"reflect"

	"github.com/charmbracelet/huh"
)

func CreateStringField(field reflect.StructField, baseField *huh.Input) CollectionField {

	// If there is a validate struct tag, use it for basic string validation (e.g., min/max length)
	baseField.Validate(func(s string) error {
		tag := field.Tag.Get("validate")
		if tag != "" {
			return runValidationWithTranslation(field.Name, s, tag)
		}
		return nil
	})

	return CollectionField{
		Field: baseField,
		Set: func(v reflect.Value) error {
			obj := baseField.GetValue()
			value, ok := obj.(string)
			if !ok {
				return errors.New("value of text field is not string")
			}
			v.SetString(value)
			return nil
		},
	}
}
