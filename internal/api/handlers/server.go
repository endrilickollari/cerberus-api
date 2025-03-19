package handlers

import (
	"errors"
	"net/http"

	"remote-server-api/internal/api/response"
	"remote-server-api/internal/domain/server"
)

// ServerHandler handles server details requests
type ServerHandler struct {
	serverService server.Service
}

// NewServerHandler creates a new server details handler
func NewServerHandler(serverService server.Service) *ServerHandler {
	return &ServerHandler{
		serverService: serverService,
	}
}

// GetBasicDetails returns basic server information
//
// @Summary Get basic server details
// @Description Retrieves basic server information like hostname, OS, kernel version, etc.
// @Tags server
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} server.ServerDetails "Server details retrieved successfully"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /server-details [get]
func (h *ServerHandler) GetBasicDetails(w http.ResponseWriter, r *http.Request) {
	// Get session ID from context
	sessionID, ok := r.Context().Value(SessionIDKey).(string)
	if !ok {
		response.Error(w, "Session not found", http.StatusUnauthorized)
		return
	}

	// Get server details
	details, err := h.serverService.GetBasicDetails(r.Context(), sessionID)
	if err != nil {
		// Handle specific errors
		switch {
		case errors.Is(err, server.ErrSessionNotFound):
			response.Error(w, "Session expired or not found", http.StatusUnauthorized)
		default:
			response.Error(w, "Failed to get server details: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Return the details
	response.JSON(w, details, http.StatusOK)
}

// GetCPUInfo returns CPU information
//
// @Summary Get CPU information
// @Description Retrieves detailed CPU information from the server
// @Tags server
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {array} server.CPUInfo "CPU information retrieved successfully"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /server-details/cpu-info [get]
func (h *ServerHandler) GetCPUInfo(w http.ResponseWriter, r *http.Request) {
	// Get session ID from context
	sessionID, ok := r.Context().Value(SessionIDKey).(string)
	if !ok {
		response.Error(w, "Session not found", http.StatusUnauthorized)
		return
	}

	// Get CPU info
	cpuInfo, err := h.serverService.GetCPUInfo(r.Context(), sessionID)
	if err != nil {
		// Handle specific errors
		switch {
		case errors.Is(err, server.ErrSessionNotFound):
			response.Error(w, "Session expired or not found", http.StatusUnauthorized)
		default:
			response.Error(w, "Failed to get CPU info: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Return the CPU info
	response.JSON(w, cpuInfo, http.StatusOK)
}

// GetDiskUsage returns disk usage information
//
// @Summary Get disk usage information
// @Description Retrieves disk usage information from the server
// @Tags server
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {array} server.DiskUsage "Disk usage information retrieved successfully"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /server-details/disk-usage [get]
func (h *ServerHandler) GetDiskUsage(w http.ResponseWriter, r *http.Request) {
	// Get session ID from context
	sessionID, ok := r.Context().Value(SessionIDKey).(string)
	if !ok {
		response.Error(w, "Session not found", http.StatusUnauthorized)
		return
	}

	// Get disk usage
	diskUsage, err := h.serverService.GetDiskUsage(r.Context(), sessionID)
	if err != nil {
		// Handle specific errors
		switch {
		case errors.Is(err, server.ErrSessionNotFound):
			response.Error(w, "Session expired or not found", http.StatusUnauthorized)
		default:
			response.Error(w, "Failed to get disk usage: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Return the disk usage
	response.JSON(w, diskUsage, http.StatusOK)
}

// GetRunningProcesses returns information about running processes
//
// @Summary Get running processes information
// @Description Retrieves information about running processes on the server
// @Tags server
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {array} server.ProcessInfo "Running processes information retrieved successfully"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /server-details/running-processes [get]
func (h *ServerHandler) GetRunningProcesses(w http.ResponseWriter, r *http.Request) {
	// Get session ID from context
	sessionID, ok := r.Context().Value(SessionIDKey).(string)
	if !ok {
		response.Error(w, "Session not found", http.StatusUnauthorized)
		return
	}

	// Get running processes
	processes, err := h.serverService.GetRunningProcesses(r.Context(), sessionID)
	if err != nil {
		// Handle specific errors
		switch {
		case errors.Is(err, server.ErrSessionNotFound):
			response.Error(w, "Session expired or not found", http.StatusUnauthorized)
		default:
			response.Error(w, "Failed to get running processes: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Return the running processes
	response.JSON(w, processes, http.StatusOK)
}
