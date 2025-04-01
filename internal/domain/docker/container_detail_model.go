package docker

import "time"

// ContainerDetail represents detailed information about a Docker container
type ContainerDetail struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Image           string            `json:"image"`
	Command         []string          `json:"command"`
	Created         time.Time         `json:"created"`
	Status          string            `json:"status"`
	Ports           map[string]string `json:"ports"`
	NetworkMode     string            `json:"network_mode"`
	RestartPolicy   string            `json:"restart_policy"`
	Platform        string            `json:"platform"`
	Mounts          []Mount           `json:"mounts"`
	Labels          map[string]string `json:"labels"`
	NetworkSettings NetworkSettings   `json:"network_settings"`
	State           ContainerState    `json:"state"`
	HostConfig      HostConfig        `json:"host_config"`
}

// Mount represents a container mount point
type Mount struct {
	Type        string `json:"type"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Mode        string `json:"mode"`
	RW          bool   `json:"rw"`
}

// NetworkSettings represents container network settings
type NetworkSettings struct {
	IPAddress    string        `json:"ip_address"`
	Gateway      string        `json:"gateway"`
	MacAddress   string        `json:"mac_address"`
	NetworkName  string        `json:"network_name"`
	EndpointID   string        `json:"endpoint_id"`
	NetworkID    string        `json:"network_id"`
	SubnetPrefix string        `json:"subnet_prefix"`
	SubnetMask   string        `json:"subnet_mask"`
	PortMappings []PortMapping `json:"port_mappings"`
}

// PortMapping represents a container port mapping
type PortMapping struct {
	ContainerPort string `json:"container_port"`
	HostPort      string `json:"host_port"`
	Protocol      string `json:"protocol"`
}

// ContainerState represents a container's state
type ContainerState struct {
	Status     string    `json:"status"`
	Running    bool      `json:"running"`
	Paused     bool      `json:"paused"`
	Restarting bool      `json:"restarting"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	ExitCode   int       `json:"exit_code"`
	Error      string    `json:"error"`
}

// HostConfig represents container host configuration
type HostConfig struct {
	AutoRemove      bool     `json:"auto_remove"`
	Privileged      bool     `json:"privileged"`
	PublishAllPorts bool     `json:"publish_all_ports"`
	RestartPolicy   string   `json:"restart_policy"`
	DNS             []string `json:"dns"`
	CapAdd          []string `json:"cap_add"`
	CapDrop         []string `json:"cap_drop"`
}
