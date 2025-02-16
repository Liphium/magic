package database

import (
	"time"

	"github.com/google/uuid"
)

type Rank struct {
	ID              uint `gorm:"primaryKey"`
	Label           string
	PermissionLevel uint

	CreatedAt time.Time
	UpdatedAt time.Time
}

// Credential types
var (
	CredentialTypeGitHub = "github"
)

type Credential struct {
	ID uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`

	Account uuid.UUID `gorm:"type:uuid;index"` // The account id of the user
	Type    string    `gorm:"index"`           // Something like "GitHub"
	Secret  string    `gorm:"index"`           // Secret like the GitHub user id
	Token   string    // User access token for GitHub

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Account struct {
	ID       uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Username string    `gorm:"index"`
	Email    string    `gorm:"index"`
	Rank     uint

	CreatedAt time.Time
	UpdatedAt time.Time
}
