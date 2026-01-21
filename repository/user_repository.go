package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wabtcdi/user_service/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User, auth *models.UserAuthentication) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.User, int, error)
	GetUserAuthentication(ctx context.Context, userID uuid.UUID) (*models.UserAuthentication, error)
}

type AccessLevelRepository interface {
	Create(ctx context.Context, accessLevel *models.AccessLevel) error
	GetByID(ctx context.Context, id int) (*models.AccessLevel, error)
	GetByName(ctx context.Context, name string) (*models.AccessLevel, error)
	List(ctx context.Context) ([]*models.AccessLevel, error)
	AssignToUser(ctx context.Context, userID uuid.UUID, accessLevelID int) error
	RemoveFromUser(ctx context.Context, userID uuid.UUID, accessLevelID int) error
	GetUserAccessLevels(ctx context.Context, userID uuid.UUID) ([]*models.AccessLevel, error)
}

type PostgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *models.User, auth *models.UserAuthentication) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Set IDs and timestamps
		user.ID = uuid.New()
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		// Create user
		if err := tx.Create(user).Error; err != nil {
			return fmt.Errorf("failed to insert user: %w", err)
		}

		// Set auth fields
		auth.ID = uuid.New()
		auth.UserID = user.ID
		auth.CreatedAt = time.Now()
		auth.UpdatedAt = time.Now()

		// Create authentication
		if err := tx.Create(auth).Error; err != nil {
			return fmt.Errorf("failed to insert authentication: %w", err)
		}

		return nil
	})
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	err := r.db.WithContext(ctx).Where("id = ?", id).First(user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// GetByIDWithAccessLevels retrieves a user by ID and preloads their access levels
func (r *PostgresUserRepository) GetByIDWithAccessLevels(ctx context.Context, id uuid.UUID) (*models.User, []*models.AccessLevel, error) {
	user := &models.User{}
	err := r.db.WithContext(ctx).Where("id = ?", id).First(user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Preload access levels using joins to avoid N+1 queries
	var accessLevels []*models.AccessLevel
	err = r.db.WithContext(ctx).
		Joins("INNER JOIN user_access_levels ON access_levels.id = user_access_levels.access_level_id").
		Where("user_access_levels.user_id = ? AND user_access_levels.deleted_at IS NULL", id).
		Find(&accessLevels).Error

	if err != nil {
		return nil, nil, fmt.Errorf("failed to preload access levels: %w", err)
	}

	return user, accessLevels, nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.WithContext(ctx).Where("email = ?", email).First(user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()
	result := r.db.WithContext(ctx).Model(user).Updates(map[string]interface{}{
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"email":        user.Email,
		"phone_number": user.PhoneNumber,
		"updated_at":   user.UpdatedAt,
	})
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *PostgresUserRepository) List(ctx context.Context, limit, offset int) ([]*models.User, int, error) {
	// Get total count
	var total int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Get users
	var users []*models.User
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, int(total), nil
}

func (r *PostgresUserRepository) GetUserAuthentication(ctx context.Context, userID uuid.UUID) (*models.UserAuthentication, error) {
	auth := &models.UserAuthentication{}
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(auth).Error
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("authentication not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get authentication: %w", err)
	}
	return auth, nil
}
