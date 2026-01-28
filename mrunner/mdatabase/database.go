package mdatabase

// TODO: Extract postgres to external methods and create this interface based on it
type DatabaseDriver interface {
	GetUniqueId() string
}

// TODO: Figure out proper environment variable handling

// All things required to create a database container
type DatabaseContainerAllocation struct {
	Name string
	Port int
}
