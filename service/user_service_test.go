package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/wabtcdi/user_service/dto"
	"github.com/wabtcdi/user_service/mocks"
	"github.com/wabtcdi/user_service/models"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockAccessLevelRepo := mocks.NewMockAccessLevelRepository(ctrl)
	service := NewUserService(mockUserRepo, mockAccessLevelRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		req := &dto.CreateUserRequest{
			FirstName:   "John",
			LastName:    "Doe",
			Email:       "john.doe@example.com",
			PhoneNumber: "123-456-7890",
			Password:    "password123",
		}

		mockUserRepo.EXPECT().
			GetByEmail(ctx, req.Email).
			Return(nil, errors.New("not found"))

		mockUserRepo.EXPECT().
			Create(ctx, gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, user *models.User, auth *models.UserAuthentication) error {
				user.ID = uuid.New()
				user.CreatedAt = time.Now()
				user.UpdatedAt = time.Now()
				return nil
			})

		mockAccessLevelRepo.EXPECT().
			GetUserAccessLevels(ctx, gomock.Any()).
			Return([]*models.AccessLevel{}, nil)

		resp, err := service.CreateUser(ctx, req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.FirstName != req.FirstName {
			t.Errorf("Expected FirstName %s, got %s", req.FirstName, resp.FirstName)
		}
		if resp.Email != req.Email {
			t.Errorf("Expected Email %s, got %s", req.Email, resp.Email)
		}
	})

	t.Run("UserAlreadyExists", func(t *testing.T) {
		req := &dto.CreateUserRequest{
			FirstName: "Jane",
			LastName:  "Smith",
			Email:     "jane.smith@example.com",
			Password:  "password123",
		}

		existingUser := &models.User{
			ID:    uuid.New(),
			Email: req.Email,
		}

		mockUserRepo.EXPECT().
			GetByEmail(ctx, req.Email).
			Return(existingUser, nil)

		_, err := service.CreateUser(ctx, req)
		if err == nil {
			t.Fatal("Expected error for duplicate user, got nil")
		}
	})

	t.Run("ValidationError_EmptyFirstName", func(t *testing.T) {
		req := &dto.CreateUserRequest{
			FirstName: "",
			LastName:  "Doe",
			Email:     "test@example.com",
			Password:  "password123",
		}

		_, err := service.CreateUser(ctx, req)
		if err == nil {
			t.Fatal("Expected validation error, got nil")
		}
	})

	t.Run("ValidationError_EmptyLastName", func(t *testing.T) {
		req := &dto.CreateUserRequest{
			FirstName: "John",
			LastName:  "",
			Email:     "test@example.com",
			Password:  "password123",
		}

		_, err := service.CreateUser(ctx, req)
		if err == nil {
			t.Fatal("Expected validation error, got nil")
		}
	})

	t.Run("ValidationError_InvalidEmail", func(t *testing.T) {
		req := &dto.CreateUserRequest{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "invalid-email",
			Password:  "password123",
		}

		_, err := service.CreateUser(ctx, req)
		if err == nil {
			t.Fatal("Expected validation error, got nil")
		}
	})

	t.Run("ValidationError_ShortPassword", func(t *testing.T) {
		req := &dto.CreateUserRequest{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "test@example.com",
			Password:  "short",
		}

		_, err := service.CreateUser(ctx, req)
		if err == nil {
			t.Fatal("Expected validation error, got nil")
		}
	})
}

