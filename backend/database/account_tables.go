package database

import (
	"time"

	"github.com/google/uuid"
)

type Rank struct {
	ID              uint
	PermissionLevel uint

	Created time.Time
	Updated time.Time
}

type Account struct {
	ID       uuid.UUID
	Username string
	Email    string
	Rank     uint

	Created time.Time
	Updated time.Time
}

type Session struct {
	ID      uuid.UUID
	Token   string
	Account uuid.UUID
	LastUse time.Time

	Created time.Time
}

type Project struct {
	ID          uuid.UUID
	Label       string
	Description string
	Creator     uuid.UUID // Account that created and owns the project

	Created time.Time
	Updated time.Time
}
