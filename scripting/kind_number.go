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
		return ValidateNumber(field, s)
	})

	return CollectionField{
		Field: baseField,
		Set: func(v reflect.Value) error {
			obj := baseField.GetValue()
			value, ok := obj.(string)
			if !ok {
				return errors.New("value of text field is not string")
			}

			return ValidateAndSetNumber(field, v, value)
		},
	}
}

// Validate and set a number in a struct field
func ValidateAndSetNumber(field reflect.StructField, structValue reflect.Value, value string) error {
	if err := ValidateNumber(field, value); err != nil {
		return err
	}
	_, err := parseStringToNumber(field.Type.Kind(), value, &structValue)
	return err
}

func ValidateNumber(field reflect.StructField, value string) error {
	if strings.ContainsFunc(value, func(r rune) bool {
		return !unicode.IsDigit(r) && r != '-' && r != '.'
	}) {
		return errors.New("value is not a number")
	}

	// If there is a validate struct tag, use validator to check the value
	tag := field.Tag.Get("validate")
	if tag != "" {
		num, err := parseStringToNumber(field.Type.Kind(), value, nil)
		if err != nil {
			return fmt.Errorf("could not parse value for validation: %w", err)
		}
		return runValidationWithTranslation(field.Name, num, tag)
	}
	return nil
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
		if fieldToSet != nil {
			fieldToSet.SetInt(n)
		}
		return n, nil

	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse number (%s) to float: %s", s, err)
		}
		if fieldToSet != nil {
			fieldToSet.SetFloat(n)
		}
		return n, nil

	default:
		return nil, fmt.Errorf("unsupported kind for parsing: %s", kind)
	}
}
