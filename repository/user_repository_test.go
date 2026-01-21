package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/wabtcdi/user_service/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(
		&models.User{},
		&models.UserAuthentication{},
		&models.AccessLevel{},
		&models.UserAccessLevel{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgresUserRepository(db)
	ctx := context.Background()

	user := &models.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}
	phone := "555-1234"
	user.PhoneNumber = &phone

	auth := &models.UserAuthentication{
		PasswordHash: "hashedpassword123",
	}

	err := repo.Create(ctx, user, auth)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Verify user was created
	if user.ID == uuid.Nil {
		t.Error("User ID was not set")
	}
	if user.CreatedAt.IsZero() {
		t.Error("CreatedAt was not set")
	}
	if user.UpdatedAt.IsZero() {
		t.Error("UpdatedAt was not set")
	}

	// Verify authentication was created
	if auth.ID == uuid.Nil {
		t.Error("Auth ID was not set")
	}
	if auth.UserID != user.ID {
		t.Errorf("Auth UserID mismatch: got %v, want %v", auth.UserID, user.ID)
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgresUserRepository(db)
	ctx := context.Background()

	// Create a test user
	user := &models.User{
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane.smith@example.com",
	}
	auth := &models.UserAuthentication{
		PasswordHash: "hashedpassword456",
	}
	err := repo.Create(ctx, user, auth)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test GetByID
	retrieved, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to get user by ID: %v", err)
	}

	if retrieved.ID != user.ID {
		t.Errorf("ID mismatch: got %v, want %v", retrieved.ID, user.ID)
	}
	if retrieved.Email != user.Email {
		t.Errorf("Email mismatch: got %v, want %v", retrieved.Email, user.Email)
	}
	if retrieved.FirstName != user.FirstName {
		t.Errorf("FirstName mismatch: got %v, want %v", retrieved.FirstName, user.FirstName)
	}
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgresUserRepository(db)
	ctx := context.Background()

	// Try to get non-existent user
	nonExistentID := uuid.New()
	_, err := repo.GetByID(ctx, nonExistentID)
	if err == nil {
		t.Error("Expected error when getting non-existent user, got nil")
	}
	if err != nil && err.Error() != "user not found" {
		t.Errorf("Expected 'user not found' error, got: %v", err)
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgresUserRepository(db)
	ctx := context.Background()

	// Create a test user
	user := &models.User{
		FirstName: "Bob",
		LastName:  "Johnson",
		Email:     "bob.johnson@example.com",
	}
	auth := &models.UserAuthentication{
		PasswordHash: "hashedpassword789",
	}
	err := repo.Create(ctx, user, auth)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test GetByEmail
	retrieved, err := repo.GetByEmail(ctx, user.Email)
	if err != nil {
		t.Fatalf("Failed to get user by email: %v", err)
	}

	if retrieved.ID != user.ID {
		t.Errorf("ID mismatch: got %v, want %v", retrieved.ID, user.ID)
	}
	if retrieved.Email != user.Email {
		t.Errorf("Email mismatch: got %v, want %v", retrieved.Email, user.Email)
	}
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgresUserRepository(db)
	ctx := context.Background()

	// Create a test user
	user := &models.User{
		FirstName: "Alice",
		LastName:  "Brown",
		Email:     "alice.brown@example.com",
	}
	auth := &models.UserAuthentication{
		PasswordHash: "hashedpassword000",
	}
	err := repo.Create(ctx, user, auth)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Update the user
	user.FirstName = "Alicia"
	user.Email = "alicia.brown@example.com"
	newPhone := "555-9999"
	user.PhoneNumber = &newPhone

	err = repo.Update(ctx, user)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	// Verify update
	retrieved, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to get updated user: %v", err)
	}

	if retrieved.FirstName != "Alicia" {
		t.Errorf("FirstName not updated: got %v, want Alicia", retrieved.FirstName)
	}
	if retrieved.Email != "alicia.brown@example.com" {
		t.Errorf("Email not updated: got %v, want alicia.brown@example.com", retrieved.Email)
	}
	if retrieved.PhoneNumber == nil || *retrieved.PhoneNumber != "555-9999" {
		t.Errorf("PhoneNumber not updated correctly")
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgresUserRepository(db)
	ctx := context.Background()

	// Create a test user
	user := &models.User{
		FirstName: "Charlie",
		LastName:  "Davis",
		Email:     "charlie.davis@example.com",
	}
	auth := &models.UserAuthentication{
		PasswordHash: "hashedpassword111",
	}
	err := repo.Create(ctx, user, auth)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Delete the user (soft delete)
	err = repo.Delete(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	// Verify user is soft-deleted (not found in normal queries)
	_, err = repo.GetByID(ctx, user.ID)
	if err == nil {
		t.Error("Expected error when getting deleted user, got nil")
	}
}

func TestUserRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgresUserRepository(db)
	ctx := context.Background()

	// Create multiple test users
	for i := 0; i < 5; i++ {
		user := &models.User{
			FirstName: "User",
			LastName:  "Test",
			Email:     "user" + string(rune('0'+i)) + "@example.com",
		}
		auth := &models.UserAuthentication{
			PasswordHash: "hashedpassword",
		}
		err := repo.Create(ctx, user, auth)
		if err != nil {
			t.Fatalf("Failed to create test user %d: %v", i, err)
		}
		// Sleep to ensure different CreatedAt timestamps
		time.Sleep(10 * time.Millisecond)
	}

	// Test List with pagination
	users, total, err := repo.List(ctx, 3, 0)
	if err != nil {
		t.Fatalf("Failed to list users: %v", err)
	}

	if total != 5 {
		t.Errorf("Total count mismatch: got %d, want 5", total)
	}
	if len(users) != 3 {
		t.Errorf("Users count mismatch: got %d, want 3", len(users))
	}

	// Test pagination offset
	users, total, err = repo.List(ctx, 3, 3)
	if err != nil {
		t.Fatalf("Failed to list users with offset: %v", err)
	}

	if total != 5 {
		t.Errorf("Total count mismatch with offset: got %d, want 5", total)
	}
	if len(users) != 2 {
		t.Errorf("Users count mismatch with offset: got %d, want 2", len(users))
	}
}

func TestUserRepository_GetUserAuthentication(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgresUserRepository(db)
	ctx := context.Background()

	// Create a test user
	user := &models.User{
		FirstName: "David",
		LastName:  "Miller",
		Email:     "david.miller@example.com",
	}
	auth := &models.UserAuthentication{
		PasswordHash: "hashedpassword222",
	}
	err := repo.Create(ctx, user, auth)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Get authentication
	retrievedAuth, err := repo.GetUserAuthentication(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to get user authentication: %v", err)
	}

	if retrievedAuth.UserID != user.ID {
		t.Errorf("UserID mismatch: got %v, want %v", retrievedAuth.UserID, user.ID)
	}
	if retrievedAuth.PasswordHash != auth.PasswordHash {
		t.Errorf("PasswordHash mismatch: got %v, want %v", retrievedAuth.PasswordHash, auth.PasswordHash)
	}
}

func TestUserRepository_GetByIDWithAccessLevels(t *testing.T) {
	db := setupTestDB(t)
	userRepo := NewPostgresUserRepository(db)
	accessLevelRepo := NewPostgresAccessLevelRepository(db)
	ctx := context.Background()

	// Create a test user
	user := &models.User{
		FirstName: "Emma",
		LastName:  "Wilson",
		Email:     "emma.wilson@example.com",
	}
	auth := &models.UserAuthentication{
		PasswordHash: "hashedpassword333",
	}
	err := userRepo.Create(ctx, user, auth)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create test access levels
	adminLevel := &models.AccessLevel{
		Name: "Admin",
	}
	desc := "Administrator access"
	adminLevel.Description = &desc
	err = accessLevelRepo.Create(ctx, adminLevel)
	if err != nil {
		t.Fatalf("Failed to create admin access level: %v", err)
	}

	userLevel := &models.AccessLevel{
		Name: "User",
	}
	err = accessLevelRepo.Create(ctx, userLevel)
	if err != nil {
		t.Fatalf("Failed to create user access level: %v", err)
	}

	// Assign access levels to user
	err = accessLevelRepo.AssignToUser(ctx, user.ID, adminLevel.ID)
	if err != nil {
		t.Fatalf("Failed to assign admin level: %v", err)
	}
	err = accessLevelRepo.AssignToUser(ctx, user.ID, userLevel.ID)
	if err != nil {
		t.Fatalf("Failed to assign user level: %v", err)
	}

	// Test GetByIDWithAccessLevels
	retrievedUser, accessLevels, err := userRepo.GetByIDWithAccessLevels(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to get user with access levels: %v", err)
	}

	if retrievedUser.ID != user.ID {
		t.Errorf("User ID mismatch: got %v, want %v", retrievedUser.ID, user.ID)
	}
	if len(accessLevels) != 2 {
		t.Errorf("Access levels count mismatch: got %d, want 2", len(accessLevels))
	}

	// Verify access levels are preloaded correctly
	foundAdmin := false
	foundUser := false
	for _, al := range accessLevels {
		if al.Name == "Admin" {
			foundAdmin = true
		}
		if al.Name == "User" {
			foundUser = true
		}
	}
	if !foundAdmin {
		t.Error("Admin access level not found in preloaded data")
	}
	if !foundUser {
		t.Error("User access level not found in preloaded data")
	}
}
