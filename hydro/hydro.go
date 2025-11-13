package hydro

// This is just a concept for now

type Config struct {
	ComposeFile string                  // URL of the compose file
	Apps        []string                // All apps that Hydro should proxy
	Helpers     map[string]HelperConfig // Service -> Helper id
	Environment Environment             // Environment variables Hydro should set for docker compose
	ProjectName string                  // The project name Hydro should use
}

type Environment = map[string]string

// Will definitely contain name="something"
type HelperConfig = map[string]string
