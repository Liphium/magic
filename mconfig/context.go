package mconfig

import (
	"fmt"
	"maps"
	"os"
	"strings"
)

type Context struct {
	appName     string       // Current app name
	profile     string       // Current profile
	projectDir  string       // Current project directory
	environment *Environment // Environment for environment variables (can be nil)
	services    []ServiceDriver
	ports       []uint // All ports the user wants to allocate
	plan        *Plan  // For later filling in with actual information
}

// The app name you set in your config.
func (c *Context) AppName() string {
	return c.appName
}

// The current profile.
//
// test = Test profile.
// default = Default profile.
// You can set the profile by passing the --m-profile flag to the executable that includes magic.
func (c *Context) Profile() string {
	return c.profile
}

func (c *Context) Environment() *Environment {
	return c.environment
}

func (c *Context) Ports() []uint {
	return c.ports
}

func (c *Context) ProjectDirectory() string {
	return c.projectDir
}

// Set the environment.
func (c *Context) WithEnvironment(env Environment) {
	if c.environment == nil {
		c.environment = &Environment{}
	}

	maps.Copy((*c.environment), env)
}

// Note: In case you use a relative path, expect it to start in the Magic directory.
func (c *Context) LoadSecretsToEnvironment(path string) error {

	// Add an environment in case there isn't one
	if c.environment == nil {
		c.environment = &Environment{}
	}

	// Load all the secrets from the file
	bytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("couldn't read file: %s", err)
	}
	content := string(bytes)
	for l := range strings.Lines(content) {
		if strings.HasPrefix(l, "#") {
			continue
		}
		args := strings.Split(l, "=")
		if len(args) < 2 {
			continue
		}
		key := strings.TrimSpace(args[0])
		value := strings.Trim(strings.TrimSpace(args[1]), "\"'")
		(*c.environment)[key] = ValueStatic(value)
	}
	return nil
}

// Get all services requested.
func (c *Context) Services() []ServiceDriver {
	return c.services
}

// Plan for later (DO NOT EXPECT THIS TO BE FILLED BEFORE DEPLOYMENT STEP)
func (c *Context) Plan() *Plan {
	return c.plan
}

// Register a service driver for a service
func (c *Context) Register(driver ServiceDriver) ServiceDriver {
	c.services = append(c.services, driver)
	return driver
}

func DefaultContext(appName string, profile string, projectDir string) *Context {
	return &Context{
		projectDir: projectDir,
		appName:    appName,
		profile:    profile,
		services:   []ServiceDriver{},
		plan:       &Plan{},
	}
}
