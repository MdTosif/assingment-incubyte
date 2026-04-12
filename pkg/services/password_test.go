package services

import (
	"testing"
)

func TestNewPasswordService(t *testing.T) {
	service := NewPasswordService()
	if service == nil {
		t.Fatal("Expected password service to be created, got nil")
	}
}

func TestPasswordService_HashPassword(t *testing.T) {
	service := NewPasswordService()
	password := "testPassword123!"

	hashedPassword, err := service.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hashedPassword == "" {
		t.Error("Expected hashed password to be non-empty")
	}

	if hashedPassword == password {
		t.Error("Expected hashed password to be different from original password")
	}

	// Verify hash starts with bcrypt prefix (should be $2a$ or $2b$)
	if len(hashedPassword) < 60 || hashedPassword[0] != '$' {
		t.Error("Expected bcrypt hash format")
	}
}

func TestPasswordService_VerifyPassword(t *testing.T) {
	service := NewPasswordService()
	password := "testPassword123!"

	// Hash password
	hashedPassword, err := service.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Verify correct password
	err = service.VerifyPassword(hashedPassword, password)
	if err != nil {
		t.Fatalf("Failed to verify correct password: %v", err)
	}

	// Verify wrong password
	err = service.VerifyPassword(hashedPassword, "wrongPassword")
	if err == nil {
		t.Error("Expected error when verifying wrong password")
	}
}

func TestPasswordService_VerifyPassword_InvalidHash(t *testing.T) {
	service := NewPasswordService()

	invalidHashes := []string{
		"invalid_hash",
		"$2a$invalid",
		"$2b$12$invalidhashformat",
		"",
		"not_a_bcrypt_hash",
	}

	password := "testPassword123!"

	for _, hash := range invalidHashes {
		err := service.VerifyPassword(hash, password)
		if err == nil {
			t.Errorf("Expected error for invalid hash: %s", hash)
		}
	}
}

func TestPasswordService_ValidatePasswordStrength(t *testing.T) {
	service := NewPasswordService()

	// Test valid passwords
	validPasswords := []string{
		"SecurePass123!",
		"MyStrongP@ssw0rd",
		"Complex#Password2024",
		"Very$ecure123Password",
		"Str0ng!PassWord",
	}

	for _, password := range validPasswords {
		err := service.ValidatePasswordStrength(password)
		if err != nil {
			t.Errorf("Expected valid password to pass validation: %s, error: %v", password, err)
		}
	}
}

func TestPasswordService_ValidatePasswordStrength_WeakPasswords(t *testing.T) {
	service := NewPasswordService()

	// Test passwords that are too short
	shortPasswords := []string{
		"short",
		"1234567",
		"Abc123!",
		"Pass1!",
	}

	for _, password := range shortPasswords {
		err := service.ValidatePasswordStrength(password)
		if err == nil {
			t.Errorf("Expected short password to fail validation: %s", password)
		}
	}

	// Test passwords without uppercase
	noUppercase := []string{
		"lowercase123!",
		"alllowercase!",
		"nouppercase1",
	}

	for _, password := range noUppercase {
		err := service.ValidatePasswordStrength(password)
		if err == nil {
			t.Errorf("Expected password without uppercase to fail validation: %s", password)
		}
	}

	// Test passwords without lowercase
	noLowercase := []string{
		"UPPERCASE123!",
		"ALLUPPERCASE!",
		"NOLOWERCASE1",
	}

	for _, password := range noLowercase {
		err := service.ValidatePasswordStrength(password)
		if err == nil {
			t.Errorf("Expected password without lowercase to fail validation: %s", password)
		}
	}

	// Test passwords without numbers
	noNumbers := []string{
		"NoNumbersHere!",
		"OnlyLetters!",
		"PasswordOnly!",
	}

	for _, password := range noNumbers {
		err := service.ValidatePasswordStrength(password)
		if err == nil {
			t.Errorf("Expected password without numbers to fail validation: %s", password)
		}
	}

	// Test passwords without special characters
	noSpecial := []string{
		"NoSpecialChars123",
		"OnlyLetters123",
		"Password123",
	}

	for _, password := range noSpecial {
		err := service.ValidatePasswordStrength(password)
		if err == nil {
			t.Errorf("Expected password without special characters to fail validation: %s", password)
		}
	}
}

