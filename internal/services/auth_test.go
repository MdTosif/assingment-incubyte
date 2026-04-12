package services

import (
	"testing"

	"github.com/tofiquem/assingment/internal/database"
	"github.com/tofiquem/assingment/internal/models"
	"github.com/tofiquem/assingment/internal/testutils"
	"gorm.io/gorm"
)

func TestNewAuthService(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewAuthService(testDB)
	if service == nil {
		t.Fatal("Expected auth service to be created, got nil")
	}

	if service.db != testDB {
		t.Error("Expected service.db to be testDB")
	}

	if service.jwtService == nil {
		t.Error("Expected jwtService to be initialized")
	}

	if service.passwordService == nil {
		t.Error("Expected passwordService to be initialized")
	}
}

func TestAuthService_Login(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewAuthService(testDB)

	// Create test user
	user := &models.User{
		Email:     "test@example.com",
		Password:  "hashedPassword",
		Role:      "hr",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
	}

	if err := testDB.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Mock password verification by updating user with properly hashed password
	hashedPassword, err := service.passwordService.HashPassword("testPassword123!")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	user.Password = hashedPassword
	if err := testDB.Save(user).Error; err != nil {
		t.Fatalf("Failed to update user password: %v", err)
	}

	// Test successful login
	response, err := service.Login("test@example.com", "testPassword123!")
	if err != nil {
		t.Fatalf("Failed to login with valid credentials: %v", err)
	}

	if response.Token == "" {
		t.Error("Expected token to be returned")
	}

	if response.User.Email != "test@example.com" {
		t.Errorf("Expected user email test@example.com, got %s", response.User.Email)
	}

	if response.User.Role != "hr" {
		t.Errorf("Expected user role hr, got %s", response.User.Role)
	}

	if response.User.Password != "" {
		t.Error("Expected password to be excluded from response")
	}
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewAuthService(testDB)

	// Create test user
	user := &models.User{
		Email:     "test@example.com",
		Password:  "hashedPassword",
		Role:      "hr",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
	}

	if err := testDB.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test login with wrong email
	_, err := service.Login("wrong@example.com", "testPassword123!")
	if err == nil {
		t.Error("Expected error for wrong email")
	}

	// Test login with wrong password
	_, err = service.Login("test@example.com", "wrongPassword")
	if err == nil {
		t.Error("Expected error for wrong password")
	}
}

func TestAuthService_Login_InactiveUser(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewAuthService(testDB)

	// Create inactive test user
	user := &models.User{
		Email:     "test@example.com",
		Password:  "hashedPassword",
		Role:      "hr",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  false,
	}

	if err := testDB.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test login with inactive user
	_, err := service.Login("test@example.com", "testPassword123!")
	if err == nil {
		t.Error("Expected error for inactive user")
	}
}

func TestAuthService_Login_NonExistentUser(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewAuthService(testDB)

	// Test login with non-existent user
	_, err := service.Login("nonexistent@example.com", "testPassword123!")
	if err == nil {
		t.Error("Expected error for non-existent user")
	}
}

func TestAuthService_CreateUser(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewAuthService(testDB)

	// Test creating HR user
	req := &models.CreateHRUserRequest{
		Email:     "newuser@example.com",
		Password:  "SecurePass123!",
		FirstName: "New",
		LastName:  "User",
		Role:      "hr",
	}

	user, err := service.CreateUser(req)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.ID == 0 {
		t.Error("Expected user ID to be set")
	}

	if user.Email != "newuser@example.com" {
		t.Errorf("Expected email newuser@example.com, got %s", user.Email)
	}

	if user.Role != "hr" {
		t.Errorf("Expected role hr, got %s", user.Role)
	}

	if !user.IsActive {
		t.Error("Expected user to be active")
	}

	if user.Password == "" {
		t.Error("Expected password to be hashed")
	}

	if user.Password == "SecurePass123!" {
		t.Error("Expected password to be hashed, not plain text")
	}
}

