package repository

import (
	"database/sql"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

// OpenTestDB opens an in-memory SQLite database using the pure Go modernc.org/sqlite driver
func OpenTestDB() (*gorm.DB, error) {
	// Open using modernc.org/sqlite driver directly
	sqlDB, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}

	// Create GORM DB using the existing connection
	return gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{})
}
