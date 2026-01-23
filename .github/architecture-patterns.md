# Architecture Patterns and Best Practices

This document outlines the architectural patterns, design decisions, and best practices used in the user_service that should be replicated in other microservices.

## Table of Contents
1. [Clean Architecture Overview](#clean-architecture-overview)
2. [Layer Responsibilities](#layer-responsibilities)
3. [Dependency Injection Pattern](#dependency-injection-pattern)
4. [Testing Strategy](#testing-strategy)
5. [Error Handling Strategy](#error-handling-strategy)
6. [Database Patterns](#database-patterns)
7. [API Design Patterns](#api-design-patterns)
8. [Security Patterns](#security-patterns)

## Clean Architecture Overview

Our microservices follow Clean Architecture principles with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────┐
│                     Handlers Layer                       │
│  (HTTP Request/Response, Input Validation, Routing)     │
└──────────────────────┬──────────────────────────────────┘
                       │ Uses Interface
                       ▼
┌─────────────────────────────────────────────────────────┐
│                     Service Layer                        │
│  (Business Logic, Validation, Orchestration)            │
└──────────────────────┬──────────────────────────────────┘
                       │ Uses Interface
                       ▼
┌─────────────────────────────────────────────────────────┐
│                   Repository Layer                       │
│  (Data Access, Database Operations, Queries)            │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────┐
│                     Database (PostgreSQL)                │
└─────────────────────────────────────────────────────────┘

      DTOs           Models
  (API Contract)  (DB Entities)
```

### Key Principles

1. **Dependency Inversion**: Handlers depend on service interfaces, services depend on repository interfaces
2. **Single Responsibility**: Each layer has one clear responsibility
3. **Testability**: Interfaces enable easy mocking and testing
4. **Separation of Concerns**: Business logic isolated from infrastructure

## Layer Responsibilities

### 1. Models Layer (`models/`)

**Responsibility**: Define database entities and their structure

**What it does**:
- Defines GORM structs with proper tags
- Implements TableName() methods
- Defines relationships between entities

**What it does NOT do**:
- Business logic
- Validation
- API serialization logic

**Example**:
```go
type User struct {
    ID          uuid.UUID      `gorm:"type:uuid;primary_key"`
    FirstName   string         `gorm:"column:first_name;size:50;not null"`
    LastName    string         `gorm:"column:last_name;size:50;not null"`
    Email       string         `gorm:"column:email;size:255;uniqueIndex;not null"`
    CreatedAt   time.Time      `gorm:"column:created_at"`
    UpdatedAt   time.Time      `gorm:"column:updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}
```

### 2. DTO Layer (`dto/`)

**Responsibility**: Define API contract structures

**What it does**:
- Defines request structures
- Defines response structures
- Includes JSON tags for serialization
- Includes validation tags

**What it does NOT do**:
- Database operations
- Business logic
- Direct model conversion (that's in service layer)

**Example**:
```go
type CreateUserRequest struct {
    FirstName string `json:"first_name" validate:"required,min=1,max=50"`
    LastName  string `json:"last_name" validate:"required,min=1,max=50"`
    Email     string `json:"email" validate:"required,email"`
    Password  string `json:"password" validate:"required,min=8"`
}

type UserResponse struct {
    ID        uuid.UUID `json:"id"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}
```

### 3. Repository Layer (`repository/`)

**Responsibility**: Abstract database operations

**What it does**:
- Defines repository interfaces
- Implements CRUD operations
- Handles database queries
- Manages transactions
- Returns models

**What it does NOT do**:
- Business logic
- Validation
- DTO conversion
- HTTP handling

**Pattern**:
```go
// Interface definition
type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
    Update(ctx context.Context, user *models.User) error
    Delete(ctx context.Context, id uuid.UUID) error
    List(ctx context.Context, limit, offset int) ([]*models.User, int, error)
}

// Implementation
type PostgresUserRepository struct {
    db *gorm.DB
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
    var user models.User
    err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
    if err == gorm.ErrRecordNotFound {
        return nil, fmt.Errorf("user not found")
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    return &user, nil
}
```

### 4. Service Layer (`service/`)

**Responsibility**: Implement business logic

**What it does**:
- Validates business rules
- Orchestrates repository calls
- Converts models to DTOs
- Handles complex operations
- Coordinates transactions

**What it does NOT do**:
- HTTP handling
- Direct database operations
- JSON parsing

**Pattern**:
```go
// Interface (in interfaces.go)
type UserServiceInterface interface {
    CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
    GetUser(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error)
}

// Implementation
type UserService struct {
    userRepo repository.UserRepository
}

func (s *UserService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
    // 1. Validate business rules
    if err := s.validateCreateUserRequest(req); err != nil {
        return nil, err
    }
    
    // 2. Check business constraints
    existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
    if existingUser != nil {
        return nil, fmt.Errorf("user already exists")
    }
    
    // 3. Create model
    user := &models.User{
        FirstName: req.FirstName,
        LastName:  req.LastName,
        Email:     req.Email,
    }
    
    // 4. Call repository
    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    // 5. Convert to DTO
    return s.toUserResponse(user), nil
}
```

### 5. Handler Layer (`handlers/`)

**Responsibility**: Handle HTTP requests and responses

**What it does**:
- Parses HTTP requests
- Validates HTTP input
- Calls service methods
- Formats HTTP responses
- Sets status codes
- Logs requests

**What it does NOT do**:
- Business logic
- Database operations
- Complex data transformations

**Pattern**:
```go
type UserHandler struct {
    service service.UserServiceInterface // Use interface!
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    // 1. Parse request
    var req dto.CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        logrus.Errorf("Failed to decode request: %v", err)
        respondWithError(w, http.StatusBadRequest, "Invalid request", err.Error())
        return
    }
    
    // 2. Call service
    user, err := h.service.CreateUser(r.Context(), &req)
    if err != nil {
        logrus.Errorf("Failed to create user: %v", err)
        respondWithError(w, http.StatusBadRequest, "Failed to create user", err.Error())
        return
    }
    
    // 3. Return response
    respondWithJSON(w, http.StatusCreated, user)
}
```

## Dependency Injection Pattern

### Why Interfaces?

Interfaces enable:
1. **Testability**: Mock dependencies in tests
2. **Flexibility**: Swap implementations easily
3. **Loose Coupling**: Layers don't depend on concrete types

### Implementation Pattern

```go
// 1. Define interface in service layer
type UserServiceInterface interface {
    CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
}

// 2. Ensure implementation
var _ UserServiceInterface = (*UserService)(nil)

// 3. Handler uses interface
type UserHandler struct {
    service UserServiceInterface // Interface, not concrete type!
}

// 4. In app initialization
userRepo := repository.NewPostgresUserRepository(db)
userService := service.NewUserService(userRepo)      // Concrete type
userHandler := handlers.NewUserHandler(userService)   // Passed as interface

// 5. In tests, use mocks
mockService := new(MockUserService)
handler := handlers.NewUserHandler(mockService)       // Mock as interface
```

## Testing Strategy

### Layer Testing Approach

#### 1. Repository Tests
- Use **SQLite in-memory database**
- Test actual database operations
- No mocking needed

```go
func TestCreateUser(t *testing.T) {
    db := setupTestDB(t)
    defer cleanupTestDB(db)
    
    repo := repository.NewPostgresUserRepository(db)
    user := &models.User{Name: "Test"}
    
    err := repo.Create(context.Background(), user)
    assert.NoError(t, err)
    assert.NotEqual(t, uuid.Nil, user.ID)
}
```

#### 2. Service Tests
- Mock repository interfaces
- Test business logic
- Test validation
- Test error handling

```go
func TestUserService_CreateUser(t *testing.T) {
    mockRepo := new(MockUserRepository)
    service := service.NewUserService(mockRepo)
    
    mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, nil)
    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)
    
    req := &dto.CreateUserRequest{Email: "test@example.com"}
    response, err := service.CreateUser(context.Background(), req)
    
    assert.NoError(t, err)
    assert.NotNil(t, response)
    mockRepo.AssertExpectations(t)
}
```

#### 3. Handler Tests
- Mock service interfaces
- Use httptest for HTTP testing
- Test all status codes
- Test error responses

```go
func TestCreateUser(t *testing.T) {
    t.Run("Success", func(t *testing.T) {
        mockService := new(MockUserService)
        handler := NewUserHandler(mockService)
        
        expectedResponse := &dto.UserResponse{ID: uuid.New()}
        mockService.On("CreateUser", mock.Anything, mock.AnythingOfType("*dto.CreateUserRequest")).
            Return(expectedResponse, nil)
        
        body, _ := json.Marshal(dto.CreateUserRequest{Email: "test@example.com"})
        request := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
        recorder := httptest.NewRecorder()
        
        handler.CreateUser(recorder, request)
        
        assert.Equal(t, http.StatusCreated, recorder.Code)
        mockService.AssertExpectations(t)
    })
}
```

### Coverage Target

- **Repository**: 100%
- **Service**: 100%
- **Handlers**: 100%
- **Overall**: 100%

## Error Handling Strategy

### Error Wrapping Pattern

Always wrap errors with context:

```go
// Good
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}

// Bad
if err != nil {
    return err
}
```

### Error Types by Layer

#### Repository Layer
```go
// Check for specific GORM errors
if err == gorm.ErrRecordNotFound {
    return nil, fmt.Errorf("user not found")
}
if err != nil {
    return nil, fmt.Errorf("database error: %w", err)
}
```

#### Service Layer
```go
// Business logic errors
if existingUser != nil {
    return nil, fmt.Errorf("user with email %s already exists", email)
}

// Validation errors
if req.Name == "" {
    return nil, fmt.Errorf("name is required")
}

// Wrapped repository errors
if err := s.repo.Create(ctx, user); err != nil {
    return nil, fmt.Errorf("failed to create user: %w", err)
}
```

#### Handler Layer
```go
// HTTP-specific errors
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    logrus.Errorf("Failed to decode request: %v", err)
    respondWithError(w, http.StatusBadRequest, "Invalid request body", err.Error())
    return
}

// Service errors mapped to HTTP status
if err := h.service.CreateUser(r.Context(), &req); err != nil {
    logrus.Errorf("Service error: %v", err)
    respondWithError(w, http.StatusBadRequest, "Failed to create user", err.Error())
    return
}
```

### HTTP Status Code Guidelines

| Error Type | Status Code | Example |
|------------|-------------|---------|
| Invalid JSON | 400 Bad Request | Malformed request body |
| Validation failure | 400 Bad Request | Missing required fields |
| Already exists | 400 Bad Request | Duplicate email |
| Not found | 404 Not Found | User ID doesn't exist |
| Auth failure | 401 Unauthorized | Invalid credentials |
| Database error | 500 Internal Error | Connection failed |

## Database Patterns

### Transaction Pattern

Use transactions for multi-step operations:

```go
func (r *PostgresUserRepository) Create(ctx context.Context, user *models.User, auth *models.UserAuthentication) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // Step 1: Create user
        user.ID = uuid.New()
        if err := tx.Create(user).Error; err != nil {
            return fmt.Errorf("failed to create user: %w", err)
        }
        
        // Step 2: Create authentication
        auth.UserID = user.ID
        auth.ID = uuid.New()
        if err := tx.Create(auth).Error; err != nil {
            return fmt.Errorf("failed to create auth: %w", err)
        }
        
        return nil
    })
}
```

### Pagination Pattern

Standard pagination implementation:

```go
func (r *PostgresUserRepository) List(ctx context.Context, limit, offset int) ([]*models.User, int, error) {
    var users []*models.User
    var total int64
    
    // Count total (without limit/offset)
    if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
        return nil, 0, fmt.Errorf("failed to count: %w", err)
    }
    
    // Get paginated results
    if err := r.db.WithContext(ctx).
        Limit(limit).
        Offset(offset).
        Order("created_at DESC").
        Find(&users).Error; err != nil {
        return nil, 0, fmt.Errorf("failed to list: %w", err)
    }
    
    return users, int(total), nil
}
```

### Soft Delete Pattern

Use GORM's soft delete:

```go
type User struct {
    ID        uuid.UUID      `gorm:"type:uuid;primary_key"`
    Name      string         `gorm:"column:name;not null"`
    DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"` // Enables soft delete
}

