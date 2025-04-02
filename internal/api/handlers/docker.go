package handlers

import (
	"errors"
	"net/http"
	"remote-server-api/internal/api/response"
	"remote-server-api/internal/domain/docker"
)

// DockerHandler handles Docker-related requests
type DockerHandler struct {
	dockerService docker.Service
}

// NewDockerHandler creates a new Docker handler
func NewDockerHandler(dockerService docker.Service) *DockerHandler {
	return &DockerHandler{
		dockerService: dockerService,
	}
}

// GetContainerInfo returns information about Docker containers
//
// @Summary Get Docker container information
// @Description Retrieves information about Docker containers from the server
// @Tags docker
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {array} docker.Container "Docker container information retrieved successfully"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /docker/containers [get]
func (h *DockerHandler) GetContainerInfo(w http.ResponseWriter, r *http.Request) {
	// Get session ID from context
	sessionID, ok := r.Context().Value(SessionIDKey).(string)
	if !ok {
		response.Error(w, "Session not found", http.StatusUnauthorized)
		return
	}

	// Get Docker containers
	containers, err := h.dockerService.GetContainers(r.Context(), sessionID)
	if err != nil {
		// Handle specific errors
		switch {
		case errors.Is(err, docker.ErrSessionNotFound):
			response.Error(w, "Session expired or not found", http.StatusUnauthorized)
		default:
			response.Error(w, "Failed to get Docker containers: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Return the Docker containers
	response.JSON(w, containers, http.StatusOK)
}
