package database

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID       uuid.UUID
	Username string
	Email    string
	Rank     uint

	Created time.Time
	Updated time.Time
}