// Delete automatically soft deletes
func (r *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
    result := r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id)
    if result.Error != nil {
        return fmt.Errorf("failed to delete: %w", result.Error)
    }
    if result.RowsAffected == 0 {
        return fmt.Errorf("user not found")
    }
    return nil
}
```

### Index Strategy

Always index:
- Primary keys (automatic)
- Foreign keys
- Unique constraints
- Frequently queried fields
- Soft delete columns

```sql
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE INDEX idx_user_access_levels_user_id ON user_access_levels(user_id);
```

## API Design Patterns

### RESTful Endpoint Structure

```
Resource Collection:
  POST   /entities           - Create new entity
  GET    /entities           - List entities (paginated)

Resource Item:
  GET    /entities/{id}      - Get single entity
  PUT    /entities/{id}      - Update entity
  DELETE /entities/{id}      - Delete entity

Resource Relationships:
  GET    /entities/{id}/related     - Get related items
  POST   /entities/{id}/related     - Add relationship
  DELETE /entities/{id}/related/{relatedId} - Remove relationship
```

### Request/Response Format

**Create Request**:
```json
{
  "name": "Example",
  "description": "Optional field"
}
```

**Success Response** (201 Created):
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "Example",
  "description": "Optional field",
  "created_at": "2026-01-22T10:00:00Z",
  "updated_at": "2026-01-22T10:00:00Z"
}
```

