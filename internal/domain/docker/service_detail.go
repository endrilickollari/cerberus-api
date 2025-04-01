package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// GetContainerDetail implements the Service interface
func (s *service) GetContainerDetail(ctx context.Context, sessionID string, containerID string) (*ContainerDetail, error) {
	// Sanitize container ID to prevent command injection
	sanitizedContainerID := sanitizeContainerID(containerID)

	// Build the docker inspect command
	// We're using a combination of format and full JSON to get both structured and raw data
	inspectCmd := fmt.Sprintf("docker inspect %s", sanitizedContainerID)

	// Execute the command
	output, err := s.sessionRepo.RunCommand(ctx, sessionID, inspectCmd)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect container: %w", err)
	}

	// Parse the output - it should be a JSON array with one element
	var inspectResults []map[string]interface{}
	if err := json.Unmarshal([]byte(output), &inspectResults); err != nil {
		return nil, fmt.Errorf("failed to parse docker inspect output: %w", err)
	}

	// Check if we got any results
	if len(inspectResults) == 0 {
		return nil, fmt.Errorf("container not found: %s", containerID)
	}

	// Get the container data from the first element
	containerData := inspectResults[0]

	// Parse the container data into our model
	containerDetail, err := parseContainerDetail(containerData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse container data: %w", err)
	}

	return containerDetail, nil
}

