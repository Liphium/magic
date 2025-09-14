package mconfig

import (
	"fmt"
	"log"
	"maps"
	"os"
	"strings"
)

type Context struct {
	appName     string       // Current app name
	profile     string       // Current profile
	directory   string       // Current working directory
	environment *Environment // Environment for environment variables (can be nil)
	databases   []*Database
	ports       []uint // All ports the user wants to allocate
	plan        **Plan // For later filling in with actual information
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

func DefaultContext(appName string, profile string) *Context {
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalln("couldn't get current working directory")
	}

	plan := &Plan{}
	return &Context{
		directory: workDir,
		appName:   appName,
		profile:   profile,
		databases: []*Database{},
		plan:      &plan,
	}
}
