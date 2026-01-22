# Mocks Directory

This directory contains all mock implementations generated using MockGen for unit testing the user_service application.

## Generated Mocks

### 1. mock_server_starter_test.go
**Source:** `cmd/app.go` (ServerStarter interface)  
**Location:** `cmd/mock_server_starter_test.go`  
**Package:** `cmd`  
**Purpose:** Mock HTTP server lifecycle for testing server startup and shutdown behavior

**Generated with:**
```bash
mockgen -source=cmd/app.go -destination=cmd/mock_server_starter_test.go -package=cmd ServerStarter
```

**Used in:**
- `TestStartServer()`
- `TestStartServer_StarterError()`
- `TestInit_ConfigError()`
- `TestInit_DatabaseError()`
- `TestStartServer_AddressFormat()`

### 2. mock_user_repository.go
**Source:** `repository/user_repository.go`  
**Package:** `mocks`  
**Size:** ~10 KB  
**Purpose:** Mock repository interfaces for user and access level operations

**Interfaces Mocked:**
- `MockUserRepository` - User CRUD operations
- `MockAccessLevelRepository` - Access level management

**Generated with:**
```bash
mockgen -source=repository/user_repository.go -destination=mocks/mock_user_repository.go -package=mocks
```

**Methods:**
- UserRepository: Create, GetByID, GetByEmail, Update, Delete, List, GetUserAuthentication
- AccessLevelRepository: Create, GetByID, GetByName, List, AssignToUser, RemoveFromUser, GetUserAccessLevels

**Future Use:**
- Service layer unit tests
- Handler unit tests with mocked services

### 3. mock_health_checker.go
**Package:** `mocks`  
**Purpose:** Mock health check handler for testing HTTP endpoints

**Interface:**
```go
type MockHealthChecker interface {
    Check(w http.ResponseWriter, r *http.Request)
}
```

**Future Use:**
- Testing router registration
- Testing middleware integration with health checks

### 4. mock_pinger.go
**Package:** `mocks`  
**Purpose:** Mock database ping operations

**Interface:**
```go
type MockPinger interface {
    Ping() error
}
```

**Future Use:**
- Testing connection health checks
- Testing retry logic

### 5. mock_sql_db.go
**Package:** `mocks`  
**Purpose:** Mock sql.DB operations for connection pool testing

**Methods:**
- `Ping() error`
- `SetMaxOpenConns(n int)`
- `SetMaxIdleConns(n int)`
- `SetConnMaxLifetime(d interface{})`
- `Close() error`

**Future Use:**
- Testing connection pool configuration
- Testing database lifecycle

### 6. gorm_db.go
**Package:** `mocks`  
**Purpose:** GORM DB interface wrapper to enable mocking

**Interface Methods:**
- DB, WithContext, Create, First, Find, Where, Update, Updates, Delete
- Transaction, Model, Joins, Order, Limit, Offset, Count, Preload

**Note:** This is an interface definition, not a mock. Used as foundation for mocking GORM operations.

## Usage Examples

### Example 1: Mocking ServerStarter
```go
func TestMyFunction(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockStarter := cmd.NewMockServerStarter(ctrl)
    mockStarter.EXPECT().Start("127.0.0.1:8080", gomock.Any()).Return(nil)

    // Use mockStarter in your test
}
```

### Example 2: Mocking UserRepository
```go
func TestMyService(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockUserRepo := mocks.NewMockUserRepository(ctrl)
    mockUserRepo.EXPECT().
        GetByID(gomock.Any(), userID).
        Return(expectedUser, nil)

    // Use mockUserRepo in your service test
}
```

### Example 3: Custom GORM Mock (in checker_test.go)
```go
// Create sqlmock
mockDB, mock, _ := sqlmock.New(sqlmock.MonitorPingsOption(true))
mock.ExpectPing()

// Wrap in GORM
dialector := &mockDialector{sqlDB: mockDB}
gormDB, _ := gorm.Open(dialector, &gorm.Config{
    DisableAutomaticPing: true,
})

// Use gormDB in tests
```

## Regenerating Mocks

When interfaces change, regenerate mocks using:

```bash
# Install mockgen if not already installed
go install github.com/golang/mock/mockgen@v1.6.0

# Regenerate all mocks
cd /path/to/user_service

# Repository mocks
mockgen -source=repository/user_repository.go -destination=mocks/mock_user_repository.go -package=mocks

# Server starter mock
mockgen -source=cmd/app.go -destination=cmd/mock_server_starter_test.go -package=cmd ServerStarter
```

## Testing Best Practices

1. **Always use gomock.NewController** in test functions
2. **Always call ctrl.Finish()** with defer
3. **Set explicit expectations** with EXPECT() before using mocks
4. **Use gomock.Any()** for flexible argument matching
5. **Verify expectations** - gomock automatically fails tests if expectations not met

## Mock Verification

Mocks automatically verify that:
- All expected method calls were made
- Method calls were made with correct arguments
- Method calls were made in expected order (if specified)

## Notes

- All mocks use gomock framework
- Mocks are generated, do not edit manually
- Re-run mockgen after interface changes
- Some complex types (like GORM) require custom dialector approach

## Directory Structure
```
mocks/
├── gorm_db.go                    # GORM interface wrapper
├── mock_health_checker.go        # Health check mock
├── mock_pinger.go                # Database ping mock
├── mock_sql_db.go                # SQL DB operations mock
└── mock_user_repository.go       # Repository interfaces mock
```

## Related Documentation

- See `UNIT_TESTS_SUMMARY.md` for complete testing documentation
- See `cmd/app_test.go` for usage examples
- See `cmd/health/checker_test.go` for custom GORM mocking pattern
