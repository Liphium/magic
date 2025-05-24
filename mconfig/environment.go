package mconfig

type Environment map[string]EnvironmentValue

// Apply all the environment variables
func (e *Environment) Generate() map[string]string {
	copy := map[string]string{}
	for value, key := range *e {
		copy[value] = key.get()
	}
	return copy
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
