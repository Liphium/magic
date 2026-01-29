package mservices

import (
	"context"

	"github.com/moby/moby/client"
)

// TODO: Extract postgres to external methods and create this interface based on it
type ServiceDriver interface {
	GetUniqueId() string
	CreateContainer(ctx context.Context, c *client.Client, a ContainerAllocation) (string, error)
	IsHealthy(ctx context.Context, c *client.Client, id string) (bool, error)
	Initialize(ctx context.Context, c *client.Client, id string) error
}

// TODO: Figure out proper environment variable handling

// All things required to create a service container
type ContainerAllocation struct {
	Name string
	Port int
}
