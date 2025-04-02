package docker

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

// RunContainer implements the Service interface
func (s *service) RunContainer(ctx context.Context, sessionID string, request ContainerRunRequest) (*ContainerRunResponse, error) {
	// Sanitize and validate input
	if err := validateContainerRunRequest(request); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Sanitize image name to prevent command injection
	sanitizedImage := sanitizeImageID(request.Image)

	// Build the docker run command
	cmd := buildDockerRunCommand(request, sanitizedImage)

	// Execute the command
	output, err := s.sessionRepo.RunCommand(ctx, sessionID, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run container: %w", err)
	}

	// Parse the output to get the container ID
	containerID := strings.TrimSpace(output)

	// If we have a container ID, get container status
	if containerID == "" {
		return nil, fmt.Errorf("failed to get container ID from output")
	}

	// Create response
	response := &ContainerRunResponse{
		ContainerID: containerID,
		Name:        request.Name,
		Status:      "created", // Default status
	}

	// If container was run in detached mode, we can get its status
	if request.Detached {
		// Get container status using docker inspect
		statusCmd := fmt.Sprintf("docker inspect --format='{{.State.Status}}' %s", containerID)
		statusOutput, err := s.sessionRepo.RunCommand(ctx, sessionID, statusCmd)
		if err == nil {
			response.Status = strings.TrimSpace(statusOutput)
		}
	}

	return response, nil
}

// validateContainerRunRequest validates a container run request
func validateContainerRunRequest(request ContainerRunRequest) error {
	// Image is required
	if request.Image == "" {
		return fmt.Errorf("image is required")
	}

	// Validate restart policy if provided
	if request.Restart != "" {
		validPolicies := map[string]bool{
			"no": true, "always": true, "on-failure": true, "unless-stopped": true,
		}
		if !validPolicies[request.Restart] {
			return fmt.Errorf("invalid restart policy: %s", request.Restart)
		}
	}

	// Validate port mappings
	for _, port := range request.Ports {
		if port.ContainerPort == "" {
			return fmt.Errorf("container port is required for port mapping")
		}

		// Validate protocol if provided
		if port.Protocol != "" && port.Protocol != "tcp" && port.Protocol != "udp" {
			return fmt.Errorf("invalid port protocol: %s", port.Protocol)
		}
	}

	// Validate volume mappings
	for _, volume := range request.Volumes {
		if volume.HostPath == "" || volume.ContainerPath == "" {
			return fmt.Errorf("both host path and container path are required for volume mapping")
		}
	}

	return nil
}

// sanitizeContainerName sanitizes a container name to prevent command injection
func sanitizeContainerName(name string) string {
	// Docker container names must match regex: [a-zA-Z0-9][a-zA-Z0-9_.-]+
	re := regexp.MustCompile(`[^a-zA-Z0-9_.-]`)
	sanitized := re.ReplaceAllString(name, "")

	// Ensure name starts with [a-zA-Z0-9]
	if len(sanitized) > 0 && !regexp.MustCompile(`^[a-zA-Z0-9]`).MatchString(sanitized) {
		sanitized = "c" + sanitized
	}

	// Limit to reasonable length
	if len(sanitized) > 64 {
		sanitized = sanitized[:64]
	}

	return sanitized
}

// sanitizeEnvironmentVariable sanitizes an environment variable to prevent command injection
func sanitizeEnvironmentVariable(key, value string) (string, string) {
	// Key should only contain alphanumeric characters and underscores
	keyRe := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	sanitizedKey := keyRe.ReplaceAllString(key, "_")

	// Value should be properly escaped
	sanitizedValue := strings.ReplaceAll(value, "'", "'\\''")

	return sanitizedKey, sanitizedValue
}

// buildDockerRunCommand builds the docker run command based on the request
func buildDockerRunCommand(request ContainerRunRequest, sanitizedImage string) string {
	cmd := "docker run"

	// Add detached mode if requested
	if request.Detached {
		cmd += " -d"
	}

	// Add container name if provided
	if request.Name != "" {
		sanitizedName := sanitizeContainerName(request.Name)
		cmd += fmt.Sprintf(" --name '%s'", sanitizedName)
	}

	// Add restart policy if provided
	if request.Restart != "" {
		cmd += fmt.Sprintf(" --restart=%s", request.Restart)
	}

	// Add network if provided
	if request.Network != "" {
		// Sanitize network name to prevent command injection
		sanitizedNetwork := regexp.MustCompile(`[^a-zA-Z0-9_.-]`).ReplaceAllString(request.Network, "")
		cmd += fmt.Sprintf(" --network=%s", sanitizedNetwork)
	}

	// Add port mappings
	for _, port := range request.Ports {
		portSpec := ""
		if port.HostPort != "" {
			portSpec = fmt.Sprintf("%s:", port.HostPort)
		}
		portSpec += port.ContainerPort

		// Add protocol if provided and not tcp (which is the default)
		if port.Protocol != "" && port.Protocol != "tcp" {
			portSpec += fmt.Sprintf("/%s", port.Protocol)
		}

		cmd += fmt.Sprintf(" -p %s", portSpec)
	}

	// Add volume mappings
	for _, volume := range request.Volumes {
		volumeSpec := fmt.Sprintf("%s:%s", volume.HostPath, volume.ContainerPath)
		if volume.ReadOnly {
			volumeSpec += ":ro"
		}
		cmd += fmt.Sprintf(" -v %s", volumeSpec)
	}

	// Add environment variables
	for key, value := range request.Environment {
		sanitizedKey, sanitizedValue := sanitizeEnvironmentVariable(key, value)
		cmd += fmt.Sprintf(" -e %s='%s'", sanitizedKey, sanitizedValue)
	}

	// Add image
	cmd += fmt.Sprintf(" %s", sanitizedImage)

	// Add command if provided
	if len(request.Command) > 0 {
		for _, arg := range request.Command {
			// Properly escape command arguments
			escapedArg := strings.ReplaceAll(arg, "'", "'\\''")
			cmd += fmt.Sprintf(" '%s'", escapedArg)
		}
	}

	return cmd
}
