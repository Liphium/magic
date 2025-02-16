package database

import (
	"time"

	"github.com/google/uuid"
)

type Preview struct {
	Forge      uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Label      string
	Repository string
	Account    uuid.UUID `gorm:"type:uuid;index"`
	LastViewed time.Time `gorm:"index"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type PreviewSecret struct {
	ID      uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Preview uuid.UUID `gorm:"type:uuid;index"`
	Name    string
	Secret  string

	CreatedAt time.Time
}

type Environment struct {
	ID      uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Account uuid.UUID `gorm:"type:uuid;index"`
	Preview uuid.UUID `gorm:"type:uuid;index"`
	Build   uuid.UUID `gorm:"type:uuid;index"`
	Node    uuid.UUID `gorm:"type:uuid;index"`
	Status  uint      `gorm:"index"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type EnvironmentFile struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Environment uuid.UUID `gorm:"type:uuid;index"`
	Asset       uuid.UUID `gorm:"type:uuid"`
}
