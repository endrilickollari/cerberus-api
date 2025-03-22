package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"remote-server-api/internal/api/response"
	"remote-server-api/internal/domain/server"
)

// FileSystemHandler handles file system-related requests
type FileSystemHandler struct {
	serverService server.Service
}

// NewFileSystemHandler creates a new file system handler
func NewFileSystemHandler(serverService server.Service) *FileSystemHandler {
	return &FileSystemHandler{
		serverService: serverService,
	}
}

// ListFileSystem returns a listing of files and directories
//
// @Summary List files and directories
// @Description Retrieves a listing of files and directories at the specified path
// @Tags filesystem
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <token>"
// @Param path query string false "Path to list (defaults to /)"
// @Param recursive query bool false "Whether to list recursively" default(false)
// @Param include_hidden query bool false "Whether to include hidden files/directories" default(false)
// @Success 200 {object} server.FileSystemListing "File system listing retrieved successfully"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /filesystem/list [get]
func (h *FileSystemHandler) ListFileSystem(w http.ResponseWriter, r *http.Request) {
	// Get session ID from context
	sessionID, ok := r.Context().Value(SessionIDKey).(string)
	if !ok {
		response.Error(w, "Session not found", http.StatusUnauthorized)
		return
	}

	// Get query parameters
	path := r.URL.Query().Get("path")
	if path == "" {
		path = "/"
	}

	recursiveStr := r.URL.Query().Get("recursive")
	recursive, _ := strconv.ParseBool(recursiveStr)

	includeHiddenStr := r.URL.Query().Get("include_hidden")
	includeHidden, _ := strconv.ParseBool(includeHiddenStr)

	// Get file system listing
	listing, err := h.serverService.ListFileSystem(r.Context(), sessionID, path, recursive, includeHidden)
	if err != nil {
		// Handle specific errors
		switch {
		case errors.Is(err, server.ErrSessionNotFound):
			response.Error(w, "Session expired or not found", http.StatusUnauthorized)
		default:
			response.Error(w, "Failed to list file system: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Return the file system listing
	response.JSON(w, listing, http.StatusOK)
}

// FileDetailsRequest represents a request to get detailed file information
type FileDetailsRequest struct {
	Path string `json:"path" validate:"required"`
}

// GetFileDetails returns detailed information about a specific file
//
// @Summary Get file details
// @Description Retrieves detailed information about a specific file or directory
// @Tags filesystem
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <token>"
// @Param path query string true "Path to the file or directory"
// @Success 200 {object} server.FileSystemEntry "File details retrieved successfully"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "File not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /filesystem/details [get]
func (h *FileSystemHandler) GetFileDetails(w http.ResponseWriter, r *http.Request) {
	// Get session ID from context
	sessionID, ok := r.Context().Value(SessionIDKey).(string)
	if !ok {
		response.Error(w, "Session not found", http.StatusUnauthorized)
		return
	}

	// Get query parameters
	path := r.URL.Query().Get("path")
	if path == "" {
		response.Error(w, "Path is required", http.StatusBadRequest)
		return
	}

	// Get file details
	fileDetails, err := h.serverService.GetFileDetails(r.Context(), sessionID, path)
	if err != nil {
		// Handle specific errors
		switch {
		case errors.Is(err, server.ErrSessionNotFound):
			response.Error(w, "Session expired or not found", http.StatusUnauthorized)
		default:
			response.Error(w, "Failed to get file details: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// If file not found, return 404
	if fileDetails == nil {
		response.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Return the file details
	response.JSON(w, fileDetails, http.StatusOK)
}

// SearchFiles searches for files matching a pattern
//
// @Summary Search for files
// @Description Searches for files matching a pattern in the specified directory
// @Tags filesystem
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <token>"
// @Param path query string true "Path to search in"
// @Param pattern query string true "Search pattern (glob or regex)"
// @Param max_depth query int false "Maximum search depth" default(10)
// @Success 200 {array} server.FileSystemEntry "Search results retrieved successfully"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /filesystem/search [get]
func (h *FileSystemHandler) SearchFiles(w http.ResponseWriter, r *http.Request) {
	// Get session ID from context
	sessionID, ok := r.Context().Value(SessionIDKey).(string)
	if !ok {
		response.Error(w, "Session not found", http.StatusUnauthorized)
		return
	}

	// Get query parameters
	path := r.URL.Query().Get("path")
	if path == "" {
		path = "/"
	}

	pattern := r.URL.Query().Get("pattern")
	if pattern == "" {
		response.Error(w, "Search pattern is required", http.StatusBadRequest)
		return
	}

	maxDepthStr := r.URL.Query().Get("max_depth")
	maxDepth := 10 // default
	if maxDepthStr != "" {
		if depth, err := strconv.Atoi(maxDepthStr); err == nil && depth > 0 {
			maxDepth = depth
		}
	}

	// Search for files
	searchResults, err := h.serverService.SearchFiles(r.Context(), sessionID, path, pattern, maxDepth)
	if err != nil {
		// Handle specific errors
		switch {
		case errors.Is(err, server.ErrSessionNotFound):
			response.Error(w, "Session expired or not found", http.StatusUnauthorized)
		default:
			response.Error(w, "Failed to search files: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Return the search results
	response.JSON(w, searchResults, http.StatusOK)
}
