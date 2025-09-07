package scripting

import (
	"reflect"

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

// baseField should be a text
func createNumberField(index int, kind reflect.Kind, baseField huh.Field) CollectionField {
	return CollectionField{
		Field: baseField,
		Set: func(v reflect.Value) {
			value := baseField.GetValue()
			//  TODO: Implement
		},
	}
}
