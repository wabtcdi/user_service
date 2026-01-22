package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/wabtcdi/user_service/dto"
)

// UserServiceInterface defines the interface for user service operations
type UserServiceInterface interface {
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUser(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	ListUsers(ctx context.Context, page, pageSize int) (*dto.ListUsersResponse, error)
	AuthenticateUser(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	AssignAccessLevels(ctx context.Context, userID uuid.UUID, req *dto.AssignAccessLevelRequest) error
	GetUserAccessLevels(ctx context.Context, userID uuid.UUID) ([]dto.AccessLevelResponse, error)
}

// AccessLevelServiceInterface defines the interface for access level service operations
type AccessLevelServiceInterface interface {
	CreateAccessLevel(ctx context.Context, req *dto.CreateAccessLevelRequest) (*dto.AccessLevelResponse, error)
	GetAccessLevel(ctx context.Context, id int) (*dto.AccessLevelResponse, error)
	ListAccessLevels(ctx context.Context) ([]dto.AccessLevelResponse, error)
}

// Ensure UserService implements UserServiceInterface
var _ UserServiceInterface = (*UserService)(nil)

// Ensure AccessLevelService implements AccessLevelServiceInterface
var _ AccessLevelServiceInterface = (*AccessLevelService)(nil)
