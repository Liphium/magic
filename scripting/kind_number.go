package scripting

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/charmbracelet/huh"
)

var numberKindsSupported = []reflect.Kind{
	reflect.Uint,
	reflect.Uint8,
	reflect.Uint16,
	reflect.Uint32,
	reflect.Uint64,
	reflect.Int,
	reflect.Int8,
	reflect.Int16,
	reflect.Int32,
	reflect.Int64,
	reflect.Float32,
	reflect.Float64,
}

func CreateNumberField(field reflect.StructField, baseField *huh.Input) CollectionField {
	baseField.Validate(func(s string) error {
		if strings.ContainsFunc(s, func(r rune) bool {
			return !unicode.IsDigit(r) && r != '-' && r != '.'
		}) {
			return errors.New("value is not a number")
		}

		// If there is a validate struct tag, use validator to check the value
		tag := field.Tag.Get("validate")
		if tag != "" {
			num, err := parseStringToNumber(field.Type.Kind(), s, nil)
			if err != nil {
				return fmt.Errorf("could not parse value for validation: %w", err)
			}
			return runValidationWithTranslation(field.Name, num, tag)
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

			_, err := parseStringToNumber(field.Type.Kind(), value, &v)
			return err
		},
	}
}

// Converts a string to a number (uint64, int64, float64) based on the reflect.Kind and set it in a struct field in case wanted
func parseStringToNumber(kind reflect.Kind, s string, fieldToSet *reflect.Value) (any, error) {
	switch kind {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		num, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse number (%s) to uint: %s", s, err)
		}
		if fieldToSet != nil {
			fieldToSet.SetUint(num)
		}
		return num, nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse number (%s) to int: %s", s, err)
		}
		return n, nil

	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse number (%s) to float: %s", s, err)
		}
		return n, nil

	default:
		return nil, fmt.Errorf("unsupported kind for parsing: %s", kind)
	}
}
