package handlers

import (
	"errors"
	"net/http"
	"remote-server-api/internal/domain/docker"
	"strconv"
	"strings"

	"remote-server-api/internal/api/response"
)

// DeleteImage DeleteImageHandler deletes a Docker image
//
// @Summary Delete a Docker image
// @Description Deletes a Docker image by ID or name
// @Tags docker
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <token>"
// @Param image_id path string true "Image ID or name:tag"
// @Param force query bool false "Force deletion even if image is being used by containers" default(false)
// @Success 200 {object} docker.ImageDeleteResponse "Docker image deleted successfully"
// @Failure 400 {object} response.Response "Invalid image ID or name"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /docker/image/{image_id} [delete]
func (h *DockerHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
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

	// Get force parameter from query
	forceStr := r.URL.Query().Get("force")
	force := false
	if forceStr != "" {
		var err error
		force, err = strconv.ParseBool(forceStr)
		if err != nil {
			response.Error(w, "Invalid force parameter: must be true or false", http.StatusBadRequest)
			return
		}
	}

	// Delete the image
	deleteResponse, err := h.dockerService.DeleteImage(r.Context(), sessionID, imageID, force)
	if err != nil {
		// Handle specific errors
		switch {
		case errors.Is(err, docker.ErrSessionNotFound):
			response.Error(w, "Session expired or not found", http.StatusUnauthorized)
		case strings.Contains(err.Error(), "No such image"):
			response.Error(w, "Image not found: "+imageID, http.StatusNotFound)
		default:
			response.Error(w, "Failed to delete image: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Return the delete response
	response.JSON(w, deleteResponse, http.StatusOK)
}
