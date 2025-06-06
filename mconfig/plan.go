package mconfig

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"unicode"
)

type Plan struct {
	Environment    map[string]string     `json:"environment"`
	DatabaseTypes  []PlannedDatabaseType `json:"database_types"`
	AllocatedPorts map[uint]uint         `json:"ports"`
}

type PlannedDatabaseType struct {
	Port      uint              `json:"port"`
	Type      DatabaseType      `json:"type"`
	Databases []PlannedDatabase `json:"databases"`
}

// Name for the database Docker container
func (p *PlannedDatabaseType) ContainerName(modName string, config string, profile string) string {
	modName = EverythingToSnakeCase(modName)
	return fmt.Sprintf("mgc-%s-%s-%s-%d", modName, config, profile, p.Type)
}

type PlannedDatabase struct {
	ConfigName string `json:"config_name"` // Name in the config
	Name       string `json:"name"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Hostname   string `json:"hostname"`

	// Just for developers to access, not included in actual plan
	Type DatabaseType `json:"-"`
	Port uint         `json:"-"`
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

// Get a database by its name.
func (p *Plan) Database(name string) (PlannedDatabase, error) {
	foundDB := PlannedDatabase{}
	found := false
	for _, t := range p.DatabaseTypes {
		for _, db := range t.Databases {
			if db.ConfigName == name {
				if found {
					return PlannedDatabase{}, errors.New("this database exists more than once")
				}
				found = true
				foundDB = db
				foundDB.Port = t.Port
			}
		}
	}
	if !found {
		return foundDB, errors.New("database not found")
	}
	return foundDB, nil
}

// Generate a connection string for the database.
func (db *PlannedDatabase) ConnectString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", db.Hostname, db.Port, db.Username, db.Password, db.Name)
}

// Convert every character except for letters and digits directly to _
func EverythingToSnakeCase(s string) string {
	newString := ""
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			newString += string(unicode.ToLower(r))
		} else if !strings.HasSuffix(newString, "_") {
			newString += "_"
		}
	}
	return newString
}
