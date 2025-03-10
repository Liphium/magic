package cloud

type ActionType = uint
type ServerType = uint

const (
	// All the differnet server types
	ServerTypeCPU2 ServerType = 0
	ServerTypeCPU4 ServerType = 1
	ServerTypeCPU8 ServerType = 2

	// All the different types of actions
	ActionTypeStart ActionType = 0
)

// Id of a server (passed to the cloud providers for future requests)
type ServerID struct{}

// Status of an action (will be displayed on the frontend)
type ActionStatus struct {
	Name string
}

// An action
type Action struct {
	Type     ActionType
	Status   ActionStatus
	Progress float32 // From 0 to 1
}

type CloudProvider interface {

	// Returns the name of the cloud provider
	name() string

	// Called when the cloud provider is initialized (to load tokens and stuff)
	init() error

	// Start a new server on the cloud provider.
	//
	// Returns a ServerID to contact the server in follow-up requests. Also returns the
	// of the create server action.
	startServer(ServerType) (ServerID, Action, error)

	// List actions
	listActions(ServerID) ([]Action, error)

	// Get the status of an action (used to update the panel in real time).
	getActionStatus(ServerID, Action) (ActionStatus, error)
}
