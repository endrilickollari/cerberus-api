package handlers

import (
	"errors"
	"net/http"
	"remote-server-api/internal/domain/docker"
	"strings"

	"remote-server-api/internal/api/response"
)

// GetContainerDetail GetContainerDetailHandler returns detailed information about a specific Docker container
//
// @Summary Get detailed Docker container information
// @Description Retrieves detailed information about a specific Docker container using docker inspect
// @Tags docker
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <token>"
// @Param container_id path string true "Container ID"
// @Success 200 {object} docker.ContainerDetail "Docker container details retrieved successfully"
// @Failure 400 {object} response.Response "Invalid container ID"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Container not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /docker/container/{container_id} [get]
func (h *DockerHandler) GetContainerDetail(w http.ResponseWriter, r *http.Request) {
	// Get session ID from context
	sessionID, ok := r.Context().Value(SessionIDKey).(string)
	if !ok {
		response.Error(w, "Session not found", http.StatusUnauthorized)
		return
	}

	// Get container ID from URL path parameter
	containerID := r.PathValue("container_id")
	if containerID == "" {
		response.Error(w, "Container ID is required", http.StatusBadRequest)
		return
	}

	// Get container details
	containerDetail, err := h.dockerService.GetContainerDetail(r.Context(), sessionID, containerID)
	if err != nil {
		// Handle specific errors
		switch {
		case errors.Is(err, docker.ErrSessionNotFound):
			response.Error(w, "Session expired or not found", http.StatusUnauthorized)
		case strings.Contains(err.Error(), "container not found"):
			response.Error(w, "Container not found: "+containerID, http.StatusNotFound)
		default:
			response.Error(w, "Failed to get container details: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Return the container details
	response.JSON(w, containerDetail, http.StatusOK)
}