**Error Response** (400 Bad Request):
```json
{
  "error": "Failed to create entity",
  "message": "entity with name Example already exists"
}
```

**List Response** (200 OK):
```json
{
  "entities": [...],
  "total": 100,
  "page": 1,
  "page_size": 10
}
```

### Pagination Standard

- Query params: `?page=1&page_size=10`
- Default: page=1, page_size=10
- Include total count in response
- Order by creation date (newest first)

## Security Patterns

### Password Hashing

```go
import "golang.org/x/crypto/bcrypt"

// Hash password
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// Verify password
err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
if err != nil {
    return fmt.Errorf("invalid credentials")
}
```

### Input Validation

Service layer validates all inputs:

```go
func (s *UserService) validateCreateUserRequest(req *dto.CreateUserRequest) error {
    if req.FirstName == "" {
        return fmt.Errorf("first name is required")
    }
    if len(req.FirstName) > 50 {
        return fmt.Errorf("first name too long")
    }
    if !isValidEmail(req.Email) {
        return fmt.Errorf("invalid email format")
    }
    if len(req.Password) < 8 {
        return fmt.Errorf("password must be at least 8 characters")
    }
    return nil
}
```

### Logging Sensitive Data

**Never log**:
- Passwords
- API keys
- Tokens
- Personal identifiable information (PII)

