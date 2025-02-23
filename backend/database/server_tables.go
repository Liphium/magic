package database

import (
	"time"

	"github.com/google/uuid"
)

// The node is currently offline
const WizardStatusOffline = 0

// Something happened on the wizard, an admin should check the situation
const WizardStatusError = 1

// The wizard is online and accepting new jobs
const WizardStatusOnline = 2

type Wizard struct {
	ID     uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Domain string    `gorm:"unique"`
	Token  string    `gorm:"index"`
	Status uint      `gorm:"index"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// All the different job types
const (
	JobTypeBuild   = "build"
	JobTypePreview = "preview"
)

type Job struct {
	ID     uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Type   string    `gorm:"index"` // Type: "build"
	Target string    // Forge id in case of type "build"

	Claimed bool      `gorm:"index"`
	Wizard  uuid.UUID `gorm:"index"`

	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
}
