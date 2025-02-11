package database

import (
	"time"

	"github.com/google/uuid"
)

type Rank struct {
	ID              uint `gorm:"primaryKey"`
	Label           string
	PermissionLevel uint `gorm:"index"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Account struct {
	ID       uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Username string    `gorm:"index"`
	Email    string    `gorm:"uniqueIndex"`
	Rank     uint      `gorm:"index"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Session struct {
	ID      uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Token   string    `gorm:"index"`
	Account uuid.UUID `gorm:"type:uuid;index"`
	LastUse time.Time `gorm:"index"`

	CreatedAt time.Time
}

type Project struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Label       string    `gorm:"index"`
	Description string
	Creator     uuid.UUID `gorm:"type:uuid;index"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
