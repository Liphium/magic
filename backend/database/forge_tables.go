package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// All provider types
const (
	ProviderTypeGitHub = "github"
)

type Forge struct {
	ID             uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Account        uuid.UUID `gorm:"type:uuid;index"`
	Provider       string    // Type of the provider ("github")
	Installation   string
	Repository     string // Identifier of the repository
	RepositoryName string // Short name of the repository
	Label          string
	LastViewed     time.Time `gorm:"index"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// Get the value for a build source with a branch name.
func BuildSourceBranch(name string) string {
	return fmt.Sprintf("branch:%s", name)
}

type Build struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Forge       uuid.UUID `gorm:"type:uuid;index"`
	DisplayName string
	Source      string // Thing being built, for example: branch:main

	CreatedAt time.Time
}

type Asset struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Build        uuid.UUID `gorm:"type:uuid;index"`
	Architecture string
	Path         string

	CreatedAt time.Time
}

type Target struct {
	ID    uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Forge uuid.UUID `gorm:"type:uuid;index"`
	Type  string    `gorm:"index"`
	Value string

	CreatedAt time.Time
}
