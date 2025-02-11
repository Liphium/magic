package database

import (
	"time"

	"github.com/google/uuid"
)

type Preview struct {
	Forge         uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Configuration uuid.UUID `gorm:"type:uuid;index"`

	Created time.Time
	Updated time.Time
}

type EnvironmentConfiguration struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Preview      uuid.UUID `gorm:"type:uuid;index"`
	BuildCommand string
	StartCommand string

	Created time.Time
}

type ConfigurationVariable struct {
	ID            uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Configuration uuid.UUID `gorm:"type:uuid;index"`
	Name          string    `gorm:"index"`
	Type          string
	Value         string

	Created time.Time
}

type ServiceConfiguration struct {
	ID            uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Configuration uuid.UUID `gorm:"type:uuid;index"`
	Type          string    `gorm:"index"`
	Version       uint
	Mappings      string

	Created time.Time
}

type Environment struct {
	ID            uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Preview       uuid.UUID `gorm:"type:uuid;index"`
	Build         uuid.UUID `gorm:"type:uuid;index"`
	Configuration uuid.UUID `gorm:"type:uuid;index"`
	Node          uuid.UUID `gorm:"type:uuid;index"`
	Status        uint      `gorm:"index"`

	Created time.Time
	Updated time.Time
}

type EnvironmentFile struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Environment uuid.UUID `gorm:"type:uuid;index"`
	Asset       uuid.UUID `gorm:"type:uuid;index"`

	Created time.Time
	Updated time.Time
}
