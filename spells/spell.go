package spells

import "encoding/json"

// All environment variables and their values
type Environment map[string]string

// Property of the database (like password) -> Environment variable
type DatabaseMapping map[string]string

type Database struct {
	Type     string          `json:"type"`
	Name     string          `json:"name"`
	Mappings DatabaseMapping `json:"mappings"`
}

type Spell struct {
	Version     uint        `json:"version"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Environment Environment `json:"environment"`
	Databases   []Database  `json:"databases"`

	// All things for processing
	secretsCached []Secret `json:"-"`
}

// Parse a Spell from bytes
func ParseSpell(encoded []byte) (*Spell, error) {

	// Parse the actual spell to JSON
	var spell Spell
	if err := json.Unmarshal(encoded, &spell); err != nil {
		return &Spell{}, err
	}

	// Return the spell as a pointer
	return &spell, nil
}

// Parse a Spell from a string
func ParseSpellString(encoded string) (*Spell, error) {
	return ParseSpell([]byte(encoded))
}
