package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"remote-server-api/internal/api/response"
	"remote-server-api/internal/domain/auth"
)

// AuthHandler handles authentication related requests
type AuthHandler struct {
	authService auth.Service
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService auth.Service) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles SSH login requests and generates a JWT token
//
// @Summary Login to SSH and generate JWT token
// @Description Authenticates user against an SSH server and returns a JWT token for subsequent API requests
// @Tags authentication
// @Accept json
// @Produce json
// @Param body body auth.LoginRequest true "SSH login credentials"
// @Success 200 {object} auth.LoginResponse "Successfully logged in and token generated"
// @Failure 400 {object} response.Response "Invalid request payload"
// @Failure 500 {object} response.Response "Failed to connect to SSH server or generate token"
// @Router /login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Attempt to login
	loginResp, err := h.authService.Login(r.Context(), req)
	if err != nil {
		// Handle specific errors
		switch {
		case errors.Is(err, auth.ErrInvalidCredentials):
			response.Error(w, "Invalid SSH credentials", http.StatusUnauthorized)
		default:
			response.Error(w, "Failed to authenticate: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Return the token
	response.JSON(w, loginResp, http.StatusOK)
}
