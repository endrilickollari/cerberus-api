package docker

// Image represents a Docker image
type Image struct {
	Repository    string `json:"repository"`
	Tag           string `json:"tag"`
	ImageID       string `json:"image_id"`
	Created       string `json:"created"`
	Size          string `json:"size"`
	Digest        string `json:"digest,omitempty"`
	Vulnerability string `json:"vulnerability,omitempty"`
}
