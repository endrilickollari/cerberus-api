package handlers

import (
	"errors"
	"net/http"
	"remote-server-api/internal/api/response"
	"remote-server-api/internal/domain/docker"
)

// GetImages GetImagesHandler returns information about Docker images
//
// @Summary Get Docker images
// @Description Retrieves information about all Docker images on the server
// @Tags docker
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {array} docker.Image "Docker images retrieved successfully"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /docker/images [get]
func (h *DockerHandler) GetImages(w http.ResponseWriter, r *http.Request) {
	// Get session ID from context
	sessionID, ok := r.Context().Value(SessionIDKey).(string)
	if !ok {
		response.Error(w, "Session not found", http.StatusUnauthorized)
		return
	}

	// Get Docker images
	images, err := h.dockerService.GetImages(r.Context(), sessionID)
	if err != nil {
		// Handle specific errors
		switch {
		case errors.Is(err, docker.ErrSessionNotFound):
			response.Error(w, "Session expired or not found", http.StatusUnauthorized)
		default:
			response.Error(w, "Failed to get Docker images: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Return the Docker images
	response.JSON(w, images, http.StatusOK)
}
