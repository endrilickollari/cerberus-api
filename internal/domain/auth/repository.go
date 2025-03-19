package auth

import (
	"context"
	"golang.org/x/crypto/ssh"
)

// Repository defines the interface for session persistence
type Repository interface {
	// StoreSession stores a new SSH session
	StoreSession(ctx context.Context, sessionID string, username string, client *ssh.Client) error

	// GetSession retrieves an SSH session by ID
	GetSession(ctx context.Context, sessionID string) (*Session, error)

	// RemoveSession removes an SSH session by ID
	RemoveSession(ctx context.Context, sessionID string) error
}
