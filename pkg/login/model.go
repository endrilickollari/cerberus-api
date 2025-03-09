package login

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// SSHLogin represents the credentials required for SSH login.
type SSHLogin struct {
	// IP address of the SSH server.
	IP string `json:"ip"`
	// Username for SSH login.
	Username string `json:"username"`
	// Port number for SSH connection (as a string).
	Port string `json:"port"`
	// Password for SSH login.
	Password string `json:"password"`
}

// Claims represents the claims embedded in the JWT token.
type Claims struct {
	// Username associated with the token.
	Username string `json:"username"`
	// Session ID related to the SSH session.
	SessionID string `json:"session_id"`
	// Standard JWT registered claims.
	jwt.RegisteredClaims
}

// Valid checks if the JWT token is valid and not expired.
func (c *Claims) Valid() error {
	if c.ExpiresAt != nil && time.Now().After(c.ExpiresAt.Time) {
		return jwt.NewValidationError("Token is expired", jwt.ValidationErrorExpired)
	}
	return nil
}
