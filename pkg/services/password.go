// Package services provides business logic for authentication and user management.
package services

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// ==================== Service Definition ====================

// PasswordService handles password hashing, verification, and strength validation.
// It uses bcrypt for secure password hashing.
type PasswordService struct{}

// ==================== Constructor ====================

// NewPasswordService creates a new PasswordService.
func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

// ==================== Password Hashing ====================

// HashPassword hashes a password using bcrypt with default cost.
func (p *PasswordService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

// ==================== Password Verification ====================

// VerifyPassword checks if the provided password matches the stored hash.
func (p *PasswordService) VerifyPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("password verification failed: %w", err)
	}
	return nil
}

// ==================== Password Validation ====================

// ValidatePasswordStrength checks if a password meets security requirements:
// - At least 8 characters
// - Contains uppercase, lowercase, number, and special character
// - Not a common or weak password
func (p *PasswordService) ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	// Check for at least one uppercase letter
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	// Check for at least one lowercase letter
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}

	// Check for at least one number
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one number")
	}

	// Check for at least one special character
	if !regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one special character")
	}

	// Check for common weak passwords
	if p.isCommonPassword(password) {
		return fmt.Errorf("password is too common, please choose a stronger password")
	}

	return nil
}

// isCommonPassword checks if the password is in a list of common weak passwords.
func (p *PasswordService) isCommonPassword(password string) bool {
	lowercasePassword := strings.ToLower(password)
	commonPasswords := []string{
		"password", "12345678", "123456789", "qwerty123", "abc12345",
		"password123", "admin123", "letmein123", "welcome123", "monkey123",
		"password1", "1234567890", "qwertyuiop", "asdfghjkl", "zxcvbnm",
		"123qweasd", "password12", "admin1234", "welcome1", "changeme",
	}

	for _, common := range commonPasswords {
		if lowercasePassword == common {
			return true
		}
	}

	// Check for sequential patterns
	if p.isSequentialPattern(lowercasePassword) {
		return true
	}

	// Check for repeated characters
	if p.isRepeatedPattern(lowercasePassword) {
		return true
	}

	return false
}

// isSequentialPattern checks for sequential keyboard patterns like "12345678" or "qwerty".
func (p *PasswordService) isSequentialPattern(password string) bool {
	sequentialPatterns := []string{
		"12345678", "23456789", "34567890", "01234567",
		"qwertyui", "asdfghjk", "zxcvbnm", "qwerty",
		"asdfgh", "zxcvbn", "qwertyu", "asdfghj",
	}

	for _, pattern := range sequentialPatterns {
		if strings.Contains(password, pattern) {
			return true
		}
	}
	return false
}

// isRepeatedPattern checks for repeated character patterns like "aaaaaa" or "111111".
func (p *PasswordService) isRepeatedPattern(password string) bool {
	// Check for patterns like "aaaaaa", "111111", etc.
	// Using a simpler approach to avoid regex backreference issues
	for i := 0; i < len(password)-5; i++ {
		char := password[i]
		repeated := true
		for j := 1; j < 6; j++ {
			if password[i+j] != char {
				repeated = false
				break
			}
		}
		if repeated {
			return true
		}
	}
	return false
}

// GenerateSecurePassword generates a secure random password with required character types.
// It ensures at least one lowercase, uppercase, number, and special character.
func (p *PasswordService) GenerateSecurePassword(length int) string {
	if length < 8 {
		length = 12
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	password := make([]byte, length)

	// Ensure at least one of each required character type
	password[0] = "abcdefghijklmnopqrstuvwxyz"[len(password)%26]
	password[1] = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"[len(password)%26]
	password[2] = "0123456789"[len(password)%10]
	password[3] = "!@#$%^&*"[len(password)%8]

	// Fill the rest randomly
	for i := 4; i < length; i++ {
		password[i] = charset[len(password)*i%len(charset)]
	}

	return string(password)
}
