package mconfig

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"unicode"
)

type Plan struct {
	AppName        string                `json:"app_name"`
	Profile        string                `json:"profile"`
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
func (p *PlannedDatabaseType) ContainerName(appName string, profile string) string {
	appName = EverythingToSnakeCase(appName)
	return fmt.Sprintf("mgc-%s-%s-%d", appName, profile, p.Type)
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
	// Pretty-print the plan as JSON
	encoded, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

// Convert back to a plan from printable form
func FromPrintable(printable string) (*Plan, error) {
	plan := &Plan{}
	err := json.Unmarshal([]byte(printable), plan)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

// Get a database by its name. Panics when it can't find the database.
func (p *Plan) Database(name string) PlannedDatabase {
	foundDB := PlannedDatabase{}
	found := false
	for _, t := range p.DatabaseTypes {
		for _, db := range t.Databases {
			if db.ConfigName == name {
				if found {
					log.Fatalln("The database", name, "exists in the config more than once.")
				}
				found = true
				foundDB = db
				foundDB.Port = t.Port
			}
		}
	}
	if !found {
		log.Fatalln("Database", name, "couldn't be found in the plan!")
	}
	return foundDB
}

// Generate a connection string for the database.
func (db PlannedDatabase) ConnectString() string {
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
