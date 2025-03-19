package docker

import (
	"context"
	"errors"
	"strings"
)

// Common errors
var (
	ErrSessionNotFound = errors.New("session not found or expired")
	ErrCommandFailed   = errors.New("command execution failed")
)

// SessionRepository defines methods to access SSH sessions
type SessionRepository interface {
	// RunCommand executes a command on the SSH session
	RunCommand(ctx context.Context, sessionID string, command string) (string, error)
}

// Service defines the Docker service
type Service interface {
	// GetContainers retrieves information about Docker containers
	GetContainers(ctx context.Context, sessionID string) ([]Container, error)
}

type service struct {
	sessionRepo SessionRepository
}

// NewService creates a new Docker service
func NewService(sessionRepo SessionRepository) Service {
	return &service{
		sessionRepo: sessionRepo,
	}
}

// GetContainers implements the Service interface
func (s *service) GetContainers(ctx context.Context, sessionID string) ([]Container, error) {
	// Execute command to get Docker containers
	dockerOutput, err := s.sessionRepo.RunCommand(ctx, sessionID, "docker ps")
	if err != nil {
		return nil, err
	}

	// Parse Docker containers
	return parseDockerContainers(dockerOutput), nil
}

// parseDockerContainers parses the output of 'docker ps'
func parseDockerContainers(dockerOutput string) []Container {
	var containers []Container
	lines := strings.Split(dockerOutput, "\n")

	// Skip the header line
	for _, line := range lines[1:] {
		if line == "" {
			continue
		}

		// Docker ps output can be tricky to parse due to variable spacing
		// Using a more robust approach

		// First, get the container ID (first column)
		fields := strings.Fields(line)
		if len(fields) < 7 {
			continue
		}

		containerID := fields[0]

		// Find indices for known columns to better handle spacing
		idxOfImage := strings.Index(line[len(containerID):], " ") + len(containerID)
		line = strings.TrimSpace(line[idxOfImage:])

		// Get image
		fields = strings.Fields(line)
		image := fields[0]
		line = strings.TrimSpace(line[len(image):])

		// Get command
		cmdStart := strings.Index(line, "\"")
		var command string
		if cmdStart >= 0 {
			cmdEnd := strings.Index(line[cmdStart+1:], "\"")
			if cmdEnd >= 0 {
				command = line[cmdStart+1 : cmdStart+1+cmdEnd]
				line = strings.TrimSpace(line[cmdStart+cmdEnd+2:])
			}
		} else {
			// If no quotes, just take the next field
			fields = strings.Fields(line)
			command = fields[0]
			line = strings.TrimSpace(line[len(command):])
		}

		// Split remaining fields carefully
		fields = strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		// Get status, ports, and name
		createdOn := fields[0] + " " + fields[1]
		status := fields[2] + " " + fields[3]

		// For ports, we need to be careful about empty port mappings
		var ports string
		var names string

		if strings.Contains(fields[4], "->") {
			// There's a port mapping
			ports = fields[4]
			if len(fields) > 5 {
				names = fields[5]
			}
		} else {
			// No port mapping
			ports = ""
			names = fields[4]
		}

		container := Container{
			ContainerID: containerID,
			Image:       image,
			Command:     command,
			CreatedOn:   createdOn,
			Status:      status,
			Ports:       ports,
			Names:       names,
		}

		containers = append(containers, container)
	}

	return containers
}
