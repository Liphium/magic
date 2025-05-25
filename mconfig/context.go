package mconfig

type Context struct {
	config      string       // Name of the current config
	profile     string       // Name of the current profile
	environment *Environment // Environment for environment variables (can be nil)
	databases   []*Database
	plan        **Plan // For later filling in with actual information
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

func DefaultContext(config string, profile string) *Context {
	plan := &Plan{}
	return &Context{
		config:    config,
		profile:   profile,
		databases: []*Database{},
		plan:      &plan,
	}
}