func TestAuthService_CreateUser_DuplicateEmail(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewAuthService(testDB)

	// Create first user
	req1 := &models.CreateHRUserRequest{
		Email:     "duplicate@example.com",
		Password:  "SecurePass123!",
		FirstName: "First",
		LastName:  "User",
		Role:      "hr",
	}

	_, err := service.CreateUser(req1)
	if err != nil {
		t.Fatalf("Failed to create first user: %v", err)
	}

	// Try to create user with same email
	req2 := &models.CreateHRUserRequest{
		Email:     "duplicate@example.com",
		Password:  "SecurePass123!",
		FirstName: "Second",
		LastName:  "User",
		Role:      "hr",
	}

	_, err = service.CreateUser(req2)
	if err == nil {
		t.Error("Expected error for duplicate email")
	}
}

func TestAuthService_CreateUser_WeakPassword(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewAuthService(testDB)

	// Test creating user with weak password
	weakPasswords := []string{
		"weak",
		"password",
		"12345678",
		"nouppercase1!",
		"NOLOWERCASE1!",
		"nonumberhere!",
		"NoSpecialChars123",
	}

	for _, password := range weakPasswords {
		req := &models.CreateHRUserRequest{
			Email:     "weak@example.com",
			Password:  password,
			FirstName: "Weak",
			LastName:  "Password",
			Role:      "hr",
		}

		_, err := service.CreateUser(req)
		if err == nil {
			t.Errorf("Expected error for weak password: %s", password)
		}
	}
}

func TestAuthService_GetUserByID(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewAuthService(testDB)

	// Create test user
	user := &models.User{
		Email:     "test@example.com",
		Password:  "hashedPassword",
		Role:      "hr",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
	}

	if err := testDB.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test getting user by ID
	foundUser, err := service.GetUserByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to get user by ID: %v", err)
	}

	if foundUser.ID != user.ID {
		t.Errorf("Expected user ID %d, got %d", user.ID, foundUser.ID)
	}

	if foundUser.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", foundUser.Email)
	}
}

func TestAuthService_GetUserByID_NonExistent(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewAuthService(testDB)

	// Test getting non-existent user
	_, err := service.GetUserByID(999)
	if err == nil {
		t.Error("Expected error for non-existent user")
	}
}

func TestAuthService_GetUserByEmail(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewAuthService(testDB)

	// Create test user
	user := &models.User{
		Email:     "test@example.com",
		Password:  "hashedPassword",
		Role:      "hr",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
	}

	if err := testDB.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test getting user by email
	foundUser, err := service.GetUserByEmail("test@example.com")
	if err != nil {
		t.Fatalf("Failed to get user by email: %v", err)
	}

	if foundUser.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", foundUser.Email)
	}

	if foundUser.ID != user.ID {
		t.Errorf("Expected user ID %d, got %d", user.ID, foundUser.ID)
	}
}

func TestAuthService_GetUserByEmail_NonExistent(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewAuthService(testDB)

	// Test getting non-existent user
	_, err := service.GetUserByEmail("nonexistent@example.com")
	if err == nil {
		t.Error("Expected error for non-existent user")
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewAuthService(testDB)

	// Create test user
	user := &models.User{
		Email:     "test@example.com",
		Password:  "hashedPassword",
		Role:      "hr",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
	}

	if err := testDB.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Generate valid token
	token, _, err := service.jwtService.GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Validate token
	validatedUser, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if validatedUser.ID != user.ID {
		t.Errorf("Expected user ID %d, got %d", user.ID, validatedUser.ID)
	}

	if validatedUser.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", validatedUser.Email)
	}
}

