package database

import (
	"time"

	"github.com/google/uuid"
)

type Forge struct {
	ID         uuid.UUID
	Project    uuid.UUID // Project the Forge was created in
	Label      string
	Repository string // Link to the repository

	Created time.Time
	Updated time.Time
}

type Build struct {
	ID     uuid.UUID
	Forge  uuid.UUID
	Target uuid.UUID
	Value  string // Could be something like the PR id when on a PR target

	Created time.Time
}

type Asset struct {
	ID           uuid.UUID
	Build        uuid.UUID
	Architecture string
	Path         string // URL in the CDN

	Created time.Time
}

type Target struct {
	ID    uuid.UUID
	Forge uuid.UUID
	Type  string
	Value string

	Created time.Time
}
