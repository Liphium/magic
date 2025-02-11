package database

import (
	"time"

	"github.com/google/uuid"
)

// The node is currently offline
const NodeStatusOffline = 0

// Something happened on the node, an admin should check the situation
const NodeStatusError = 1

// The node is online and accepting new deployments
const NodeStatusOnline = 2

type Node struct {
	ID     uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Domain string    `gorm:"uniqueIndex"`
	Status uint      `gorm:"index"`

	Created time.Time
	Updated time.Time
}
