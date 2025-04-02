package docker

import (
	"context"
	"fmt"
	"strings"
)

// DeleteImage implements the Service interface
func (s *service) DeleteImage(ctx context.Context, sessionID string, imageID string, force bool) (*ImageDeleteResponse, error) {
	// Sanitize image ID to prevent command injection
	sanitizedImageID := sanitizeImageID(imageID)

	// Build the docker rmi command
	deleteCmd := "docker rmi"
	if force {
		deleteCmd += " -f" // Force removal (will remove even if image has containers using it)
	}
	deleteCmd += " " + sanitizedImageID

	// Execute the command
	output, err := s.sessionRepo.RunCommand(ctx, sessionID, deleteCmd)

	// Parse the response
	response := &ImageDeleteResponse{}

	// If there was an error, check if it's a specific Docker error we can handle
	if err != nil {
		// If the error is because the image has dependent containers, we can extract that info
		if strings.Contains(err.Error(), "image is being used by running container") {
			response.Errors = append(response.Errors, "Cannot delete image that is being used by running containers")
			return response, nil
		}

		// If it's a conflict error, we can also extract that info
		if strings.Contains(err.Error(), "conflict:") {
			lines := strings.Split(err.Error(), "\n")
			for _, line := range lines {
				if strings.Contains(line, "conflict:") {
					response.Errors = append(response.Errors, strings.TrimSpace(line))
				}
			}
			return response, nil
		}

		// For other errors, just return the error
		return nil, fmt.Errorf("failed to delete image: %w", err)
	}

	// Parse the output lines
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Look for "Deleted:" lines
		if strings.HasPrefix(line, "Deleted:") {
			deletedLayer := strings.TrimSpace(strings.TrimPrefix(line, "Deleted:"))
			response.Deleted = append(response.Deleted, deletedLayer)
			continue
		}

		// Look for "Untagged:" lines
		if strings.HasPrefix(line, "Untagged:") {
			untaggedImage := strings.TrimSpace(strings.TrimPrefix(line, "Untagged:"))
			response.Untagged = append(response.Untagged, untaggedImage)
			continue
		}
	}

	return response, nil
}