func TestUserService_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockAccessLevelRepo := mocks.NewMockAccessLevelRepository(ctrl)
	service := NewUserService(mockUserRepo, mockAccessLevelRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()
		user := &models.User{
			ID:        userID,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockUserRepo.EXPECT().
			GetByID(ctx, userID).
			Return(user, nil)

		mockAccessLevelRepo.EXPECT().
			GetUserAccessLevels(ctx, userID).
			Return([]*models.AccessLevel{}, nil)

		resp, err := service.GetUser(ctx, userID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.ID != userID {
			t.Errorf("Expected ID %v, got %v", userID, resp.ID)
		}
		if resp.Email != user.Email {
			t.Errorf("Expected Email %s, got %s", user.Email, resp.Email)
		}
	})

	t.Run("UserNotFound", func(t *testing.T) {
		userID := uuid.New()

		mockUserRepo.EXPECT().
			GetByID(ctx, userID).
			Return(nil, errors.New("user not found"))

		_, err := service.GetUser(ctx, userID)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockAccessLevelRepo := mocks.NewMockAccessLevelRepository(ctrl)
	service := NewUserService(mockUserRepo, mockAccessLevelRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()
		user := &models.User{
			ID:        userID,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
		}

		req := &dto.UpdateUserRequest{
			FirstName: "Jane",
			LastName:  "Smith",
		}

		mockUserRepo.EXPECT().
			GetByID(ctx, userID).
			Return(user, nil)

		mockUserRepo.EXPECT().
			Update(ctx, gomock.Any()).
			Return(nil)

		mockAccessLevelRepo.EXPECT().
			GetUserAccessLevels(ctx, userID).
			Return([]*models.AccessLevel{}, nil)

		resp, err := service.UpdateUser(ctx, userID, req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.FirstName != req.FirstName {
			t.Errorf("Expected FirstName %s, got %s", req.FirstName, resp.FirstName)
		}
	})

	t.Run("EmailAlreadyTaken", func(t *testing.T) {
		userID := uuid.New()
		otherUserID := uuid.New()

		user := &models.User{
			ID:    userID,
			Email: "john@example.com",
		}

		existingUser := &models.User{
			ID:    otherUserID,
			Email: "jane@example.com",
		}

		req := &dto.UpdateUserRequest{
			Email: "jane@example.com",
		}

		mockUserRepo.EXPECT().
			GetByID(ctx, userID).
			Return(user, nil)

		mockUserRepo.EXPECT().
			GetByEmail(ctx, req.Email).
			Return(existingUser, nil)

		_, err := service.UpdateUser(ctx, userID, req)
		if err == nil {
			t.Fatal("Expected error for email already taken, got nil")
		}
	})

	t.Run("UserNotFound", func(t *testing.T) {
		userID := uuid.New()
		req := &dto.UpdateUserRequest{
			FirstName: "John",
		}

		mockUserRepo.EXPECT().
			GetByID(ctx, userID).
			Return(nil, errors.New("user not found"))

		_, err := service.UpdateUser(ctx, userID, req)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockAccessLevelRepo := mocks.NewMockAccessLevelRepository(ctrl)
	service := NewUserService(mockUserRepo, mockAccessLevelRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()

		mockUserRepo.EXPECT().
			Delete(ctx, userID).
			Return(nil)

		err := service.DeleteUser(ctx, userID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("Error", func(t *testing.T) {
		userID := uuid.New()

		mockUserRepo.EXPECT().
			Delete(ctx, userID).
			Return(errors.New("delete failed"))

		err := service.DeleteUser(ctx, userID)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestUserService_ListUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockAccessLevelRepo := mocks.NewMockAccessLevelRepository(ctrl)
	service := NewUserService(mockUserRepo, mockAccessLevelRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		users := []*models.User{
			{
				ID:        uuid.New(),
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
			},
			{
				ID:        uuid.New(),
				FirstName: "Jane",
				LastName:  "Smith",
				Email:     "jane@example.com",
			},
		}

		mockUserRepo.EXPECT().
			List(ctx, 10, 0).
			Return(users, 2, nil)

		mockAccessLevelRepo.EXPECT().
			GetUserAccessLevels(ctx, users[0].ID).
			Return([]*models.AccessLevel{}, nil)

		mockAccessLevelRepo.EXPECT().
			GetUserAccessLevels(ctx, users[1].ID).
			Return([]*models.AccessLevel{}, nil)

		resp, err := service.ListUsers(ctx, 1, 10)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(resp.Users) != 2 {
			t.Errorf("Expected 2 users, got %d", len(resp.Users))
		}
		if resp.Total != 2 {
			t.Errorf("Expected total 2, got %d", resp.Total)
		}
	})

	t.Run("DefaultPagination", func(t *testing.T) {
		mockUserRepo.EXPECT().
			List(ctx, 10, 0).
			Return([]*models.User{}, 0, nil)

		resp, err := service.ListUsers(ctx, 0, 0)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.Page != 1 {
			t.Errorf("Expected page 1, got %d", resp.Page)
		}
		if resp.PageSize != 10 {
			t.Errorf("Expected pageSize 10, got %d", resp.PageSize)
		}
	})

	t.Run("MaxPageSize", func(t *testing.T) {
		mockUserRepo.EXPECT().
			List(ctx, 10, 0).
			Return([]*models.User{}, 0, nil)

		resp, err := service.ListUsers(ctx, 1, 200)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.PageSize != 10 {
			t.Errorf("Expected pageSize capped at 10, got %d", resp.PageSize)
		}
	})
}

func TestUserService_AuthenticateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockAccessLevelRepo := mocks.NewMockAccessLevelRepository(ctrl)
	service := NewUserService(mockUserRepo, mockAccessLevelRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := &models.User{
			ID:        userID,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
		}

		auth := &models.UserAuthentication{
			UserID:       userID,
			PasswordHash: string(hashedPassword),
		}

		req := &dto.LoginRequest{
			Email:    "john@example.com",
			Password: password,
		}

		mockUserRepo.EXPECT().
			GetByEmail(ctx, req.Email).
			Return(user, nil)

		mockUserRepo.EXPECT().
			GetUserAuthentication(ctx, userID).
			Return(auth, nil)

		mockAccessLevelRepo.EXPECT().
			GetUserAccessLevels(ctx, userID).
			Return([]*models.AccessLevel{}, nil)

		resp, err := service.AuthenticateUser(ctx, req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.User.Email != user.Email {
			t.Errorf("Expected email %s, got %s", user.Email, resp.User.Email)
		}
		if resp.Message != "Login successful" {
			t.Errorf("Expected success message, got %s", resp.Message)
		}
	})

	t.Run("UserNotFound", func(t *testing.T) {
		req := &dto.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}

		mockUserRepo.EXPECT().
			GetByEmail(ctx, req.Email).
			Return(nil, errors.New("user not found"))

		_, err := service.AuthenticateUser(ctx, req)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		userID := uuid.New()
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

		user := &models.User{
			ID:    userID,
			Email: "john@example.com",
		}

		auth := &models.UserAuthentication{
			UserID:       userID,
			PasswordHash: string(hashedPassword),
		}

		req := &dto.LoginRequest{
			Email:    "john@example.com",
			Password: "wrongpassword",
		}

		mockUserRepo.EXPECT().
			GetByEmail(ctx, req.Email).
			Return(user, nil)

		mockUserRepo.EXPECT().
			GetUserAuthentication(ctx, userID).
			Return(auth, nil)

		_, err := service.AuthenticateUser(ctx, req)
		if err == nil {
			t.Fatal("Expected error for invalid password, got nil")
		}
	})

	t.Run("AuthenticationNotFound", func(t *testing.T) {
		userID := uuid.New()
		user := &models.User{
			ID:    userID,
			Email: "john@example.com",
		}

		req := &dto.LoginRequest{
			Email:    "john@example.com",
			Password: "password123",
		}

		mockUserRepo.EXPECT().
			GetByEmail(ctx, req.Email).
			Return(user, nil)

		mockUserRepo.EXPECT().
			GetUserAuthentication(ctx, userID).
			Return(nil, errors.New("auth not found"))

		_, err := service.AuthenticateUser(ctx, req)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestUserService_AssignAccessLevels(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockAccessLevelRepo := mocks.NewMockAccessLevelRepository(ctrl)
	service := NewUserService(mockUserRepo, mockAccessLevelRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()
		user := &models.User{
			ID: userID,
		}

		accessLevel := &models.AccessLevel{
			ID:   1,
			Name: "Admin",
		}

		req := &dto.AssignAccessLevelRequest{
			AccessLevelIDs: []int{1},
		}

		mockUserRepo.EXPECT().
			GetByID(ctx, userID).
			Return(user, nil)

		mockAccessLevelRepo.EXPECT().
			GetByID(ctx, 1).
			Return(accessLevel, nil)

		mockAccessLevelRepo.EXPECT().
			AssignToUser(ctx, userID, 1).
			Return(nil)

		err := service.AssignAccessLevels(ctx, userID, req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("UserNotFound", func(t *testing.T) {
		userID := uuid.New()
		req := &dto.AssignAccessLevelRequest{
			AccessLevelIDs: []int{1},
		}

		mockUserRepo.EXPECT().
			GetByID(ctx, userID).
			Return(nil, errors.New("user not found"))

		err := service.AssignAccessLevels(ctx, userID, req)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("AccessLevelNotFound", func(t *testing.T) {
		userID := uuid.New()
		user := &models.User{
			ID: userID,
		}

		req := &dto.AssignAccessLevelRequest{
			AccessLevelIDs: []int{999},
		}

		mockUserRepo.EXPECT().
			GetByID(ctx, userID).
			Return(user, nil)

		mockAccessLevelRepo.EXPECT().
			GetByID(ctx, 999).
			Return(nil, errors.New("access level not found"))

		err := service.AssignAccessLevels(ctx, userID, req)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestUserService_RemoveAccessLevel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockAccessLevelRepo := mocks.NewMockAccessLevelRepository(ctrl)
	service := NewUserService(mockUserRepo, mockAccessLevelRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()

		mockAccessLevelRepo.EXPECT().
			RemoveFromUser(ctx, userID, 1).
			Return(nil)

		err := service.RemoveAccessLevel(ctx, userID, 1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("Error", func(t *testing.T) {
		userID := uuid.New()

		mockAccessLevelRepo.EXPECT().
			RemoveFromUser(ctx, userID, 1).
			Return(errors.New("removal failed"))

		err := service.RemoveAccessLevel(ctx, userID, 1)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestUserService_GetUserAccessLevels(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockAccessLevelRepo := mocks.NewMockAccessLevelRepository(ctrl)
	service := NewUserService(mockUserRepo, mockAccessLevelRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()
		desc := "Admin access"
		accessLevels := []*models.AccessLevel{
			{
				ID:          1,
				Name:        "Admin",
				Description: &desc,
			},
		}

		mockAccessLevelRepo.EXPECT().
			GetUserAccessLevels(ctx, userID).
			Return(accessLevels, nil)

		resp, err := service.GetUserAccessLevels(ctx, userID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(resp) != 1 {
			t.Fatalf("Expected 1 access level, got %d", len(resp))
		}
		if resp[0].Name != "Admin" {
			t.Errorf("Expected name Admin, got %s", resp[0].Name)
		}
		if resp[0].Description != desc {
			t.Errorf("Expected description %s, got %s", desc, resp[0].Description)
		}
	})

	t.Run("NoAccessLevels", func(t *testing.T) {
		userID := uuid.New()

		mockAccessLevelRepo.EXPECT().
			GetUserAccessLevels(ctx, userID).
			Return([]*models.AccessLevel{}, nil)

		resp, err := service.GetUserAccessLevels(ctx, userID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(resp) != 0 {
			t.Errorf("Expected 0 access levels, got %d", len(resp))
		}
	})

	t.Run("Error", func(t *testing.T) {
		userID := uuid.New()

		mockAccessLevelRepo.EXPECT().
			GetUserAccessLevels(ctx, userID).
			Return(nil, errors.New("database error"))

		_, err := service.GetUserAccessLevels(ctx, userID)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}
