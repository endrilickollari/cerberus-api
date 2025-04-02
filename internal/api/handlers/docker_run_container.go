package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"remote-server-api/internal/domain/docker"
	"strings"

	"remote-server-api/internal/api/response"
)

// RunContainer RunContainerHandler runs a Docker container
//
// @Summary Run a Docker container
// @Description Creates and runs a new Docker container from an image with optional configuration
// @Tags docker
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <token>"
// @Param request body docker.ContainerRunRequest true "Container configuration"
// @Success 201 {object} docker.ContainerRunResponse "Container created successfully"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Image not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /docker/image/run [post]
func (h *DockerHandler) RunContainer(w http.ResponseWriter, r *http.Request) {
	// Get session ID from context
	sessionID, ok := r.Context().Value(SessionIDKey).(string)
	if !ok {
		response.Error(w, "Session not found", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var request docker.ContainerRunRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Set defaults if not provided
	if request.Detached == false {
		// Default to detached mode
		request.Detached = true
	}

	// Run the container
	containerResponse, err := h.dockerService.RunContainer(r.Context(), sessionID, request)
	if err != nil {
		// Handle specific errors
		switch {
		case errors.Is(err, docker.ErrSessionNotFound):
			response.Error(w, "Session expired or not found", http.StatusUnauthorized)
		case strings.Contains(err.Error(), "invalid request"):
			response.Error(w, err.Error(), http.StatusBadRequest)
		case strings.Contains(err.Error(), "No such image"):
			response.Error(w, "Image not found: "+request.Image, http.StatusNotFound)
		default:
			response.Error(w, "Failed to run container: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Return the container response
	response.JSON(w, containerResponse, http.StatusCreated)
}
