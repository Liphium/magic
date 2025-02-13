package spells

import (
	"errors"
	"strings"
)

// Errors in secret processing
var (
	ErrSecretsNotSet = errors.New("secrets aren't set, please call SetSecrets() to set them")
)

// A secret used by Magic to hide credentials and other stuff from the public (typically something like %jwt_secret%)
type Secret struct {
	Name     string // Name of the secret
	Variable string // Environment variable where the secret is used
	value    string // Value of the secret (not known at parse time)
}

// Set the value of the secret
func (s *Secret) SetValue(value string) {
	s.value = value
}

// Extract all secrets used (they will also be cached, you can use this as often as you want)
func (s *Spell) Secrets() []Secret {

	// Add all the secrets if they aren't there yet
	if s.secretsCached == nil {
		s.secretsCached = []Secret{}

		// Search all environment variables for secrets
		for k, v := range s.Environment {
			split := strings.Split(v, "%")

			// Find all of the secrets in the value of the environment variable
			for i, f := range split {
				if strings.TrimSpace(f) == f && i > 0 && i < len(split)-1 {
					s.secretsCached = append(s.secretsCached, Secret{
						Name:     f,
						Variable: k,
					})
				}
			}
		}
	}

	return s.secretsCached
}

// Replace the secrets in the cache with these new ones
func (s *Spell) SetSecrets(secrets []Secret) {
	s.secretsCached = secrets
}

// Build the environment, this will also replace secrets and stuff
func (s *Spell) AddSecretsToEnvironment() error {

	// Make sure there are secrets
	if s.secretsCached == nil {
		return ErrSecretsNotSet
	}

	// Replace all the secrets in the environment variables
	for k, v := range s.Environment {
		for _, sec := range s.secretsCached {
			if sec.Variable == k {
				v = strings.ReplaceAll(v, "%"+sec.Name+"%", sec.value)
			}
		}
		s.Environment[k] = v
	}

	return nil
}
