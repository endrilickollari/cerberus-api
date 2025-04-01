package docker

import (
	"context"
	"strings"
)

// GetImages implements the Service interface
func (s *service) GetImages(ctx context.Context, sessionID string) ([]Image, error) {
	// Execute command to get Docker images
	// We use the -a flag to show all images, including intermediate images
	// Format: repository tag image_id created size
	imagesOutput, err := s.sessionRepo.RunCommand(ctx, sessionID, "docker images --format \"{{.Repository}}|{{.Tag}}|{{.ID}}|{{.CreatedSince}}|{{.Size}}|{{.Digest}}\"")
	if err != nil {
		return nil, err
	}

	// Parse Docker images
	return parseDockerImages(imagesOutput), nil
}

// parseDockerImages parses the output of 'docker images' command with custom format
func parseDockerImages(imagesOutput string) []Image {
	var images []Image
	lines := strings.Split(imagesOutput, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		// Parse the line based on our custom format
		parts := strings.Split(line, "|")
		if len(parts) < 5 {
			continue
		}

		// Extract basic information
		repository := parts[0]
		tag := parts[1]
		imageID := parts[2]
		created := parts[3]
		size := parts[4]

		// Extract optional digest if present
		digest := ""
		if len(parts) > 5 {
			digest = parts[5]
		}

		image := Image{
			Repository: repository,
			Tag:        tag,
			ImageID:    imageID,
			Created:    created,
			Size:       size,
			Digest:     digest,
		}

		images = append(images, image)
	}

	return images
}
