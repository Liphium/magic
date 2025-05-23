package mconfig

import "os"

type Environment map[string]EnvironmentValue

// Apply all the environment variables
func (e *Environment) Apply() error {
	for value, key := range *e {
		if err := os.Setenv(value, key.get()); err != nil {
			return err
		}
	}
	return nil
}

type EnvironmentValue struct {
	get func() string
}

// Create a new static environment value.
func ValueStatic(value string) EnvironmentValue {
	return EnvironmentValue{
		get: func() string {
			return value
		},
	}
}
