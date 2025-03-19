package server

// Library represents information about an installed library
type Library struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Status  string `json:"status"`
	Arch    string `json:"architecture"`
}
