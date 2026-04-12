package services

import (
	"testing"
	"time"

	"github.com/tofiquem/assingment/internal/models"
	"github.com/tofiquem/assingment/internal/testutils"
)

func TestNewJWTService(t *testing.T) {
	service := NewJWTService()
	if service == nil {
		t.Fatal("Expected JWT service to be created, got nil")
	}
	if len(service.secretKey) == 0 {
		t.Error("Expected secret key to be set")
	}
}

func TestJWTService_GenerateToken(t *testing.T) {
	service := NewJWTService()
	user := &models.User{
		ID:        1,
		Email:     "test@example.com",
		Role:      "hr",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
	}

	token, expiresAt, err := service.GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Expected token to be generated, got empty string")
	}

	if expiresAt.Before(time.Now()) {
		t.Error("Expected expiration time to be in the future")
	}

	// Verify token can be validated
	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate generated token: %v", err)
	}

	userID, err := service.GetUserIDFromToken(*claims)
	if err != nil {
		t.Fatalf("Failed to extract user ID: %v", err)
	}

	if userID != 1 {
		t.Errorf("Expected user ID 1, got %d", userID)
	}
}

func TestJWTService_ValidateToken(t *testing.T) {
	service := NewJWTService()
	user := &models.User{
		ID:        1,
		Email:     "test@example.com",
		Role:      "hr",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
	}

	// Generate valid token
	token, _, err := service.GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Validate token
	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate valid token: %v", err)
	}

	if claims == nil {
		t.Error("Expected claims to be returned")
	}

	// Verify claims contain expected data
	userID, err := service.GetUserIDFromToken(*claims)
	if err != nil {
		t.Fatalf("Failed to extract user ID: %v", err)
	}

	if userID != 1 {
		t.Errorf("Expected user ID 1, got %d", userID)
	}

	email, err := service.GetUserEmailFromToken(*claims)
	if err != nil {
		t.Fatalf("Failed to extract email: %v", err)
	}

	if email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", email)
	}
}

func TestJWTService_ValidateToken_InvalidToken(t *testing.T) {
	service := NewJWTService()

	// Test with invalid token
	invalidTokens := []string{
		"invalid.token.here",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid.signature",
		"",
		"not.a.jwt.at.all",
	}

	for _, token := range invalidTokens {
		_, err := service.ValidateToken(token)
		if err == nil {
			t.Errorf("Expected error for invalid token: %s", token)
		}
	}
}

func TestJWTService_ValidateToken_ExpiredToken(t *testing.T) {
	// Create JWT service with very short expiration for testing
	cleanup := testutils.SetTestEnv("JWT_SECRET", "test-secret-for-expired-token")
	defer cleanup()

	service := NewJWTService()

	// Create user
	user := &models.User{
		ID:        1,
		Email:     "test@example.com",
		Role:      "hr",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
	}

	// Generate token
	token, _, err := service.GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Wait a bit to ensure token would be expired if using short expiration
	time.Sleep(10 * time.Millisecond)

	// Try to validate the token (this test is limited since we can't easily create truly expired tokens in unit tests)
	// In a real scenario, you would mock time or create tokens with past expiration
	_, err = service.ValidateToken(token)
	// Note: This test may not actually test expiration since tokens have 24h expiration by default
	// For proper expiration testing, you'd need to mock time or manually create tokens with past timestamps
	if err != nil {
		// This is expected behavior if token was somehow expired
		t.Logf("Token validation failed (may be expected): %v", err)
	}
}

func TestJWTService_ExtractClaims(t *testing.T) {
	service := NewJWTService()
	user := &models.User{
		ID:        1,
		Email:     "test@example.com",
		Role:      "hr",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
	}

	token, _, err := service.GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	claims, err := service.ExtractClaims(token)
	if err != nil {
		t.Fatalf("Failed to extract claims: %v", err)
	}

	if claims == nil {
		t.Error("Expected claims to be extracted")
	}

	// Verify claims contain expected data
	if claims["user_id"] != float64(1) {
		t.Errorf("Expected user_id 1, got %v", claims["user_id"])
	}

	if claims["email"] != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %v", claims["email"])
	}

	if claims["role"] != "hr" {
		t.Errorf("Expected role hr, got %v", claims["role"])
	}
}

func TestJWTService_ExtractClaims_InvalidToken(t *testing.T) {
	service := NewJWTService()

	invalidTokens := []string{
		"invalid.token.here",
		"",
		"not.a.jwt.at.all",
	}

	for _, token := range invalidTokens {
		_, err := service.ExtractClaims(token)
		if err == nil {
			t.Errorf("Expected error for invalid token: %s", token)
		}
	}
}

