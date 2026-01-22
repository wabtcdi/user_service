package mocks

import (
	"database/sql"

	"gorm.io/gorm"
)

// GormDB is a wrapper interface for *gorm.DB to enable mocking
type GormDB interface {
	DB() (*sql.DB, error)
	WithContext(ctx interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Find(dest interface{}, conds ...interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
	Update(column string, value interface{}) *gorm.DB
	Updates(values interface{}) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB
	Transaction(fc func(tx *gorm.DB) error) error
	Model(value interface{}) *gorm.DB
	Joins(query string, args ...interface{}) *gorm.DB
	Order(value interface{}) *gorm.DB
	Limit(limit int) *gorm.DB
	Offset(offset int) *gorm.DB
	Count(count *int64) *gorm.DB
	Preload(query string, args ...interface{}) *gorm.DB
}
