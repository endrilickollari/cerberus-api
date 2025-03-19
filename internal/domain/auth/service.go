package auth

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
)

// Common errors
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrSessionNotFound    = errors.New("session not found or expired")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

// TokenService defines methods for JWT token operations
type TokenService interface {
	GenerateToken(username, sessionID string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

// SSHClient defines methods for SSH operations
type SSHClient interface {
	Connect(ip, username, port, password string) (*ssh.Client, error)
}

// Service defines the authentication service
type Service interface {
	// Login authenticates a user via SSH and returns a token
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)

	// ValidateToken validates a JWT token and returns the claims
	ValidateToken(tokenString string) (*Claims, error)

	// GetSession retrieves an active session by ID
	GetSession(ctx context.Context, sessionID string) (*Session, error)
}

type service struct {
	repo         Repository
	sshClient    SSHClient
	tokenService TokenService
}

// NewService creates a new authentication service
func NewService(repo Repository, sshClient SSHClient, tokenService TokenService) Service {
	return &service{
		repo:         repo,
		sshClient:    sshClient,
		tokenService: tokenService,
	}
}

// Login implements the Service interface
func (s *service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// Connect to SSH server
	client, err := s.sshClient.Connect(req.IP, req.Username, req.Port, req.Password)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidCredentials, err)
	}

	// Generate a unique session ID (in production, use a proper UUID library)
	sessionID := fmt.Sprintf("session_%s_%s", req.Username, req.IP)

	// Store the session
	if err := s.repo.StoreSession(ctx, sessionID, req.Username, client); err != nil {
		return nil, fmt.Errorf("failed to store session: %w", err)
	}

	// Generate token
	token, err := s.tokenService.GenerateToken(req.Username, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResponse{Token: token}, nil
}

// ValidateToken implements the Service interface
func (s *service) ValidateToken(tokenString string) (*Claims, error) {
	return s.tokenService.ValidateToken(tokenString)
}

// GetSession implements the Service interface
func (s *service) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	return s.repo.GetSession(ctx, sessionID)
}
