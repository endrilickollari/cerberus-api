package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"remote-server-api/internal/domain/auth"
)

// JWTService implements token management using JWT
type JWTService struct {
	secretKey []byte
	expiresIn time.Duration
}

// NewJWTService creates a new JWT service
func NewJWTService(secretKey []byte, expiresIn time.Duration) *JWTService {
	return &JWTService{
		secretKey: secretKey,
		expiresIn: expiresIn,
	}
}

// GenerateToken generates a new JWT token
func (s *JWTService) GenerateToken(username, sessionID string) (string, error) {
	expirationTime := time.Now().Add(s.expiresIn)
	claims := &auth.Claims{
		Username:  username,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *JWTService) ValidateToken(tokenString string) (*auth.Claims, error) {
	claims := &auth.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
