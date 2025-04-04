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

	// GetContainerDetail retrieves detailed information about a specific Docker container
	GetContainerDetail(ctx context.Context, sessionID string, containerID string) (*ContainerDetail, error)

	// GetImages retrieves information about Docker images
	GetImages(ctx context.Context, sessionID string) ([]Image, error)

	// GetImageDetail retrieves detailed information about a specific Docker image
	GetImageDetail(ctx context.Context, sessionID string, imageID string) (*ImageDetail, error)

	// DeleteImage deletes a Docker image
	DeleteImage(ctx context.Context, sessionID string, imageID string, force bool) (*ImageDeleteResponse, error)

	// RunContainer runs a Docker container from an image
	RunContainer(ctx context.Context, sessionID string, request ContainerRunRequest) (*ContainerRunResponse, error)
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
	dockerOutput, err := s.sessionRepo.RunCommand(ctx, sessionID, "docker ps -a")
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

		// Extract container ID (first field)
		fields := strings.Fields(line)
		if len(fields) < 4 { // Need at least ID, image, command, and created
			continue
		}

		containerID := fields[0]
		image := fields[1]

		// Extract command - look for quoted text
		var command string
		var cmdEndIndex int

		// Find start and end of command (it's usually in quotes)
		cmdStartIndex := strings.Index(line, "\"")
		if cmdStartIndex >= 0 {
			// Find the closing quote
			cmdEndIndex = strings.Index(line[cmdStartIndex+1:], "\"")
			if cmdEndIndex >= 0 {
				cmdEndIndex += cmdStartIndex + 1 // Adjust for the offset in the substring
				command = line[cmdStartIndex+1 : cmdEndIndex]
			}
		}

		// If we couldn't find a quoted command, try a different approach
		if command == "" {
			// Skip ID and image fields to get to where command should be
			remainingLine := line
			// Skip container ID
			remainingLine = strings.TrimSpace(remainingLine[len(containerID):])
			// Skip image
			remainingLine = strings.TrimSpace(remainingLine[len(image):])

			// Take the next chunk as command
			cmdFields := strings.Fields(remainingLine)
			if len(cmdFields) > 0 {
				command = cmdFields[0]
				// Remove any quote characters
				command = strings.Trim(command, "\"")
			}
		}

		// Extract created time and status, which can vary in format
		remainingFields := fields[2:] // Skip ID and image
		var createdTime, status string
		var i int

		// Skip the command if we already parsed it
		if cmdStartIndex >= 0 && cmdEndIndex >= 0 {
			// Find which field index corresponds to after the command
			for i = 0; i < len(remainingFields); i++ {
				if strings.Contains(remainingFields[i], "\"") && strings.HasSuffix(remainingFields[i], "\"") {
					// This is the end of the command
					i++
					break
				}
			}
			remainingFields = remainingFields[i:]
		} else if command != "" {
			// Skip the field we identified as the command
			remainingFields = remainingFields[1:]
		}

		// Parse created time and status
		if len(remainingFields) >= 4 {
			// Created time is typically "X time ago" format (e.g., "4 days ago")
			createdTime = remainingFields[0] + " " + remainingFields[1]
			if remainingFields[1] != "ago" && len(remainingFields) > 2 && remainingFields[2] == "ago" {
				createdTime += " " + remainingFields[2]
				// Status typically starts with "Up" or "Exited"
				status = strings.Join(remainingFields[3:], " ")
			} else {
				// If "ago" is not a separate field, then status starts right after "created"
				status = strings.Join(remainingFields[2:], " ")
			}
		}

		// Extract ports and names from status
		var ports, names string
		statusFields := strings.Fields(status)

		// Ports and names are the last fields - extract them
		if len(statusFields) >= 2 {
			// Check if there's a port mapping
			portIndex := -1
			for i, field := range statusFields {
				if strings.Contains(field, "->") || strings.Contains(field, ":") || strings.HasSuffix(field, "/tcp") || strings.HasSuffix(field, "/udp") {
					portIndex = i
					break
				}
			}

			if portIndex >= 0 {
				// We found a port
				ports = statusFields[portIndex]

				// The name is likely the last field
				if portIndex < len(statusFields)-1 {
					names = statusFields[len(statusFields)-1]
				}

				// Recalculate status without ports and names
				status = strings.Join(statusFields[:portIndex], " ")
			} else {
				// No port mapping, the last field is probably the name
				names = statusFields[len(statusFields)-1]

				// Recalculate status without the name
				status = strings.Join(statusFields[:len(statusFields)-1], " ")
			}
		}

		// Create the container object
		container := Container{
			ContainerID: containerID,
			Image:       image,
			Command:     command,
			CreatedOn:   createdTime,
			Status:      status,
			Ports:       ports,
			Names:       names,
		}

		containers = append(containers, container)
	}

	return containers
}
