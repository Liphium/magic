package database

import (
	"time"

	"github.com/google/uuid"
)

type Preview struct {
	Forge         uuid.UUID // Forge the preview belongs to (only one possible per Forge)
	Configuration uuid.UUID // Configuration of the preview that will be used for all environments spun up

	Created time.Time
	Updated time.Time
}

type EnvironmentConfiguration struct {
	ID           uuid.UUID
	Preview      uuid.UUID
	BuildCommand string
	StartCommand string

	Created time.Time
}

type ConfigurationVariable struct {
	ID            uuid.UUID
	Configuration uuid.UUID
	Name          string
	Type          string
	Value         string

	Created time.Time
}

type ServiceConfiguration struct {
	ID            uuid.UUID
	Configuration uuid.UUID
	Type          string
	Version       uint
	Mappings      string

	Created time.Time
}

type Environment struct {
	ID            uuid.UUID
	Preview       uuid.UUID // Forge ID because the Preview also belongs to the Forge
	Build         uuid.UUID
	Configuration uuid.UUID // Configuration used to build the environment
	Node          uuid.UUID
	Status        uint

	Created time.Time
	Updated time.Time
}

type EnvironmentFile struct {
	ID          uuid.UUID
	Environment uuid.UUID
	Asset       uuid.UUID

	Created time.Time
	Updated time.Time
}
