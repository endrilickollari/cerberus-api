package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/ssh"
)

// LoginRequest represents the credentials required for SSH login
type LoginRequest struct {
	IP       string `json:"ip" validate:"required,ip"`
	Username string `json:"username" validate:"required"`
	Port     string `json:"port" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the successful login response
type LoginResponse struct {
	Token string `json:"token"`
}

// Claims represents the claims embedded in the JWT token
type Claims struct {
	Username  string `json:"username"`
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}

// Session represents an active SSH session
type Session struct {
	ID       string
	Username string
	Client   *ssh.Client
}