```go
// Good
logrus.Infof("User created: ID=%s", user.ID)

// Bad
logrus.Infof("User created: %+v", user) // May contain sensitive data
```

## Configuration Management

### Environment-Specific Configs

- `local.yaml` - Development
- `cloud.yaml` - Production
- `test.yaml` - Testing

### Environment Variable Support

Configuration files support environment variable expansion:

```yaml
database:
  host: ${DB_HOST}           # Expands $DB_HOST
  password: ${DB_PASSWORD}   # Expands $DB_PASSWORD
```

### Connection Pooling

Configure based on resources:

```go
maxOpenConns := cfg.Resources.Threads * 2
maxIdleConns := cfg.Resources.Threads / 2
connMaxLifetime := 5 * time.Minute

db.SetMaxOpenConns(maxOpenConns)
db.SetMaxIdleConns(maxIdleConns)
db.SetConnMaxLifetime(connMaxLifetime)
```

## Health Checks

### Liveness Probe

Simple check that service is running:

```go
func livenessHandler(w http.ResponseWriter, _ *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "OK")
}
```

### Readiness Probe

Checks database connectivity:

```go
type Checker struct {
    DB *gorm.DB
}

func (c *Checker) Check(w http.ResponseWriter, _ *http.Request) {
    sqlDB, _ := c.DB.DB()
    if err := sqlDB.Ping(); err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        return
    }
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "OK")
}
```

## Summary

Follow these patterns for consistent, maintainable, testable microservices:

1. ✅ Use clean architecture with clear layer separation
2. ✅ Define interfaces for all services and repositories
3. ✅ Use dependency injection throughout
4. ✅ Achieve 100% test coverage
5. ✅ Wrap all errors with context
6. ✅ Use transactions for multi-step operations
7. ✅ Implement soft deletes with DeletedAt
8. ✅ Standard pagination (page, page_size)
9. ✅ RESTful API design
10. ✅ Proper HTTP status codes
11. ✅ Secure password handling
12. ✅ Environment-based configuration
13. ✅ Health check endpoints
14. ✅ Structured logging with logrus
15. ✅ Database migrations with Goose
