# GORM Quick Reference for User Service

## Database Connection

The application now uses GORM instead of database/sql. The connection is established in `cmd/app.go`:

```go
import "gorm.io/gorm"

// Database is *gorm.DB, not *sql.DB
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

// For Goose migrations, we get the underlying *sql.DB
sqlDB, err := db.DB()
```

## Working with Models

### Pointer Fields for Nullable Columns
```go
// Old (sql.NullString)
user.PhoneNumber = sql.NullString{String: "123-456", Valid: true}

// New (pointer)
phoneNumber := "123-456"
user.PhoneNumber = &phoneNumber

// Check if set
if user.PhoneNumber != nil {
    fmt.Println(*user.PhoneNumber)
}
```

### Soft Deletes
GORM automatically filters soft-deleted records:
```go
// This automatically excludes records where deleted_at IS NOT NULL
db.Where("email = ?", email).First(&user)

// To include soft-deleted records
db.Unscoped().Where("email = ?", email).First(&user)

// To permanently delete
db.Unscoped().Delete(&user)
```

## Common GORM Operations

### Create
```go
user := &models.User{
    FirstName: "John",
    LastName:  "Doe",
    Email:     "john@example.com",
}

// GORM automatically sets ID, CreatedAt, UpdatedAt
err := db.WithContext(ctx).Create(user).Error
```

### Create with Transaction
```go
err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(user).Error; err != nil {
        return err // auto-rollback
    }
    if err := tx.Create(auth).Error; err != nil {
        return err // auto-rollback
    }
    return nil // auto-commit
})
```

### Read (Find One)
```go
user := &models.User{}

// By ID
err := db.WithContext(ctx).Where("id = ?", id).First(user).Error

// By email
err := db.WithContext(ctx).Where("email = ?", email).First(user).Error

// Check for not found
if err == gorm.ErrRecordNotFound {
    return nil, fmt.Errorf("user not found")
}
```

### Read (Find Many)
```go
var users []*models.User

// Simple find all
err := db.WithContext(ctx).Find(&users).Error

// With conditions
err := db.WithContext(ctx).
    Where("email LIKE ?", "%@example.com").
    Order("created_at DESC").
    Limit(10).
    Find(&users).Error
```

### Update
```go
// Update specific fields using map
result := db.WithContext(ctx).Model(user).Updates(map[string]interface{}{
    "first_name": user.FirstName,
    "last_name":  user.LastName,
    "email":      user.Email,
    "updated_at": time.Now(),
})

// Check if record was found
if result.RowsAffected == 0 {
    return fmt.Errorf("user not found")
}
```

### Delete (Soft Delete)
```go
// Sets deleted_at to current timestamp
result := db.WithContext(ctx).Delete(&models.User{}, "id = ?", id)

if result.RowsAffected == 0 {
    return fmt.Errorf("user not found")
}
```

### Count
```go
var total int64
err := db.WithContext(ctx).Model(&models.User{}).Count(&total).Error
```

### Pagination
```go
limit := 10
offset := 20

var users []*models.User
err := db.WithContext(ctx).
    Order("created_at DESC").
    Limit(limit).
    Offset(offset).
    Find(&users).Error
```

### Joins
```go
var accessLevels []*models.AccessLevel
err := db.WithContext(ctx).
    Joins("INNER JOIN user_access_levels ON access_levels.id = user_access_levels.access_level_id").
    Where("user_access_levels.user_id = ?", userID).
    Order("access_levels.name ASC").
    Find(&accessLevels).Error
```

### Upsert (ON CONFLICT)
```go
import "gorm.io/gorm/clause"

userAccessLevel := &models.UserAccessLevel{
    UserID:        userID,
    AccessLevelID: accessLevelID,
    CreatedAt:     time.Now(),
    UpdatedAt:     time.Now(),
}

err := db.WithContext(ctx).Clauses(clause.OnConflict{
    Columns:   []clause.Column{{Name: "user_id"}, {Name: "access_level_id"}},
    DoUpdates: clause.Assignments(map[string]interface{}{
        "deleted_at": nil,
        "updated_at": time.Now(),
    }),
}).Create(userAccessLevel).Error
```

## Context Support

Always use `.WithContext(ctx)` for proper request cancellation and timeouts:

```go
// Good
db.WithContext(ctx).Where("id = ?", id).First(user)

// Bad (no cancellation support)
db.Where("id = ?", id).First(user)
```

## Error Handling

```go
err := db.WithContext(ctx).First(user).Error

// Check specific GORM errors
if err == gorm.ErrRecordNotFound {
    return nil, fmt.Errorf("user not found")
}

// Check for other errors
if err != nil {
    return nil, fmt.Errorf("database error: %w", err)
}
```

## Common Patterns in This Project

### Repository Pattern
All database operations are in `repository/` package:
- `user_repository.go` - User and Authentication operations
- `access_level_repository.go` - Access Level operations

### Service Layer
Business logic in `service/` package:
- `user_service.go` - User business logic
- `access_level_service.go` - Access level business logic

### Handler Layer
HTTP handlers in `handlers/` package (unchanged by GORM migration)

## Migration Notes

### Goose Migrations Still Used
- Migration files: `database/migrations/*.sql`
- Migrations run automatically on startup
- Uses underlying `*sql.DB` from GORM via `db.DB()`

### Tests
- Database tests are currently commented out (marked TODO)
- Consider using `gorm.io/driver/sqlite` with `:memory:` for tests
- Or use integration tests with test database

## Performance Tips

### Connection Pooling
```go
sqlDB, err := db.DB()

// SetMaxIdleConns sets the maximum number of connections in the idle pool
sqlDB.SetMaxIdleConns(10)

// SetMaxOpenConns sets the maximum number of open connections
sqlDB.SetMaxOpenConns(100)

// SetConnMaxLifetime sets the maximum time a connection may be reused
sqlDB.SetConnMaxLifetime(time.Hour)
```

### Select Specific Columns
```go
// Only fetch needed columns
var users []struct {
    ID    uuid.UUID
    Email string
}
db.Model(&models.User{}).Select("id", "email").Find(&users)
```

### Batch Operations
```go
users := []*models.User{user1, user2, user3}

// Insert in batches of 100
db.CreateInBatches(users, 100)
```

## Debugging

### Enable SQL Logging
```go
import "gorm.io/gorm/logger"

db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})
```

### See Generated SQL
```go
// Use .Debug() to see SQL in logs
db.Debug().Where("email = ?", email).First(&user)
```

## Resources

- GORM Documentation: https://gorm.io/docs/
- GORM PostgreSQL Driver: https://github.com/go-gorm/postgres
- Goose Migrations: https://github.com/pressly/goose

## Getting Started

1. Run migrations: `go run main.go -config=local`
2. Migrations run automatically on startup
3. Check `GORM_MIGRATION_SUMMARY.md` for migration details
