package database

import (
	"time"

	"github.com/google/uuid"
)

type Preview struct {
	Forge         uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Account       uuid.UUID `gorm:"type:uuid;index"`
	Configuration uuid.UUID `gorm:"type:uuid;index"`
	LastViewed    time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

type EnvironmentConfiguration struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Preview      uuid.UUID `gorm:"type:uuid;index"`
	BuildCommand string
	StartCommand string

	CreatedAt time.Time
}

type ConfigurationVariable struct {
	ID            uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Configuration uuid.UUID `gorm:"type:uuid;index"`
	Name          string    `gorm:"index"`
	Type          string
	Value         string
}

type ServiceConfiguration struct {
	ID            uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Configuration uuid.UUID `gorm:"type:uuid;index"`
	Type          string    `gorm:"index"`
	Version       uint
	Mappings      string
}

type Environment struct {
	ID            uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Account       uuid.UUID `gorm:"type:uuid;index"`
	Preview       uuid.UUID `gorm:"type:uuid;index"`
	Build         uuid.UUID `gorm:"type:uuid;index"`
	Configuration uuid.UUID `gorm:"type:uuid;index"`
	Node          uuid.UUID `gorm:"type:uuid;index"`
	Status        uint      `gorm:"index"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type EnvironmentFile struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Environment uuid.UUID `gorm:"type:uuid;index"`
	Asset       uuid.UUID `gorm:"type:uuid;index"`
}
