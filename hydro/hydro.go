package hydro

// This is just a concept for now

// Idea:
// Bring Magic's nice developer experience to a wider audience by building a runtime that can support Magic's featureset for any
// programming language. Use Docker Compose to deploy everything Magic needs, but add extra nice features on top like automatic
// testing environments, easily clearing all databases and migrating databases automatically (postgres 17->18).
// Hydro is the engine that should power that.
//
// How that could work:
// 1. Copy Docker Compose file
//
// Migration phase:
// 2. Hydro compares the file to last deployed compoes file (should be saved locally)
// 3. Hydro verifies that all services stayed the same or migrates in case nessecary (postgres 17->18, using helper config)
// 4. Go to deploy phase
//
// If not there, before deploying, Hydro should put all services into one Docker network, so the services don't disturb the host system.
// Important for being able to run multiple projects at the same time without port collisions.
//
// Deploy phase:
// 2. Replace all services that are apps (from Hydro config) with a proxy service
// 3. Start all of the containers that can be started and don't depend on the app services
// 4. Give the signal to whatever is controlling Hydro that apps can be started and proxy the app's traffic through a tunnel into the Docker
// network making it act like a container without actually being one, make sure all the app's exposed ports are covered, map the ports from
// the host system to whatever they really are inside the container + support UDP,TCP
// 5. Repeat until all apps are started and all dependency conflicts are solved
// 6. Finished, that's Hydro, running Docker compose files in an isolated environment for it to be testable and also great DX

type Config struct {
	ComposeFile     string                  // URL of the compose file
	Apps            []string                // All apps that Hydro should proxy
	Helpers         map[string]HelperConfig // Service -> Helper id
	TestEnvironment Environment             // Environment variables Hydro should set for docker compose (for the test environment)
	ProjectName     string                  // The project name Hydro should use
}

type Environment = map[string]string

// Will definitely contain name="something"
type HelperConfig = map[string]string
