package docker

import "time"

// ImageDetail represents detailed information about a Docker image
type ImageDetail struct {
	ID            string              `json:"id"`
	RepoTags      []string            `json:"repo_tags"`
	RepoDigests   []string            `json:"repo_digests"`
	Created       time.Time           `json:"created"`
	Size          int64               `json:"size"`
	VirtualSize   int64               `json:"virtual_size"`
	SharedSize    int64               `json:"shared_size"`
	Architecture  string              `json:"architecture"`
	OS            string              `json:"os"`
	Author        string              `json:"author,omitempty"`
	Container     string              `json:"container,omitempty"`
	DockerVersion string              `json:"docker_version,omitempty"`
	Labels        map[string]string   `json:"labels,omitempty"`
	Env           []string            `json:"env,omitempty"`
	Cmd           []string            `json:"cmd,omitempty"`
	Entrypoint    []string            `json:"entrypoint,omitempty"`
	WorkingDir    string              `json:"working_dir,omitempty"`
	Volumes       map[string]struct{} `json:"volumes,omitempty"`
	ExposedPorts  map[string]struct{} `json:"exposed_ports,omitempty"`
	Layers        []string            `json:"layers,omitempty"`
	History       []ImageHistory      `json:"history,omitempty"`
}

// ImageHistory represents a layer in a Docker image
type ImageHistory struct {
	Created    time.Time `json:"created"`
	CreatedBy  string    `json:"created_by"`
	EmptyLayer bool      `json:"empty_layer"`
	Comment    string    `json:"comment,omitempty"`
}
