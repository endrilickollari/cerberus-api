package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// GetImageDetail implements the Service interface
func (s *service) GetImageDetail(ctx context.Context, sessionID string, imageID string) (*ImageDetail, error) {
	// Sanitize image ID to prevent command injection
	sanitizedImageID := sanitizeImageID(imageID)

	// Build the docker inspect command for images
	inspectCmd := fmt.Sprintf("docker image inspect %s", sanitizedImageID)

	// Execute the command
	output, err := s.sessionRepo.RunCommand(ctx, sessionID, inspectCmd)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect image: %w", err)
	}

	// Parse the output - it should be a JSON array with one element
	var inspectResults []map[string]interface{}
	if err := json.Unmarshal([]byte(output), &inspectResults); err != nil {
		return nil, fmt.Errorf("failed to parse docker image inspect output: %w", err)
	}

	// Check if we got any results
	if len(inspectResults) == 0 {
		return nil, fmt.Errorf("image not found: %s", imageID)
	}

	// Get the image data from the first element
	imageData := inspectResults[0]

	// Parse the image data into our model
	imageDetail, err := parseImageDetail(imageData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse image data: %w", err)
	}

	return imageDetail, nil
}

// parseImageDetail parses the image data from docker image inspect into our model
func parseImageDetail(data map[string]interface{}) (*ImageDetail, error) {
	imageDetail := &ImageDetail{}

	// Extract basic information
	if id, ok := data["Id"].(string); ok {
		imageDetail.ID = id
	}

	// Extract repo tags
	if repoTags, ok := data["RepoTags"].([]interface{}); ok {
		for _, tag := range repoTags {
			if tagStr, ok := tag.(string); ok {
				imageDetail.RepoTags = append(imageDetail.RepoTags, tagStr)
			}
		}
	}

	// Extract repo digests
	if repoDigests, ok := data["RepoDigests"].([]interface{}); ok {
		for _, digest := range repoDigests {
			if digestStr, ok := digest.(string); ok {
				imageDetail.RepoDigests = append(imageDetail.RepoDigests, digestStr)
			}
		}
	}

	// Extract created time
	if createdStr, ok := data["Created"].(string); ok {
		created, err := time.Parse(time.RFC3339, createdStr)
		if err == nil {
			imageDetail.Created = created
		}
	}

	// Extract size information
	if size, ok := data["Size"].(float64); ok {
		imageDetail.Size = int64(size)
	}

	if virtualSize, ok := data["VirtualSize"].(float64); ok {
		imageDetail.VirtualSize = int64(virtualSize)
	}

	if sharedSize, ok := data["SharedSize"].(float64); ok {
		imageDetail.SharedSize = int64(sharedSize)
	}

	// Extract architecture and OS
	if architecture, ok := data["Architecture"].(string); ok {
		imageDetail.Architecture = architecture
	}

	if os, ok := data["Os"].(string); ok {
		imageDetail.OS = os
	}

	// Extract author and container info
	if author, ok := data["Author"].(string); ok {
		imageDetail.Author = author
	}

	if container, ok := data["Container"].(string); ok {
		imageDetail.Container = container
	}

	if dockerVersion, ok := data["DockerVersion"].(string); ok {
		imageDetail.DockerVersion = dockerVersion
	}

	// Extract container config information
	if config, ok := data["Config"].(map[string]interface{}); ok {
		// Extract labels
		if labels, ok := config["Labels"].(map[string]interface{}); ok {
			imageDetail.Labels = make(map[string]string)
			for k, v := range labels {
				if str, ok := v.(string); ok {
					imageDetail.Labels[k] = str
				}
			}
		}

		// Extract environment variables
		if env, ok := config["Env"].([]interface{}); ok {
			for _, e := range env {
				if envStr, ok := e.(string); ok {
					imageDetail.Env = append(imageDetail.Env, envStr)
				}
			}
		}

		// Extract command
		if cmd, ok := config["Cmd"].([]interface{}); ok {
			for _, c := range cmd {
				if cmdStr, ok := c.(string); ok {
					imageDetail.Cmd = append(imageDetail.Cmd, cmdStr)
				}
			}
		}

		// Extract entrypoint
		if entrypoint, ok := config["Entrypoint"].([]interface{}); ok {
			for _, e := range entrypoint {
				if entryStr, ok := e.(string); ok {
					imageDetail.Entrypoint = append(imageDetail.Entrypoint, entryStr)
				}
			}
		}

		// Extract working directory
		if workingDir, ok := config["WorkingDir"].(string); ok {
			imageDetail.WorkingDir = workingDir
		}

		// Extract volumes
		if volumes, ok := config["Volumes"].(map[string]interface{}); ok {
			imageDetail.Volumes = make(map[string]struct{})
			for k := range volumes {
				imageDetail.Volumes[k] = struct{}{}
			}
		}

		// Extract exposed ports
		if exposedPorts, ok := config["ExposedPorts"].(map[string]interface{}); ok {
			imageDetail.ExposedPorts = make(map[string]struct{})
			for k := range exposedPorts {
				imageDetail.ExposedPorts[k] = struct{}{}
			}
		}
	}

	// Extract layers
	if rootFS, ok := data["RootFS"].(map[string]interface{}); ok {
		if layers, ok := rootFS["Layers"].([]interface{}); ok {
			for _, layer := range layers {
				if layerStr, ok := layer.(string); ok {
					imageDetail.Layers = append(imageDetail.Layers, layerStr)
				}
			}
		}
	}

	// Extract history
	if history, ok := data["History"].([]interface{}); ok {
		for _, h := range history {
			if historyMap, ok := h.(map[string]interface{}); ok {
				historyEntry := ImageHistory{}

				if createdStr, ok := historyMap["created"].(string); ok {
					created, err := time.Parse(time.RFC3339, createdStr)
					if err == nil {
						historyEntry.Created = created
					}
				}

				if createdBy, ok := historyMap["created_by"].(string); ok {
					historyEntry.CreatedBy = createdBy
				}

				if emptyLayer, ok := historyMap["empty_layer"].(bool); ok {
					historyEntry.EmptyLayer = emptyLayer
				}

				if comment, ok := historyMap["comment"].(string); ok {
					historyEntry.Comment = comment
				}

				imageDetail.History = append(imageDetail.History, historyEntry)
			}
		}
	}

	return imageDetail, nil
}

// sanitizeImageID sanitizes an image ID or name to prevent command injection
func sanitizeImageID(imageID string) string {
	// For image IDs (sha256:...)
	if strings.HasPrefix(imageID, "sha256:") {
		// Extract the hexadecimal part
		hexPart := strings.TrimPrefix(imageID, "sha256:")

		// Allow only hexadecimal characters
		validChars := "0123456789abcdefABCDEF"

		var result strings.Builder
		for _, char := range hexPart {
			if strings.ContainsRune(validChars, char) {
				result.WriteRune(char)
			}
		}

		return "sha256:" + result.String()
	}

	// For image names (repository:tag)
	// Allow only characters valid in a Docker image name
	re := regexp.MustCompile(`[^a-zA-Z0-9-_\/.:]`)
	sanitized := re.ReplaceAllString(imageID, "")

	// Limit to reasonable length to avoid overflow
	if len(sanitized) > 128 {
		sanitized = sanitized[:128]
	}

	return sanitized
}
