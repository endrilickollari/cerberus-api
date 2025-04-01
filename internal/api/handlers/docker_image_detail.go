package handlers

import (
	"errors"
	"net/http"
	"strings"

	"remote-server-api/internal/api/response"
	"remote-server-api/internal/domain/docker"
)

// GetImageDetail GetImageDetailHandler returns detailed information about a specific Docker image
//
// @Summary Get detailed Docker image information
// @Description Retrieves detailed information about a specific Docker image using docker image inspect
// @Tags docker
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <token>"
// @Param image_id path string true "Image ID or name:tag"
// @Success 200 {object} docker.ImageDetail "Docker image details retrieved successfully"
// @Failure 400 {object} response.Response "Invalid image ID"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Image not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /docker/image/{image_id} [get]
func (h *DockerHandler) GetImageDetail(w http.ResponseWriter, r *http.Request) {
	// Get session ID from context
	sessionID, ok := r.Context().Value(SessionIDKey).(string)
	if !ok {
		response.Error(w, "Session not found", http.StatusUnauthorized)
		return
	}

	// Get image ID from URL path parameter
	imageID := r.PathValue("image_id")
	if imageID == "" {
		response.Error(w, "Image ID or name is required", http.StatusBadRequest)
		return
	}

	// Get image details
	imageDetail, err := h.dockerService.GetImageDetail(r.Context(), sessionID, imageID)
	if err != nil {
		// Handle specific errors
		switch {
		case errors.Is(err, docker.ErrSessionNotFound):
			response.Error(w, "Session expired or not found", http.StatusUnauthorized)
		case strings.Contains(err.Error(), "image not found"):
			response.Error(w, "Image not found: "+imageID, http.StatusNotFound)
		default:
			response.Error(w, "Failed to get image details: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Return the image details
	response.JSON(w, imageDetail, http.StatusOK)
}