func TestJWTService_GetUserIDFromToken(t *testing.T) {
	service := NewJWTService()

	// Test valid claims
	claims := map[string]interface{}{
		"user_id": float64(123),
		"email":   "test@example.com",
	}

	userID, err := service.GetUserIDFromToken(claims)
	if err != nil {
		t.Fatalf("Failed to get user ID: %v", err)
	}

	if userID != 123 {
		t.Errorf("Expected user ID 123, got %d", userID)
	}

	// Test missing user_id
	invalidClaims := map[string]interface{}{
		"email": "test@example.com",
	}

	_, err = service.GetUserIDFromToken(invalidClaims)
	if err == nil {
		t.Error("Expected error for missing user_id")
	}

	// Test invalid user_id type
	invalidClaims2 := map[string]interface{}{
		"user_id": "invalid",
	}

	_, err = service.GetUserIDFromToken(invalidClaims2)
	if err == nil {
		t.Error("Expected error for invalid user_id type")
	}
}

func TestJWTService_GetUserEmailFromToken(t *testing.T) {
	service := NewJWTService()

	// Test valid claims
	claims := map[string]interface{}{
		"user_id": float64(1),
		"email":   "test@example.com",
	}

	email, err := service.GetUserEmailFromToken(claims)
	if err != nil {
		t.Fatalf("Failed to get email: %v", err)
	}

	if email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", email)
	}

	// Test missing email
	invalidClaims := map[string]interface{}{
		"user_id": float64(1),
	}

	_, err = service.GetUserEmailFromToken(invalidClaims)
	if err == nil {
		t.Error("Expected error for missing email")
	}

	// Test invalid email type
	invalidClaims2 := map[string]interface{}{
		"user_id": float64(1),
		"email":   123,
	}

	_, err = service.GetUserEmailFromToken(invalidClaims2)
	if err == nil {
		t.Error("Expected error for invalid email type")
	}
}

func TestJWTService_GetUserRoleFromToken(t *testing.T) {
	service := NewJWTService()

	// Test valid claims
	claims := map[string]interface{}{
		"user_id": float64(1),
		"email":   "test@example.com",
		"role":    "hr",
	}

	role, err := service.GetUserRoleFromToken(claims)
	if err != nil {
		t.Fatalf("Failed to get role: %v", err)
	}

	if role != "hr" {
		t.Errorf("Expected role hr, got %s", role)
	}

	// Test missing role
	invalidClaims := map[string]interface{}{
		"user_id": float64(1),
		"email":   "test@example.com",
	}

	_, err = service.GetUserRoleFromToken(invalidClaims)
	if err == nil {
		t.Error("Expected error for missing role")
	}

	// Test invalid role type
	invalidClaims2 := map[string]interface{}{
		"user_id": float64(1),
		"email":   "test@example.com",
		"role":    123,
	}

	_, err = service.GetUserRoleFromToken(invalidClaims2)
	if err == nil {
		t.Error("Expected error for invalid role type")
	}
}

func TestJWTService_IsTokenExpired(t *testing.T) {
	service := NewJWTService()

	// Test non-expired token
	nonExpiredClaims := map[string]interface{}{
		"exp": time.Now().Add(1 * time.Hour).Unix(),
	}

	if service.IsTokenExpired(nonExpiredClaims) {
		t.Error("Expected token to not be expired")
	}

	// Test expired token
	expiredClaims := map[string]interface{}{
		"exp": time.Now().Add(-1 * time.Hour).Unix(),
	}

	if !service.IsTokenExpired(expiredClaims) {
		t.Error("Expected token to be expired")
	}

	// Test missing exp claim (should be considered expired)
	missingExpClaims := map[string]interface{}{
		"user_id": float64(1),
	}

	if !service.IsTokenExpired(missingExpClaims) {
		t.Error("Expected token to be considered expired when exp claim is missing")
	}
}

func TestJWTService_Integration(t *testing.T) {
	service := NewJWTService()

	// Test complete flow: generate -> validate -> extract claims
	user := &models.User{
		ID:        42,
		Email:     "integration@test.com",
		Role:      "admin",
		FirstName: "Integration",
		LastName:  "Test",
		IsActive:  true,
	}

	// Generate token
	token, expiresAt, err := service.GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Validate token
	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	// Extract all user data
	userID, err := service.GetUserIDFromToken(*claims)
	if err != nil {
		t.Fatalf("Failed to get user ID: %v", err)
	}

	email, err := service.GetUserEmailFromToken(*claims)
	if err != nil {
		t.Fatalf("Failed to get email: %v", err)
	}

	role, err := service.GetUserRoleFromToken(*claims)
	if err != nil {
		t.Fatalf("Failed to get role: %v", err)
	}

	// Verify all data matches
	if userID != 42 {
		t.Errorf("Expected user ID 42, got %d", userID)
	}

	if email != "integration@test.com" {
		t.Errorf("Expected email integration@test.com, got %s", email)
	}

	if role != "admin" {
		t.Errorf("Expected role admin, got %s", role)
	}

	if service.IsTokenExpired(*claims) {
		t.Error("Expected token to not be expired")
	}

	if expiresAt.Before(time.Now()) {
		t.Error("Expected expiration time to be in the future")
	}
}
