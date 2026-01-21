package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	FirstName   string `json:"first_name" validate:"required,min=1,max=50"`
	LastName    string `json:"last_name" validate:"required,min=1,max=50"`
	Email       string `json:"email" validate:"required,email,max=255"`
	PhoneNumber string `json:"phone_number,omitempty" validate:"omitempty,max=20"`
	Password    string `json:"password" validate:"required,min=8"`
}

// UpdateUserRequest represents the request to update an existing user
type UpdateUserRequest struct {
	FirstName   string `json:"first_name,omitempty" validate:"omitempty,min=1,max=50"`
	LastName    string `json:"last_name,omitempty" validate:"omitempty,min=1,max=50"`
	Email       string `json:"email,omitempty" validate:"omitempty,email,max=255"`
	PhoneNumber string `json:"phone_number,omitempty" validate:"omitempty,max=20"`
}

// UserResponse represents the user data returned in API responses
type UserResponse struct {
	ID           uuid.UUID             `json:"id"`
	FirstName    string                `json:"first_name"`
	LastName     string                `json:"last_name"`
	Email        string                `json:"email"`
	PhoneNumber  string                `json:"phone_number,omitempty"`
	AccessLevels []AccessLevelResponse `json:"access_levels,omitempty"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
}

// LoginRequest represents authentication credentials
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the response after successful authentication
type LoginResponse struct {
	User    UserResponse `json:"user"`
	Message string       `json:"message"`
}

// AssignAccessLevelRequest represents the request to assign access levels to a user
type AssignAccessLevelRequest struct {
	AccessLevelIDs []int `json:"access_level_ids" validate:"required,min=1"`
}

// AccessLevelResponse represents an access level in API responses
type AccessLevelResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// ListUsersResponse represents paginated list of users
type ListUsersResponse struct {
	Users    []UserResponse `json:"users"`
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

// CreateAccessLevelRequest represents the request to create an access level
type CreateAccessLevelRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=50"`
	Description string `json:"description,omitempty"`
}
