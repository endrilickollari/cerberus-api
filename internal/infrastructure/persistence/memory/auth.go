package memory

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/crypto/ssh"
	"remote-server-api/internal/domain/auth"
	sshClient "remote-server-api/internal/infrastructure/ssh"
)

// SessionRepository manages SSH sessions in memory
type SessionRepository struct {
	sessions map[string]*auth.Session
	mu       sync.RWMutex
}

// NewSessionRepository creates a new in-memory session repository
func NewSessionRepository() *SessionRepository {
	return &SessionRepository{
		sessions: make(map[string]*auth.Session),
	}
}

// StoreSession stores a new SSH session
func (r *SessionRepository) StoreSession(ctx context.Context, sessionID string, username string, client *ssh.Client) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.sessions[sessionID] = &auth.Session{
		ID:       sessionID,
		Username: username,
		Client:   client,
	}

	return nil
}

// GetSession retrieves an SSH session by ID
func (r *SessionRepository) GetSession(ctx context.Context, sessionID string) (*auth.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	session, exists := r.sessions[sessionID]
	if !exists {
		return nil, errors.New("session not found")
	}

	return session, nil
}

// RemoveSession removes an SSH session by ID
func (r *SessionRepository) RemoveSession(ctx context.Context, sessionID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.sessions[sessionID]; !exists {
		return errors.New("session not found")
	}

	delete(r.sessions, sessionID)
	return nil
}

// RunCommand executes a command on an SSH session
func (r *SessionRepository) RunCommand(ctx context.Context, sessionID string, command string) (string, error) {
	r.mu.RLock()
	session, exists := r.sessions[sessionID]
	r.mu.RUnlock()

	if !exists {
		return "", errors.New("session not found")
	}

	return sshClient.RunCommand(session.Client, command)
}
