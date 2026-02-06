package mservices

import (
	"context"

	"github.com/moby/moby/client"
)

// An instruction to do something with a container.
//
// This is used by Magic to for example tell database providers to clear their databases.
type Instruction string

const (
	InstructionDropTables  Instruction = "database:drop_tables"
	InstructionClearTables Instruction = "database:clear_tables"
)

// A service driver is a manager for containers running a particular service image.
//
// That can be databases or literally anything you could imagine. It provides a unified interface for Magic to be able to properly control those Docker containers.
type ServiceDriver interface {
	GetUniqueId() string

	// Should return the amount of ports required to start the container.
	GetRequiredPortAmount() int

	// Should return the image. Magic will
	GetImage() string

	// Create a new container for this type of service
	CreateContainer(ctx context.Context, c *client.Client, a ContainerAllocation) (string, error)

	// This method should check if the container with the id is healthy for this service
	IsHealthy(ctx context.Context, c *client.Client, container ContainerInformation) (bool, error)

	// Called to initialize the container when it is healthy
	Initialize(ctx context.Context, c *client.Client, container ContainerInformation) error

	// An instruction sent down from Magic to potentially do something with the container (not every service has to handle every instruction).
	//
	// When implementing, please look into the instructions you can support.
	HandleInstruction(ctx context.Context, c *client.Client, container ContainerInformation, instruction Instruction) error
}

// All things required to create a service container
type ContainerAllocation struct {
	Name  string
	Ports []int
}

type ContainerInformation struct {
	ID    string
	Name  string
	Ports []int
}
