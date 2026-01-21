package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/wabtcdi/user_service/models"
)

func TestAccessLevelRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgresAccessLevelRepository(db)
	ctx := context.Background()

	accessLevel := &models.AccessLevel{
		Name: "Manager",
	}
	desc := "Manager access level"
	accessLevel.Description = &desc

	err := repo.Create(ctx, accessLevel)
	if err != nil {
		t.Fatalf("Failed to create access level: %v", err)
	}

	if accessLevel.ID == 0 {
		t.Error("AccessLevel ID was not set")
	}
	if accessLevel.CreatedAt.IsZero() {
		t.Error("CreatedAt was not set")
	}
}

func TestAccessLevelRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgresAccessLevelRepository(db)
	ctx := context.Background()

	// Create a test access level
	accessLevel := &models.AccessLevel{
		Name: "Editor",
	}
	err := repo.Create(ctx, accessLevel)
	if err != nil {
		t.Fatalf("Failed to create test access level: %v", err)
	}

	// Test GetByID
	retrieved, err := repo.GetByID(ctx, accessLevel.ID)
	if err != nil {
		t.Fatalf("Failed to get access level by ID: %v", err)
	}

	if retrieved.ID != accessLevel.ID {
		t.Errorf("ID mismatch: got %d, want %d", retrieved.ID, accessLevel.ID)
	}
	if retrieved.Name != accessLevel.Name {
		t.Errorf("Name mismatch: got %s, want %s", retrieved.Name, accessLevel.Name)
	}
}

func TestAccessLevelRepository_GetByName(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgresAccessLevelRepository(db)
	ctx := context.Background()

	// Create a test access level
	accessLevel := &models.AccessLevel{
		Name: "Viewer",
	}
	err := repo.Create(ctx, accessLevel)
	if err != nil {
		t.Fatalf("Failed to create test access level: %v", err)
	}

	// Test GetByName
	retrieved, err := repo.GetByName(ctx, "Viewer")
	if err != nil {
		t.Fatalf("Failed to get access level by name: %v", err)
	}

	if retrieved.Name != "Viewer" {
		t.Errorf("Name mismatch: got %s, want Viewer", retrieved.Name)
	}
}

func TestAccessLevelRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgresAccessLevelRepository(db)
	ctx := context.Background()

	// Create multiple access levels
	levels := []string{"Admin", "Editor", "Viewer", "Guest"}
	for _, name := range levels {
		accessLevel := &models.AccessLevel{
			Name: name,
		}
		err := repo.Create(ctx, accessLevel)
		if err != nil {
			t.Fatalf("Failed to create access level %s: %v", name, err)
		}
	}

	// Test List
	retrieved, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("Failed to list access levels: %v", err)
	}

	if len(retrieved) != 4 {
		t.Errorf("Count mismatch: got %d, want 4", len(retrieved))
	}

	// Verify sorted by name
	if len(retrieved) > 0 && retrieved[0].Name != "Admin" {
		t.Errorf("First item should be 'Admin' (sorted), got %s", retrieved[0].Name)
	}
}

func TestAccessLevelRepository_AssignToUser(t *testing.T) {
	db := setupTestDB(t)
	accessLevelRepo := NewPostgresAccessLevelRepository(db)
	userRepo := NewPostgresUserRepository(db)
	ctx := context.Background()

	// Create a test user
	user := &models.User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "test.user@example.com",
	}
	auth := &models.UserAuthentication{
		PasswordHash: "hashedpassword",
	}
	err := userRepo.Create(ctx, user, auth)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create an access level
	accessLevel := &models.AccessLevel{
		Name: "TestLevel",
	}
	err = accessLevelRepo.Create(ctx, accessLevel)
	if err != nil {
		t.Fatalf("Failed to create access level: %v", err)
	}

	// Assign to user
	err = accessLevelRepo.AssignToUser(ctx, user.ID, accessLevel.ID)
	if err != nil {
		t.Fatalf("Failed to assign access level to user: %v", err)
	}

	// Verify assignment
	userAccessLevels, err := accessLevelRepo.GetUserAccessLevels(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to get user access levels: %v", err)
	}

	if len(userAccessLevels) != 1 {
		t.Errorf("Access levels count mismatch: got %d, want 1", len(userAccessLevels))
	}
	if len(userAccessLevels) > 0 && userAccessLevels[0].Name != "TestLevel" {
		t.Errorf("Access level name mismatch: got %s, want TestLevel", userAccessLevels[0].Name)
	}
}

