package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/meshyampratap01/letStayInn/internal/config"
)

// Claims struct
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a JWT for a user
func GenerateJWT(userID, username, role string, expirationMinutes int) (string, error) {
	expirationTime := time.Now().Add(time.Duration(expirationMinutes) * time.Minute)
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWTSecret))
}

// ValidateJWT parses and validates a JWT
func ValidateJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC and specifically HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenUnverifiable
		}
		return []byte(config.JWTSecret), nil
	})
	if err != nil {
		return nil, err // Return the specific error for debugging
	}
	if !token.Valid {
		return nil, jwt.ErrTokenSignatureInvalid
	}
	// Check for expiration explicitly
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, jwt.ErrTokenExpired
	}
	return claims, nil
}
