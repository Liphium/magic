package mconfig

import (
	"fmt"
	"os"
	"strings"
)

type Context struct {
	config      string       // Name of the current config
	profile     string       // Name of the current profile
	module      string       // Name of the current module
	magicDir    string       // Current Magic directory
	environment *Environment // Environment for environment variables (can be nil)
	databases   []*Database
	plan        **Plan // For later filling in with actual information
}

func (c *Context) Module() string {
	return c.module
}

func (c *Context) Config() string {
	return c.config
}

func (c *Context) Profile() string {
	return c.profile
}

func (c *Context) Environment() *Environment {
	return c.environment
}

// Set the environment.
func (c *Context) WithEnvironment(env *Environment) {
	c.environment = env
}

// Note: In case you use a relative path, expect it to start in the Magic directory.
func (c *Context) LoadSecretsToEnvironment(path string) error {
	oldWd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("couldn't get working directory: %s", err)
	}

	// Change to magic directory
	if err := os.Chdir(c.magicDir); err != nil {
		return fmt.Errorf("couldn't change to magic directory: %s", err)
	}

	// Add an environment in case there isn't one
	if c.environment == nil {
		c.WithEnvironment(&Environment{})
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

	// Change back into the old working directory
	if err := os.Chdir(oldWd); err != nil {
		return fmt.Errorf("couldn't change back to old working directory: %s", err)
	}
	return nil
}

// Get the databases.
func (c *Context) Databases() []*Database {
	return c.databases
}

// Plan for later (DO NOT EXPECT THIS TO BE FILLED BEFORE DEPLOYMENT STEP)
func (c *Context) Plan() *Plan {
	return *c.plan
}

// Apply a plan for the environment in the config
func (c *Context) ApplyPlan(plan *Plan) {
	*c.plan = plan
}

// Add a new database.
func (c *Context) AddDatabase(database *Database) {
	c.databases = append(c.databases, database)
}

func DefaultContext(module string, config string, profile string, magicDir string) *Context {
	plan := &Plan{}
	return &Context{
		module:    module,
		config:    config,
		profile:   profile,
		magicDir:  magicDir,
		databases: []*Database{},
		plan:      &plan,
	}
}
