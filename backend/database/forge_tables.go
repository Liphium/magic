package database

import (
	"time"

	"github.com/google/uuid"
)

type Forge struct {
	ID         uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Project    uuid.UUID `gorm:"type:uuid;index"`
	Label      string    `gorm:"index"`
	Repository string

	Created time.Time
	Updated time.Time
}

type Build struct {
	ID     uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Forge  uuid.UUID `gorm:"type:uuid;index"`
	Target uuid.UUID `gorm:"type:uuid;index"`
	Value  string

	Created time.Time
}

type Asset struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Build        uuid.UUID `gorm:"type:uuid;index"`
	Architecture string    `gorm:"index"`
	Path         string

	Created time.Time
}

type Target struct {
	ID    uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Forge uuid.UUID `gorm:"type:uuid;index"`
	Type  string    `gorm:"index"`
	Value string

	Created time.Time
}
