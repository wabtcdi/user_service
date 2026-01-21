# Testing Guide for GORM

This document explains how to test GORM repositories in the user service.

## Testing Approaches

We provide two approaches for testing GORM code:

### Approach 1: SQLite In-Memory Database (Recommended for CI/CD)

**Location:** `repository/*_test.go`

**Pros:**
- Fast execution (no external dependencies)
- Isolated tests (each test gets a fresh database)
- No cleanup required
- Works great in CI/CD pipelines

**Cons:**
- Requires CGO enabled (C compiler needed)
- SQLite may have minor behavioral differences from PostgreSQL

**Setup:**
```bash
# On Linux/Mac (requires GCC)
CGO_ENABLED=1 go test ./repository/...

# On Windows (requires MinGW or similar)
# Install TDM-GCC or MinGW-w64 first
set CGO_ENABLED=1
go test ./repository/...
```

**Files:**
- `repository/user_repository_test.go` - User repository tests
- `repository/access_level_repository_test.go` - Access level tests

### Approach 2: Integration Tests with Test Database (Alternative)

**Pros:**
- Tests against actual PostgreSQL database
- 100% production-like behavior
- No CGO required

**Cons:**
- Slower than in-memory tests
- Requires test database setup
- Needs cleanup between tests

**Setup:**
```bash
# Start test database (Docker)
docker run --name test-postgres \
  -e POSTGRES_PASSWORD=testpass \
  -e POSTGRES_DB=testdb \
  -p 5433:5432 \
  -d postgres:15-alpine

# Run tests
TEST_DB_DSN="host=localhost port=5433 user=postgres password=testpass dbname=testdb sslmode=disable" \
  go test ./repository/...
```

## Current Test Status

### ✅ Tests Created
- 19 comprehensive tests covering all repository methods
- Tests for user CRUD operations
- Tests for access level management
- Tests for user-access level relationships
- Tests for edge cases (not found, duplicates, etc.)

### ⚠️ Current Limitation
The SQLite tests require CGO which needs a C compiler (GCC/MinGW) that's not available in the current Windows environment. 

### Solutions

**Option 1: Install MinGW/GCC on Windows**
```bash
# Install TDM-GCC from https://jmeubank.github.io/tdm-gcc/
# Or use Chocolatey:
choco install mingw

# Then run tests:
$env:CGO_ENABLED=1
go test ./repository/... -v
```

**Option 2: Use Integration Tests**
Create a test database and run integration tests (see Approach 2 above).

**Option 3: Run Tests in Docker/WSL**
```bash
# In WSL or Docker container
CGO_ENABLED=1 go test ./repository/... -v
```

**Option 4: Use CI/CD**
Run tests in GitHub Actions, GitLab CI, or similar where GCC is available.

## Test Coverage

Our tests cover:

### User Repository
- ✅ Create user with authentication
- ✅ Get user by ID
- ✅ Get user by email
- ✅ Update user information
- ✅ Soft delete user
- ✅ List users with pagination
- ✅ Get user authentication
- ✅ Get user with preloaded access levels (demonstrates GORM relationships)
- ✅ Not found scenarios
- ✅ Error handling

### Access Level Repository
- ✅ Create access level
- ✅ Get by ID
- ✅ Get by name
- ✅ List all access levels
- ✅ Assign access level to user
- ✅ Handle duplicate assignments (upsert)
- ✅ Remove access level from user
- ✅ Get user's access levels
- ✅ Handle multiple access levels per user
- ✅ Not found scenarios

## Running Specific Tests

```bash
# Run all repository tests
go test ./repository/... -v

# Run specific test
go test ./repository/... -v -run TestUserRepository_Create

# Run with coverage
go test ./repository/... -cover

# Run with detailed coverage
go test ./repository/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Test Database Schema

The test helper function `setupTestDB()` automatically creates the schema using GORM AutoMigrate:

```go
db.AutoMigrate(
    &models.User{},
    &models.UserAuthentication{},
    &models.AccessLevel{},
    &models.UserAccessLevel{},
)
```

This ensures tests have a consistent schema matching the production Goose migrations.

## Example Test

```go
func TestUserRepository_Create(t *testing.T) {
    db := setupTestDB(t) // Creates in-memory database
    repo := NewPostgresUserRepository(db)
    ctx := context.Background()

    user := &models.User{
        FirstName: "John",
        LastName:  "Doe",
        Email:     "john.doe@example.com",
    }
    auth := &models.UserAuthentication{
        PasswordHash: "hashedpassword123",
    }

    err := repo.Create(ctx, user, auth)
    if err != nil {
        t.Fatalf("Failed to create user: %v", err)
    }

    // Assertions...
}
```

## Best Practices

1. **Use context.Background()** for tests
2. **Create fresh database per test** via `setupTestDB()`
3. **Test both success and error paths**
4. **Use descriptive test names** (`TestUserRepository_Create`, not `TestCreate`)
5. **Verify all returned data**, not just errors
6. **Test edge cases** (not found, duplicates, empty results)
7. **Add sleep between operations** if testing timestamps

## Continuous Integration Example

```yaml
# .github/workflows/test.yml
name: Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.25'
      
      - name: Install dependencies
        run: go mod download
      
      - name: Run tests
        run: CGO_ENABLED=1 go test ./... -v -coverprofile=coverage.out
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
```

## Next Steps

1. **Install GCC/MinGW** to run SQLite tests locally
2. **Set up test database** for integration testing
3. **Add CI/CD pipeline** for automated testing
4. **Increase coverage** by adding more edge case tests
5. **Add benchmarks** to measure GORM performance

## Resources

- [GORM Testing Guide](https://gorm.io/docs/testing.html)
- [SQLite Driver](https://github.com/glebarez/sqlite) - Pure Go alternative (no CGO)
- [Testify](https://github.com/stretchr/testify) - Advanced testing assertions
- [Dockertest](https://github.com/ory/dockertest) - Integration testing with Docker
