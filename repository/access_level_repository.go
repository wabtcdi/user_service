package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wabtcdi/user_service/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresAccessLevelRepository struct {
	db *gorm.DB
}

func NewPostgresAccessLevelRepository(db *gorm.DB) *PostgresAccessLevelRepository {
	return &PostgresAccessLevelRepository{db: db}
}

func (r *PostgresAccessLevelRepository) Create(ctx context.Context, accessLevel *models.AccessLevel) error {
	accessLevel.CreatedAt = time.Now()
	accessLevel.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Create(accessLevel).Error; err != nil {
		return fmt.Errorf("failed to create access level: %w", err)
	}
	return nil
}

func (r *PostgresAccessLevelRepository) GetByID(ctx context.Context, id int) (*models.AccessLevel, error) {
	accessLevel := &models.AccessLevel{}
	err := r.db.WithContext(ctx).Where("id = ?", id).First(accessLevel).Error
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("access level not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get access level: %w", err)
	}
	return accessLevel, nil
}

func (r *PostgresAccessLevelRepository) GetByName(ctx context.Context, name string) (*models.AccessLevel, error) {
	accessLevel := &models.AccessLevel{}
	err := r.db.WithContext(ctx).Where("name = ?", name).First(accessLevel).Error
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("access level not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get access level: %w", err)
	}
	return accessLevel, nil
}

func (r *PostgresAccessLevelRepository) List(ctx context.Context) ([]*models.AccessLevel, error) {
	var accessLevels []*models.AccessLevel
	err := r.db.WithContext(ctx).Order("name ASC").Find(&accessLevels).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list access levels: %w", err)
	}

	return accessLevels, nil
}

func (r *PostgresAccessLevelRepository) AssignToUser(ctx context.Context, userID uuid.UUID, accessLevelID int) error {
	now := time.Now()
	userAccessLevel := &models.UserAccessLevel{
		UserID:        userID,
		AccessLevelID: accessLevelID,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Use GORM's Clauses with OnConflict to handle upsert
	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "access_level_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"deleted_at": nil, "updated_at": now}),
	}).Create(userAccessLevel).Error

	if err != nil {
		return fmt.Errorf("failed to assign access level to user: %w", err)
	}
	return nil
}

func (r *PostgresAccessLevelRepository) RemoveFromUser(ctx context.Context, userID uuid.UUID, accessLevelID int) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND access_level_id = ?", userID, accessLevelID).
		Delete(&models.UserAccessLevel{})

	if result.Error != nil {
		return fmt.Errorf("failed to remove access level from user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user access level not found")
	}
	return nil
}

func (r *PostgresAccessLevelRepository) GetUserAccessLevels(ctx context.Context, userID uuid.UUID) ([]*models.AccessLevel, error) {
	var accessLevels []*models.AccessLevel
	err := r.db.WithContext(ctx).
		Joins("INNER JOIN user_access_levels ON access_levels.id = user_access_levels.access_level_id").
		Where("user_access_levels.user_id = ? AND user_access_levels.deleted_at IS NULL", userID).
		Order("access_levels.name ASC").
		Find(&accessLevels).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get user access levels: %w", err)
	}

	return accessLevels, nil
}
