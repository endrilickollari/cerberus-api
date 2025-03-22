package server

import "time"

// FileSystemEntry represents a file or directory in the file system
type FileSystemEntry struct {
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	Type         string    `json:"type"`           // "file", "directory", "symlink", etc.
	Size         int64     `json:"size,omitempty"` // Size in bytes, only applicable for files
	Permissions  string    `json:"permissions"`
	Owner        string    `json:"owner"`
	Group        string    `json:"group"`
	LastModified time.Time `json:"last_modified"`
	IsHidden     bool      `json:"is_hidden"`           // Whether the file/directory is hidden (starts with .)
	MimeType     string    `json:"mime_type,omitempty"` // MIME type of the file (only for files)
	Preview      string    `json:"preview,omitempty"`   // Preview of text files (first few lines)
}

// FileSystemListing contains the result of a directory listing
type FileSystemListing struct {
	Path      string            `json:"path"`
	Entries   []FileSystemEntry `json:"entries"`
	Recursive bool              `json:"recursive"`
}