func TestAccessLevelRepository_AssignToUser_Duplicate(t *testing.T) {
	db := setupTestDB(t)
	accessLevelRepo := NewPostgresAccessLevelRepository(db)
	userRepo := NewPostgresUserRepository(db)
	ctx := context.Background()

	// Create a test user
	user := &models.User{
		FirstName: "Test",
		LastName:  "User2",
		Email:     "test.user2@example.com",
	}
	auth := &models.UserAuthentication{
		PasswordHash: "hashedpassword",
	}
	err := userRepo.Create(ctx, user, auth)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create an access level
	accessLevel := &models.AccessLevel{
		Name: "DuplicateTest",
	}
	err = accessLevelRepo.Create(ctx, accessLevel)
	if err != nil {
		t.Fatalf("Failed to create access level: %v", err)
	}

	// Assign to user first time
	err = accessLevelRepo.AssignToUser(ctx, user.ID, accessLevel.ID)
	if err != nil {
		t.Fatalf("Failed to assign access level to user: %v", err)
	}

	// Try to assign again (should use upsert logic)
	err = accessLevelRepo.AssignToUser(ctx, user.ID, accessLevel.ID)
	if err != nil {
		t.Fatalf("Failed to assign duplicate access level: %v", err)
	}

	// Verify still only one assignment
	userAccessLevels, err := accessLevelRepo.GetUserAccessLevels(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to get user access levels: %v", err)
	}

	if len(userAccessLevels) != 1 {
		t.Errorf("Access levels count mismatch after duplicate assign: got %d, want 1", len(userAccessLevels))
	}
}

func TestAccessLevelRepository_RemoveFromUser(t *testing.T) {
	db := setupTestDB(t)
	accessLevelRepo := NewPostgresAccessLevelRepository(db)
	userRepo := NewPostgresUserRepository(db)
	ctx := context.Background()

	// Create a test user
	user := &models.User{
		FirstName: "Remove",
		LastName:  "Test",
		Email:     "remove.test@example.com",
	}
	auth := &models.UserAuthentication{
		PasswordHash: "hashedpassword",
	}
	err := userRepo.Create(ctx, user, auth)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create an access level
	accessLevel := &models.AccessLevel{
		Name: "RemoveLevel",
	}
	err = accessLevelRepo.Create(ctx, accessLevel)
	if err != nil {
		t.Fatalf("Failed to create access level: %v", err)
	}

	// Assign to user
	err = accessLevelRepo.AssignToUser(ctx, user.ID, accessLevel.ID)
	if err != nil {
		t.Fatalf("Failed to assign access level: %v", err)
	}

	// Remove from user
	err = accessLevelRepo.RemoveFromUser(ctx, user.ID, accessLevel.ID)
	if err != nil {
		t.Fatalf("Failed to remove access level from user: %v", err)
	}

	// Verify removal (soft delete)
	userAccessLevels, err := accessLevelRepo.GetUserAccessLevels(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to get user access levels: %v", err)
	}

	if len(userAccessLevels) != 0 {
		t.Errorf("Access levels should be empty after removal: got %d", len(userAccessLevels))
	}
}

func TestAccessLevelRepository_RemoveFromUser_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostgresAccessLevelRepository(db)
	ctx := context.Background()

	// Try to remove non-existent assignment
	err := repo.RemoveFromUser(ctx, uuid.New(), 999)
	if err == nil {
		t.Error("Expected error when removing non-existent access level assignment, got nil")
	}
}

func TestAccessLevelRepository_GetUserAccessLevels_Empty(t *testing.T) {
	db := setupTestDB(t)
	accessLevelRepo := NewPostgresAccessLevelRepository(db)
	userRepo := NewPostgresUserRepository(db)
	ctx := context.Background()

	// Create a user without access levels
	user := &models.User{
		FirstName: "Empty",
		LastName:  "Levels",
		Email:     "empty.levels@example.com",
	}
	auth := &models.UserAuthentication{
		PasswordHash: "hashedpassword",
	}
	err := userRepo.Create(ctx, user, auth)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Get access levels for user with none assigned
	accessLevels, err := accessLevelRepo.GetUserAccessLevels(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to get user access levels: %v", err)
	}

	if len(accessLevels) != 0 {
		t.Errorf("Expected empty access levels, got %d", len(accessLevels))
	}
}

func TestAccessLevelRepository_GetUserAccessLevels_Multiple(t *testing.T) {
	db := setupTestDB(t)
	accessLevelRepo := NewPostgresAccessLevelRepository(db)
	userRepo := NewPostgresUserRepository(db)
	ctx := context.Background()

	// Create a test user
	user := &models.User{
		FirstName: "Multi",
		LastName:  "Levels",
		Email:     "multi.levels@example.com",
	}
	auth := &models.UserAuthentication{
		PasswordHash: "hashedpassword",
	}
	err := userRepo.Create(ctx, user, auth)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create multiple access levels
	levelNames := []string{"Alpha", "Beta", "Gamma"}
	for _, name := range levelNames {
		accessLevel := &models.AccessLevel{
			Name: name,
		}
		err = accessLevelRepo.Create(ctx, accessLevel)
		if err != nil {
			t.Fatalf("Failed to create access level %s: %v", name, err)
		}

		// Assign to user
		err = accessLevelRepo.AssignToUser(ctx, user.ID, accessLevel.ID)
		if err != nil {
			t.Fatalf("Failed to assign access level %s: %v", name, err)
		}
	}

	// Get all access levels for user
	accessLevels, err := accessLevelRepo.GetUserAccessLevels(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to get user access levels: %v", err)
	}

	if len(accessLevels) != 3 {
		t.Errorf("Access levels count mismatch: got %d, want 3", len(accessLevels))
	}

	// Verify they're sorted by name
	if len(accessLevels) > 0 && accessLevels[0].Name != "Alpha" {
		t.Errorf("First access level should be 'Alpha' (sorted), got %s", accessLevels[0].Name)
	}
}
