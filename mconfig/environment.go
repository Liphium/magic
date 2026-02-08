package mconfig

import (
	"fmt"
	"log"
	"slices"
)

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

func ValueFunction(get func() string) EnvironmentValue {
	return EnvironmentValue{get}
}

// Create a new environment value based on other environment values.
//
// The index in the values array matches the output of the environment value.
func ValueWithBase(values []EnvironmentValue, builder func([]string) string) EnvironmentValue {
	return EnvironmentValue{
		get: func() string {

			// Evaluate all the environment values
			evaluated := make([]string, len(values))
			for i, value := range values {
				evaluated[i] = value.get()
			}

			// Build the output
			return builder(evaluated)
		},
	}
}

// Allocate a new port for the container (and parse it as a environment variable).
func (c *Context) ValuePort(preferredPort uint) EnvironmentValue {

	// Make sure the ports slice is initalized
	if c.ports == nil {
		c.ports = []uint{}
	}

	// Make sure the port isn't already allocated
	if slices.Contains(c.ports, preferredPort) {
		log.Fatalln("Port", preferredPort, "is already taken: taken ports: ", c.ports)
	}
	c.ports = append(c.ports, preferredPort)

	return EnvironmentValue{
		get: func() string {
			allocatedPort, ok := c.Plan().AllocatedPorts[preferredPort]
			if !ok {
				log.Fatalln("Couldn't find port", preferredPort, "in final plan")
			}
			return fmt.Sprintf("%d", allocatedPort)
		},
	}
}
