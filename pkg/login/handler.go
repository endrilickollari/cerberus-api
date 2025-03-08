package login

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var sshLogin SSHLogin

	err := json.NewDecoder(r.Body).Decode(&sshLogin)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Connect to SSH
	client, err := ConnectToSSH(sshLogin.IP, sshLogin.Username, sshLogin.Port, sshLogin.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to SSH server: %v", err), http.StatusInternalServerError)
		return
	}

	// Generate token
	sessionID := "some_unique_id"

	// Store the SSH client in memory
	StoreSession(sessionID, client)

	// Generate JWT token with the session ID
	token, err := GenerateToken(sshLogin.Username, sessionID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return the token to the user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"token":"%s"}`, token)))
}

func TokenValidationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		// Extract the token from the "Bearer <token>" format
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil // return the secret key for validation
		})

		if err != nil || !tkn.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Call the next handler if everything is valid
		next.ServeHTTP(w, r)
	}
}