func TestPasswordService_ValidatePasswordStrength_CommonPasswords(t *testing.T) {
	service := NewPasswordService()

	// Test common passwords
	commonPasswords := []string{
		"password",
		"12345678",
		"123456789",
		"qwerty123",
		"abc12345",
		"password123",
		"admin123",
		"letmein123",
		"welcome123",
		"monkey123",
		"password1",
		"1234567890",
		"qwertyuiop",
		"asdfghjkl",
		"zxcvbnm",
		"123qweasd",
		"password12",
		"admin1234",
		"welcome1",
		"changeme",
	}

	for _, password := range commonPasswords {
		err := service.ValidatePasswordStrength(password)
		if err == nil {
			t.Errorf("Expected common password to fail validation: %s", password)
		}
	}
}

func TestPasswordService_ValidatePasswordStrength_SequentialPatterns(t *testing.T) {
	service := NewPasswordService()

	// Test sequential patterns
	sequentialPasswords := []string{
		"MyPassword12345678",
		"SecureQwerty123!",
		"PasswordAsdfgh123",
		"TestZxcvbnm456",
		"Sequential123456",
	}

	for _, password := range sequentialPasswords {
		err := service.ValidatePasswordStrength(password)
		if err == nil {
			t.Errorf("Expected password with sequential pattern to fail validation: %s", password)
		}
	}
}

func TestPasswordService_ValidatePasswordStrength_RepeatedPatterns(t *testing.T) {
	service := NewPasswordService()

	// Test repeated character patterns
	repeatedPasswords := []string{
		"Passwordaaaaaa",
		"Secure111111",
		"Test@@@@@@@",
		"Password!!!!!!",
		"Secure$$$$$$",
	}

	for _, password := range repeatedPasswords {
		err := service.ValidatePasswordStrength(password)
		if err == nil {
			t.Errorf("Expected password with repeated pattern to fail validation: %s", password)
		}
	}
}

func TestPasswordService_GenerateSecurePassword(t *testing.T) {
	service := NewPasswordService()

	// Test default length
	password := service.GenerateSecurePassword(0)
	if len(password) < 12 {
		t.Errorf("Expected password length at least 12, got %d", len(password))
	}

	// Test custom length
	password = service.GenerateSecurePassword(16)
	if len(password) != 16 {
		t.Errorf("Expected password length 16, got %d", len(password))
	}

	// Test short length (should default to 12)
	password = service.GenerateSecurePassword(6)
	if len(password) < 12 {
		t.Errorf("Expected password length at least 12 for short input, got %d", len(password))
	}

	// Verify generated password meets strength requirements
	err := service.ValidatePasswordStrength(password)
	if err != nil {
		t.Errorf("Generated password should meet strength requirements: %v", err)
	}
}

func TestPasswordService_Integration(t *testing.T) {
	service := NewPasswordService()
	password := "IntegrationTest123!"

	// Hash password
	hashedPassword, err := service.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Verify password
	err = service.VerifyPassword(hashedPassword, password)
	if err != nil {
		t.Fatalf("Failed to verify password: %v", err)
	}

	// Validate password strength
	err = service.ValidatePasswordStrength(password)
	if err != nil {
		t.Fatalf("Password should meet strength requirements: %v", err)
	}

	// Verify wrong password fails
	err = service.VerifyPassword(hashedPassword, "WrongPassword123!")
	if err == nil {
		t.Error("Expected verification to fail for wrong password")
	}
}

func TestPasswordService_EdgeCases(t *testing.T) {
	service := NewPasswordService()

	// Test empty password
	_, err := service.HashPassword("")
	if err != nil {
		t.Errorf("Should be able to hash empty password: %v", err)
	}

	// Test very long password
	longPassword := "VeryLongPassword123!WithManyCharactersToTestTheSystem"
	err = service.ValidatePasswordStrength(longPassword)
	if err != nil {
		t.Errorf("Long password should be valid: %v", err)
	}

	// Test password with special characters
	specialPassword := "P@ssw0rd!#$%^&*()_+-=[]{}|;:,.<>?"
	err = service.ValidatePasswordStrength(specialPassword)
	if err != nil {
		t.Errorf("Password with many special characters should be valid: %v", err)
	}

	// Test password with Unicode characters
	unicodePassword := "Pásswórd123!"
	err = service.ValidatePasswordStrength(unicodePassword)
	if err != nil {
		t.Logf("Unicode password validation result: %v", err)
	}
}