func TestAuthService_ValidateToken_InvalidToken(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewAuthService(testDB)

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

func TestAuthService_ValidateToken_InactiveUser(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewAuthService(testDB)

	// Create test user first (will be active by default)
	user := &models.User{
		Email:     "test@example.com",
		Password:  "hashedPassword",
		Role:      "hr",
		FirstName: "Test",
		LastName:  "User",
	}

	if err := testDB.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Then update user to be inactive
	if err := testDB.Model(user).Update("is_active", false).Error; err != nil {
		t.Fatalf("Failed to update user to inactive: %v", err)
	}

	// Refresh user from database to get updated values
	if err := testDB.First(user, user.ID).Error; err != nil {
		t.Fatalf("Failed to refresh user: %v", err)
	}

	// Verify user is actually inactive
	if user.IsActive {
		t.Fatalf("Expected user to be inactive, but IsActive is %v", user.IsActive)
	}

	// Generate token for inactive user
	token, _, err := service.jwtService.GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Try to validate token for inactive user
	_, err = service.ValidateToken(token)
	if err == nil {
		t.Error("Expected error for inactive user")
	}
}

func TestAuthService_UpdateUser(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewAuthService(testDB)

	// Create test user
	user := &models.User{
		Email:     "test@example.com",
		Password:  "hashedPassword",
		Role:      "hr",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
	}

	if err := testDB.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Update user
	newFirstName := "Updated"
	newRole := "admin"
	req := &models.UpdateUserRequest{
		FirstName: &newFirstName,
		Role:      &newRole,
	}

	updatedUser, err := service.UpdateUser(user.ID, req)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	if updatedUser.FirstName != "Updated" {
		t.Errorf("Expected first name Updated, got %s", updatedUser.FirstName)
	}

	if updatedUser.Role != "admin" {
		t.Errorf("Expected role admin, got %s", updatedUser.Role)
	}

	// Verify unchanged fields
	if updatedUser.Email != "test@example.com" {
		t.Errorf("Expected email to remain test@example.com, got %s", updatedUser.Email)
	}
}

func TestAuthService_DeleteUser(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewAuthService(testDB)

	// Create test user
	user := &models.User{
		Email:     "test@example.com",
		Password:  "hashedPassword",
		Role:      "hr",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
	}

	if err := testDB.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Delete user
	err := service.DeleteUser(user.ID)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	// Verify user is deleted
	var deletedUser models.User
	err = testDB.First(&deletedUser, user.ID).Error
	if err == nil {
		t.Error("Expected user to be deleted")
	}

	if err != gorm.ErrRecordNotFound {
		t.Errorf("Expected record not found error, got: %v", err)
	}
}

func TestAuthService_ListUsers(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewAuthService(testDB)

	// Create test users
	users := []models.User{
		{Email: "user1@example.com", Password: "hashed1", Role: "hr", FirstName: "User", LastName: "One", IsActive: true},
		{Email: "user2@example.com", Password: "hashed2", Role: "admin", FirstName: "User", LastName: "Two", IsActive: true},
		{Email: "user3@example.com", Password: "hashed3", Role: "hr", FirstName: "User", LastName: "Three", IsActive: true},
	}

	for _, user := range users {
		if err := testDB.Create(&user).Error; err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}

	// List users
	userList, total, err := service.ListUsers(1, 10)
	if err != nil {
		t.Fatalf("Failed to list users: %v", err)
	}

	if len(userList) != 3 {
		t.Errorf("Expected 3 users, got %d", len(userList))
	}

	if total != 3 {
		t.Errorf("Expected total 3, got %d", total)
	}

	// Verify passwords are excluded
	for _, user := range userList {
		if user.Password != "" {
			t.Error("Expected password to be excluded from user list")
		}
	}
}

func TestAuthService_ChangePassword(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewAuthService(testDB)

	// Create test user
	hashedPassword, err := service.passwordService.HashPassword("oldPassword123!")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	user := &models.User{
		Email:     "test@example.com",
		Password:  hashedPassword,
		Role:      "hr",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
	}

	if err := testDB.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Change password
	err = service.ChangePassword(user.ID, "oldPassword123!", "newPassword456!")
	if err != nil {
		t.Fatalf("Failed to change password: %v", err)
	}

	// Verify password was changed
	updatedUser, err := service.GetUserByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to get updated user: %v", err)
	}

	// Verify old password doesn't work
	err = service.passwordService.VerifyPassword(updatedUser.Password, "oldPassword123!")
	if err == nil {
		t.Error("Expected old password to be invalid")
	}

	// Verify new password works
	err = service.passwordService.VerifyPassword(updatedUser.Password, "newPassword456!")
	if err != nil {
		t.Errorf("Expected new password to be valid: %v", err)
	}
}
