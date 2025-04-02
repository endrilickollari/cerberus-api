package docker

// ContainerRunRequest represents a request to run a Docker container
type ContainerRunRequest struct {
	Image       string            `json:"image" validate:"required"` // Docker image to run
	Name        string            `json:"name,omitempty"`            // Optional container name
	Ports       []PortMapping     `json:"ports,omitempty"`           // Port mappings
	Volumes     []VolumeMapping   `json:"volumes,omitempty"`         // Volume mappings
	Environment map[string]string `json:"environment,omitempty"`     // Environment variables
	Detached    bool              `json:"detached" default:"true"`   // Run in detached mode
	Restart     string            `json:"restart,omitempty"`         // Restart policy (no, always, on-failure, unless-stopped)
	Network     string            `json:"network,omitempty"`         // Network to connect to
	Command     []string          `json:"command,omitempty"`         // Command to run (overrides default)
}

// RunPortMapping represents a port mapping for a container
type RunPortMapping struct {
	HostPort      string `json:"host_port"`                        // Port on the host
	ContainerPort string `json:"container_port"`                   // Port in the container
	Protocol      string `json:"protocol,omitempty" default:"tcp"` // Protocol (tcp or udp)
}

// VolumeMapping represents a volume mapping for a container
type VolumeMapping struct {
	HostPath      string `json:"host_path"`                           // Path on the host
	ContainerPath string `json:"container_path"`                      // Path in the container
	ReadOnly      bool   `json:"read_only,omitempty" default:"false"` // Mount as read-only
}

// ContainerRunResponse represents a response from running a container
type ContainerRunResponse struct {
	ContainerID string   `json:"container_id"`       // ID of the created container
	Name        string   `json:"name"`               // Name of the container
	Status      string   `json:"status"`             // Container status (e.g., "created", "running")
	Warnings    []string `json:"warnings,omitempty"` // Any warnings generated during container creation
}
