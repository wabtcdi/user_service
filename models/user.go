package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	FirstName   string         `json:"first_name" gorm:"column:first_name;size:50;not null"`
	LastName    string         `json:"last_name" gorm:"column:last_name;size:50;not null"`
	Email       string         `json:"email" gorm:"column:email;size:255;uniqueIndex;not null"`
	PhoneNumber *string        `json:"phone_number,omitempty" gorm:"column:phone_number;size:20"`
	CreatedAt   time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

func (User) TableName() string {
	return "users"
}

type UserAuthentication struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	UserID       uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	PasswordHash string         `json:"-" gorm:"column:password_hash;size:255;not null"`
	CreatedAt    time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
	User         *User          `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (UserAuthentication) TableName() string {
	return "user_authentications"
}

type AccessLevel struct {
	ID          int            `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string         `json:"name" gorm:"column:name;size:50;uniqueIndex;not null"`
	Description *string        `json:"description,omitempty" gorm:"column:description;type:text"`
	CreatedAt   time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

func (AccessLevel) TableName() string {
	return "access_levels"
}

type UserAccessLevel struct {
	UserID        uuid.UUID      `json:"user_id" gorm:"type:uuid;primaryKey;index"`
	AccessLevelID int            `json:"access_level_id" gorm:"primaryKey;index"`
	CreatedAt     time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
	User          *User          `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	AccessLevel   *AccessLevel   `json:"-" gorm:"foreignKey:AccessLevelID;constraint:OnDelete:CASCADE"`
}

func (UserAccessLevel) TableName() string {
	return "user_access_levels"
}
