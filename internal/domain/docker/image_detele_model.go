package docker

// ImageDeleteResponse represents the response from a Docker image deletion
type ImageDeleteResponse struct {
	Deleted  []string `json:"deleted,omitempty"`  // List of deleted image layers
	Untagged []string `json:"untagged,omitempty"` // List of untagged images
	Errors   []string `json:"errors,omitempty"`   // List of errors encountered during deletion
}