// parseContainerDetail parses the container data from docker inspect into our model
func parseContainerDetail(data map[string]interface{}) (*ContainerDetail, error) {
	containerDetail := &ContainerDetail{}

	// Extract basic information
	if id, ok := data["Id"].(string); ok {
		containerDetail.ID = id
	}

	if name, ok := data["Name"].(string); ok {
		// Remove leading slash from name
		containerDetail.Name = strings.TrimPrefix(name, "/")
	}

	// Extract configuration information
	if config, ok := data["Config"].(map[string]interface{}); ok {
		if image, ok := config["Image"].(string); ok {
			containerDetail.Image = image
		}

		// Extract command (can be string or array)
		if cmd, ok := config["Cmd"].([]interface{}); ok {
			for _, c := range cmd {
				if str, ok := c.(string); ok {
					containerDetail.Command = append(containerDetail.Command, str)
				}
			}
		}

		// Extract labels
		if labels, ok := config["Labels"].(map[string]interface{}); ok {
			containerDetail.Labels = make(map[string]string)
			for k, v := range labels {
				if str, ok := v.(string); ok {
					containerDetail.Labels[k] = str
				}
			}
		}
	}

	// Extract creation time
	if createdStr, ok := data["Created"].(string); ok {
		created, err := time.Parse(time.RFC3339, createdStr)
		if err == nil {
			containerDetail.Created = created
		}
	}

	// Extract state information
	if state, ok := data["State"].(map[string]interface{}); ok {
		containerState := ContainerState{}

		if status, ok := state["Status"].(string); ok {
			containerState.Status = status
		}

		if running, ok := state["Running"].(bool); ok {
			containerState.Running = running
		}

		if paused, ok := state["Paused"].(bool); ok {
			containerState.Paused = paused
		}

		if restarting, ok := state["Restarting"].(bool); ok {
			containerState.Restarting = restarting
		}

		if startedAt, ok := state["StartedAt"].(string); ok {
			started, err := time.Parse(time.RFC3339, startedAt)
			if err == nil {
				containerState.StartedAt = started
			}
		}

		if finishedAt, ok := state["FinishedAt"].(string); ok {
			finished, err := time.Parse(time.RFC3339, finishedAt)
			if err == nil {
				containerState.FinishedAt = finished
			}
		}

		if exitCode, ok := state["ExitCode"].(float64); ok {
			containerState.ExitCode = int(exitCode)
		}

		if error, ok := state["Error"].(string); ok {
			containerState.Error = error
		}

		containerDetail.State = containerState
	}

	// Extract host configuration
	if hostConfig, ok := data["HostConfig"].(map[string]interface{}); ok {
		hostConfigDetail := HostConfig{}

		if autoRemove, ok := hostConfig["AutoRemove"].(bool); ok {
			hostConfigDetail.AutoRemove = autoRemove
		}

		if privileged, ok := hostConfig["Privileged"].(bool); ok {
			hostConfigDetail.Privileged = privileged
		}

		if publishAllPorts, ok := hostConfig["PublishAllPorts"].(bool); ok {
			hostConfigDetail.PublishAllPorts = publishAllPorts
		}

		if restartPolicy, ok := hostConfig["RestartPolicy"].(map[string]interface{}); ok {
			if name, ok := restartPolicy["Name"].(string); ok {
				hostConfigDetail.RestartPolicy = name
			}
		}

		// Extract capabilities
		if capAdd, ok := hostConfig["CapAdd"].([]interface{}); ok {
			for _, cap := range capAdd {
				if str, ok := cap.(string); ok {
					hostConfigDetail.CapAdd = append(hostConfigDetail.CapAdd, str)
				}
			}
		}

		if capDrop, ok := hostConfig["CapDrop"].([]interface{}); ok {
			for _, cap := range capDrop {
				if str, ok := cap.(string); ok {
					hostConfigDetail.CapDrop = append(hostConfigDetail.CapDrop, str)
				}
			}
		}

		// Extract DNS settings
		if dns, ok := hostConfig["Dns"].([]interface{}); ok {
			for _, d := range dns {
				if str, ok := d.(string); ok {
					hostConfigDetail.DNS = append(hostConfigDetail.DNS, str)
				}
			}
		}

		containerDetail.HostConfig = hostConfigDetail
	}

	// Extract network settings
	if networkSettings, ok := data["NetworkSettings"].(map[string]interface{}); ok {
		networkSettingsDetail := NetworkSettings{}

		if ipAddress, ok := networkSettings["IPAddress"].(string); ok {
			networkSettingsDetail.IPAddress = ipAddress
		}

		if gateway, ok := networkSettings["Gateway"].(string); ok {
			networkSettingsDetail.Gateway = gateway
		}

		if macAddress, ok := networkSettings["MacAddress"].(string); ok {
			networkSettingsDetail.MacAddress = macAddress
		}

		// Extract port mappings
		if ports, ok := networkSettings["Ports"].(map[string]interface{}); ok {
			for containerPort, hostPorts := range ports {
				if hostPortsList, ok := hostPorts.([]interface{}); ok {
					for _, hostPortMapping := range hostPortsList {
						if hostPortMap, ok := hostPortMapping.(map[string]interface{}); ok {
							mapping := PortMapping{
								ContainerPort: containerPort,
							}

							if hostPort, ok := hostPortMap["HostPort"].(string); ok {
								mapping.HostPort = hostPort
							}

							// Extract protocol from the container port (e.g., "80/tcp")
							parts := strings.Split(containerPort, "/")
							if len(parts) > 1 {
								mapping.Protocol = parts[1]
							}

							networkSettingsDetail.PortMappings = append(networkSettingsDetail.PortMappings, mapping)
						}
					}
				}
			}
		}

		// Extract network information
		if networks, ok := networkSettings["Networks"].(map[string]interface{}); ok {
			// Just take the first network for simplicity
			for networkName, networkInfo := range networks {
				networkSettingsDetail.NetworkName = networkName

				if network, ok := networkInfo.(map[string]interface{}); ok {
					if ipAddress, ok := network["IPAddress"].(string); ok {
						networkSettingsDetail.IPAddress = ipAddress
					}

					if gateway, ok := network["Gateway"].(string); ok {
						networkSettingsDetail.Gateway = gateway
					}

					if endpointID, ok := network["EndpointID"].(string); ok {
						networkSettingsDetail.EndpointID = endpointID
					}

					if networkID, ok := network["NetworkID"].(string); ok {
						networkSettingsDetail.NetworkID = networkID
					}

					if ipPrefixLen, ok := network["IPPrefixLen"].(float64); ok {
						networkSettingsDetail.SubnetPrefix = fmt.Sprintf("/%d", int(ipPrefixLen))
					}
				}

				// Just use the first network
				break
			}
		}

		containerDetail.NetworkSettings = networkSettingsDetail
	}

	// Extract mounts
	if mounts, ok := data["Mounts"].([]interface{}); ok {
		for _, m := range mounts {
			if mountInfo, ok := m.(map[string]interface{}); ok {
				mount := Mount{}

				if mountType, ok := mountInfo["Type"].(string); ok {
					mount.Type = mountType
				}

				if source, ok := mountInfo["Source"].(string); ok {
					mount.Source = source
				}

				if destination, ok := mountInfo["Destination"].(string); ok {
					mount.Destination = destination
				}

				if mode, ok := mountInfo["Mode"].(string); ok {
					mount.Mode = mode
				}

				if rw, ok := mountInfo["RW"].(bool); ok {
					mount.RW = rw
				}

				containerDetail.Mounts = append(containerDetail.Mounts, mount)
			}
		}
	}

	// Extract platform
	if platform, ok := data["Platform"].(string); ok {
		containerDetail.Platform = platform
	}

	return containerDetail, nil
}

// sanitizeContainerID sanitizes a container ID to prevent command injection
func sanitizeContainerID(containerID string) string {
	// Allow only hexadecimal characters and short IDs (first 12 chars of a full ID)
	validChars := "0123456789abcdefABCDEF"

	var result strings.Builder
	for _, char := range containerID {
		if strings.ContainsRune(validChars, char) {
			result.WriteRune(char)
		}
	}

	// Limit to reasonable length to avoid overflow
	sanitized := result.String()
	if len(sanitized) > 64 {
		sanitized = sanitized[:64]
	}

	return sanitized
}
