package mdatabase

// TODO: Extract postgres to external methods and create this interface based on it
type DatabaseDriver interface {
	GetUniqueId() string
}
