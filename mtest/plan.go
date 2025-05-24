package mtest

import (
	"encoding/base64"
	"encoding/json"

	"github.com/Liphium/magic/mconfig"
)

type Plan struct {
	Environment   map[string]string     `json:"environment"`
	DatabaseTypes []PlannedDatabaseType `json:"database_types"`
}

type PlannedDatabaseType struct {
	Port      uint                 `json:"port"`
	Type      mconfig.DatabaseType `json:"type"`
	Databases []PlannedDatabase    `json:"databases"`
}

type PlannedDatabase struct {
	ConfigName string `json:"config_name"` // Name in the config
	Name       string `json:"name"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Hostname   string `json:"hostname"`
}

// Turn the plan into printable form
func (p *Plan) ToPrintable() (string, error) {
	encoded, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encoded), nil
}

// Convert back to a plan from printable form
func FromPrintable(printable string) (*Plan, error) {
	decoded, err := base64.StdEncoding.DecodeString(printable)
	if err != nil {
		return nil, err
	}
	plan := &Plan{}
	err = json.Unmarshal(decoded, plan)
	if err != nil {
		return nil, err
	}
	return plan, nil
}
