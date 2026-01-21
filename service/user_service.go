package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/wabtcdi/user_service/dto"
	"github.com/wabtcdi/user_service/models"
	"github.com/wabtcdi/user_service/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo        repository.UserRepository
	accessLevelRepo repository.AccessLevelRepository
}

func NewUserService(userRepo repository.UserRepository, accessLevelRepo repository.AccessLevelRepository) *UserService {
	return &UserService{
		userRepo:        userRepo,
		accessLevelRepo: accessLevelRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Validate input
	if err := s.validateCreateUserRequest(req); err != nil {
		return nil, err
	}

	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user model
	user := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = &req.PhoneNumber
	}

	// Create authentication model
	auth := &models.UserAuthentication{
		PasswordHash: string(hashedPassword),
	}

	// Save to database
	err = s.userRepo.Create(ctx, user, auth)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return s.toUserResponse(ctx, user), nil
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toUserResponse(ctx, user), nil
}

func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Email != "" {
		// Check if new email is already taken by another user
		existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
		if existingUser != nil && existingUser.ID != id {
			return nil, fmt.Errorf("email %s is already taken", req.Email)
		}
		user.Email = req.Email
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = &req.PhoneNumber
	}

	// Save updates
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return s.toUserResponse(ctx, user), nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context, page, pageSize int) (*dto.ListUsersResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	users, total, err := s.userRepo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, err
	}

	userResponses := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, *s.toUserResponse(ctx, user))
	}

	return &dto.ListUsersResponse{
		Users:    userResponses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *UserService) AuthenticateUser(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Get authentication
	auth, err := s.userRepo.GetUserAuthentication(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(auth.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	return &dto.LoginResponse{
		User:    *s.toUserResponse(ctx, user),
		Message: "Login successful",
	}, nil
}

func (s *UserService) AssignAccessLevels(ctx context.Context, userID uuid.UUID, req *dto.AssignAccessLevelRequest) error {
	// Verify user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Assign each access level
	for _, accessLevelID := range req.AccessLevelIDs {
		// Verify access level exists
		_, err := s.accessLevelRepo.GetByID(ctx, accessLevelID)
		if err != nil {
			return fmt.Errorf("access level %d not found", accessLevelID)
		}

		err = s.accessLevelRepo.AssignToUser(ctx, userID, accessLevelID)
		if err != nil {
			return fmt.Errorf("failed to assign access level %d: %w", accessLevelID, err)
		}
	}

	return nil
}

func (s *UserService) RemoveAccessLevel(ctx context.Context, userID uuid.UUID, accessLevelID int) error {
	return s.accessLevelRepo.RemoveFromUser(ctx, userID, accessLevelID)
}

func (s *UserService) GetUserAccessLevels(ctx context.Context, userID uuid.UUID) ([]dto.AccessLevelResponse, error) {
	accessLevels, err := s.accessLevelRepo.GetUserAccessLevels(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.AccessLevelResponse, 0, len(accessLevels))
	for _, al := range accessLevels {
		desc := ""
		if al.Description != nil {
			desc = *al.Description
		}
		responses = append(responses, dto.AccessLevelResponse{
			ID:          al.ID,
			Name:        al.Name,
			Description: desc,
		})
	}

	return responses, nil
}

func (s *UserService) toUserResponse(ctx context.Context, user *models.User) *dto.UserResponse {
	response := &dto.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if user.PhoneNumber != nil {
		response.PhoneNumber = *user.PhoneNumber
	}

	// Get access levels
	accessLevels, err := s.accessLevelRepo.GetUserAccessLevels(ctx, user.ID)
	if err == nil && len(accessLevels) > 0 {
		response.AccessLevels = make([]dto.AccessLevelResponse, 0, len(accessLevels))
		for _, al := range accessLevels {
			desc := ""
			if al.Description != nil {
				desc = *al.Description
			}
			response.AccessLevels = append(response.AccessLevels, dto.AccessLevelResponse{
				ID:          al.ID,
				Name:        al.Name,
				Description: desc,
			})
		}
	}

	return response
}

func (s *UserService) validateCreateUserRequest(req *dto.CreateUserRequest) error {
	if strings.TrimSpace(req.FirstName) == "" {
		return fmt.Errorf("first name is required")
	}
	if strings.TrimSpace(req.LastName) == "" {
		return fmt.Errorf("last name is required")
	}
	if strings.TrimSpace(req.Email) == "" {
		return fmt.Errorf("email is required")
	}
	if !strings.Contains(req.Email, "@") {
		return fmt.Errorf("invalid email format")
	}
	if len(req.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	return nil
}
