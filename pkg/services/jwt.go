// Package services provides business logic for authentication and user management.
package services

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tofiquem/assingment/pkg/models"
)

// ==================== Service Definition ====================

// JWTService handles JWT token generation and validation.
// It uses HMAC-SHA256 for signing tokens.
type JWTService struct {
	secretKey []byte
}

// ==================== Constructor ====================

// NewJWTService creates a new JWTService with the secret key from environment.
// Defaults to a warning message if JWT_SECRET is not set.
func NewJWTService() *JWTService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-super-secret-jwt-key-change-in-production"
		fmt.Println("WARNING: Using default JWT secret, please set JWT_SECRET environment variable")
	}
	return &JWTService{
		secretKey: []byte(secret),
	}
}

// ==================== Token Generation ====================

// GenerateToken creates a new JWT token for the given user.
// The token expires after 24 hours and contains user claims.
func (j *JWTService) GenerateToken(user *models.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(24 * time.Hour) // 24 hours expiration

	claims := jwt.MapClaims{
		"user_id":    user.ID,
		"email":      user.Email,
		"role":       user.Role,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"exp":        expiresAt.Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// ==================== Token Validation ====================

// ValidateToken validates a JWT token string and returns the claims.
// It verifies the token signature and structure.
func (j *JWTService) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ExtractClaims extracts claims from a token without full validation.
// Useful for inspecting token content before validation.
func (j *JWTService) ExtractClaims(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("no claims found in token")
}

// ==================== Claims Helpers ====================

// GetUserIDFromToken extracts the user ID from token claims.
func (j *JWTService) GetUserIDFromToken(claims jwt.MapClaims) (uint, error) {
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("user_id not found in token")
	}
	return uint(userID), nil
}

// GetUserEmailFromToken extracts the email from token claims.
func (j *JWTService) GetUserEmailFromToken(claims jwt.MapClaims) (string, error) {
	email, ok := claims["email"].(string)
	if !ok {
		return "", fmt.Errorf("email not found in token")
	}
	return email, nil
}

// GetUserRoleFromToken extracts the role from token claims.
func (j *JWTService) GetUserRoleFromToken(claims jwt.MapClaims) (string, error) {
	role, ok := claims["role"].(string)
	if !ok {
		return "", fmt.Errorf("role not found in token")
	}
	return role, nil
}

// IsTokenExpired checks if the token expiration time has passed.
func (j *JWTService) IsTokenExpired(claims jwt.MapClaims) bool {
	expClaim, exists := claims["exp"]
	if !exists {
		return true // If no exp claim, consider expired
	}

	var expTime int64
	switch v := expClaim.(type) {
	case float64:
		expTime = int64(v)
	case int64:
		expTime = v
	default:
		return true // Invalid exp claim format, consider expired
	}

	return time.Now().Unix() > expTime
}
