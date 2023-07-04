package monitoring

import "github.com/NethermindEth/egn/internal/monitoring/services/types"

// ServiceAPI defines the interface for managing a monitoring service. It provides methods for
// adding and removing targets, retrieving environment variables, setting up the service, and initializing the service.
type ServiceAPI interface {
	// AddTarget adds a new target to the service's configuration given the endpoint of the new node.
	// The instanceID is used to identify the node in the service's configuration.
	AddTarget(endpoint, instanceID string) error

	// RemoveTarget removes a target from the service's configuration given the endpoint of the node to be removed.
	RemoveTarget(endpoint string) error

	// DotEnv returns a map of the service's environment variables and their default values.
	DotEnv() map[string]string

	// Setup configures the service given a map of options. The options should include the values for the environment variables.
	Setup(options map[string]string) error

	// Init initializes the service with the given ServiceOptions.
	Init(types.ServiceOptions) error

	// SetContainerIP sets the container IP of the service.
	SetContainerIP(ip, containerName string)

	// ContainerName returns the name of the service's container.
	ContainerName() string
}
