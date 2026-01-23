# GitHub Copilot Instructions for Go Microservices

## Project Context
This is a Go-based RESTful microservice following clean architecture patterns with:
- **Framework**: Gorilla Mux for routing
- **Database**: PostgreSQL with GORM ORM
- **Migrations**: Goose for database migrations
- **Testing**: testify for assertions and mocking
- **Logging**: logrus for structured logging
- **Configuration**: YAML-based with environment variable support

## Architecture Layers

### 1. Models (`models/`)
- GORM entities with proper tags
- Use `gorm.DeletedAt` for soft deletes
- UUID primary keys for main entities
- Timestamps: `CreatedAt`, `UpdatedAt`, `DeletedAt`
- Proper foreign key relationships with cascade

**Example Structure:**
```go
type Entity struct {
    ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
    Name      string         `json:"name" gorm:"column:name;size:50;not null"`
    CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
    UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
    DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

func (Entity) TableName() string {
    return "entities"
}
```

### 2. DTOs (`dto/`)
- Request and response structures separate from models
- JSON tags for API serialization
- Validation tags for input validation
- Keep business logic out of DTOs

**Naming Convention:**
- `Create{Entity}Request`
- `Update{Entity}Request`
- `{Entity}Response`
- `List{Entity}Response` (for paginated lists)
- `ErrorResponse`

**Example:**
```go
type CreateEntityRequest struct {
    Name        string `json:"name" validate:"required,min=1,max=50"`
    Description string `json:"description,omitempty" validate:"omitempty,max=255"`
}

type EntityResponse struct {
    ID          uuid.UUID `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message,omitempty"`
}
```

### 3. Repository Layer (`repository/`)
- Interface-based for testability
- GORM-based implementations
- Context support for all methods
- Transaction handling in Create operations
- Proper error wrapping with context

**Interface Pattern:**
```go
type EntityRepository interface {
    Create(ctx context.Context, entity *models.Entity) error
    GetByID(ctx context.Context, id uuid.UUID) (*models.Entity, error)
    GetByName(ctx context.Context, name string) (*models.Entity, error)
    Update(ctx context.Context, entity *models.Entity) error
    Delete(ctx context.Context, id uuid.UUID) error
    List(ctx context.Context, limit, offset int) ([]*models.Entity, int, error)
}

type PostgresEntityRepository struct {
    db *gorm.DB
}

func NewPostgresEntityRepository(db *gorm.DB) *PostgresEntityRepository {
    return &PostgresEntityRepository{db: db}
}
```

### 4. Service Layer (`service/`)
- Business logic and validation
- Interface definitions for mocking (`interfaces.go`)
- Orchestrates repository calls
- Error handling with meaningful messages
- Transaction coordination when needed

**Service Interface Pattern:**
```go
// In service/interfaces.go
type EntityServiceInterface interface {
    CreateEntity(ctx context.Context, req *dto.CreateEntityRequest) (*dto.EntityResponse, error)
    GetEntity(ctx context.Context, id uuid.UUID) (*dto.EntityResponse, error)
    UpdateEntity(ctx context.Context, id uuid.UUID, req *dto.UpdateEntityRequest) (*dto.EntityResponse, error)
    DeleteEntity(ctx context.Context, id uuid.UUID) error
    ListEntities(ctx context.Context, page, pageSize int) (*dto.ListEntitiesResponse, error)
}

// Ensure implementation
var _ EntityServiceInterface = (*EntityService)(nil)

// In service/entity_service.go
type EntityService struct {
    repo EntityRepository
}

func NewEntityService(repo EntityRepository) *EntityService {
    return &EntityService{repo: repo}
}
```

### 5. Handlers (`handlers/`)
- HTTP request/response handling
- Input validation
- Status code management
- Error response formatting
- Use interfaces for services (for testability)

**Handler Pattern:**
```go
type EntityHandler struct {
    service service.EntityServiceInterface
}

func NewEntityHandler(service service.EntityServiceInterface) *EntityHandler {
    return &EntityHandler{service: service}
}

func (h *EntityHandler) CreateEntity(w http.ResponseWriter, r *http.Request) {
    var req dto.CreateEntityRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        logrus.Errorf("Failed to decode request: %v", err)
        respondWithError(w, http.StatusBadRequest, "Invalid request body", err.Error())
        return
    }

    entity, err := h.service.CreateEntity(r.Context(), &req)
    if err != nil {
        logrus.Errorf("Failed to create entity: %v", err)
        respondWithError(w, http.StatusBadRequest, "Failed to create entity", err.Error())
        return
    }

    respondWithJSON(w, http.StatusCreated, entity)
}
```

## Code Standards

### Error Handling
- Always wrap errors with context: `fmt.Errorf("operation failed: %w", err)`
- Log errors at handler level with logrus
- Return meaningful error messages to clients
- Use appropriate HTTP status codes
- Check for `gorm.ErrRecordNotFound` specifically

**Example:**
```go
func (r *PostgresEntityRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Entity, error) {
    entity := &models.Entity{}
    err := r.db.WithContext(ctx).Where("id = ?", id).First(entity).Error
    if err == gorm.ErrRecordNotFound {
        return nil, fmt.Errorf("entity not found")
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get entity: %w", err)
    }
    return entity, nil
}
```

### Logging
- Use logrus for structured logging
- Log at handler level for request/response
- Use appropriate log levels: Debug, Info, Warn, Error
- Include context in log messages
- Never log sensitive data (passwords, tokens)

### Testing
- 100% code coverage target
- Unit tests for all layers
- Use testify/mock for mocking
- Use testify/assert for assertions
- Test success, error, and edge cases
- Name tests descriptively with subtests

**Test Pattern:**
```go
func TestCreateEntity(t *testing.T) {
    t.Run("Success", func(t *testing.T) {
        // Arrange
        mockService := new(MockEntityService)
        handler := NewEntityHandler(mockService)
        expectedResponse := &dto.EntityResponse{
            ID:   uuid.New(),
            Name: "Test Entity",
        }
        mockService.On("CreateEntity", mock.Anything, mock.AnythingOfType("*dto.CreateEntityRequest")).Return(expectedResponse, nil)
        
        body, _ := json.Marshal(dto.CreateEntityRequest{Name: "Test Entity"})
        request := httptest.NewRequest(http.MethodPost, "/entities", bytes.NewReader(body))
        recorder := httptest.NewRecorder()
        
        // Act
        handler.CreateEntity(recorder, request)
        
        // Assert
        assert.Equal(t, http.StatusCreated, recorder.Code)
        var response dto.EntityResponse
        json.Unmarshal(recorder.Body.Bytes(), &response)
        assert.Equal(t, expectedResponse.ID, response.ID)
        mockService.AssertExpectations(t)
    })
    
    t.Run("Invalid Request Body", func(t *testing.T) {
        mockService := new(MockEntityService)
        handler := NewEntityHandler(mockService)
        
        request := httptest.NewRequest(http.MethodPost, "/entities", bytes.NewReader([]byte("invalid json")))
        recorder := httptest.NewRecorder()
        
        handler.CreateEntity(recorder, request)
        
        assert.Equal(t, http.StatusBadRequest, recorder.Code)
    })
    
    t.Run("Service Error", func(t *testing.T) {
        mockService := new(MockEntityService)
        handler := NewEntityHandler(mockService)
        mockService.On("CreateEntity", mock.Anything, mock.AnythingOfType("*dto.CreateEntityRequest")).Return(nil, errors.New("creation failed"))
        
        body, _ := json.Marshal(dto.CreateEntityRequest{Name: "Test"})
        request := httptest.NewRequest(http.MethodPost, "/entities", bytes.NewReader(body))
        recorder := httptest.NewRecorder()
        
        handler.CreateEntity(recorder, request)
        
        assert.Equal(t, http.StatusBadRequest, recorder.Code)
        mockService.AssertExpectations(t)
    })
}
```

### Database Migrations
- Use Goose for migrations
- Migrations in `database/migrations/`
- Timestamp-based naming: `YYYYMMDDHHMMSS_description.sql`
- Include both Up and Down migrations
- Create indexes for foreign keys and frequently queried fields

**Example Migration:**
```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE entities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_entities_name ON entities(name);
CREATE INDEX idx_entities_deleted_at ON entities(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS entities;
-- +goose StatementEnd
```

### Configuration
- YAML-based configuration files
- Support for environment variable expansion with `$VAR` or `${VAR}`
- Separate configs for local, cloud, test
- Configuration struct in `cmd/properties.go`

**Config Structure:**
```go
type Config struct {
    Database struct {
        Host     string `yaml:"host"`
        Port     int    `yaml:"port"`
        User     string `yaml:"user"`
        Password string `yaml:"password"`
        Name     string `yaml:"name"`
    } `yaml:"database"`
    Server struct {
        Host          string `yaml:"host"`
        Port          int    `yaml:"port"`
        ReadinessPath string `yaml:"readinessPath"`
        LivenessPath  string `yaml:"livenessPath"`
    } `yaml:"server"`
    Resources struct {
        Memory  string `yaml:"memory"`
        CPU     string `yaml:"cpu"`
        Storage string `yaml:"storage"`
        Threads int    `yaml:"threads"`
    } `yaml:"resources"`
    Logging struct {
        Level  string `yaml:"level"`
        Format string `yaml:"format"`
    } `yaml:"logging"`
}
```

### API Design
- RESTful endpoints
- Proper HTTP methods (GET, POST, PUT, DELETE)
- UUID path parameters for resources
- Query parameters for pagination (page, page_size)
- Consistent response format

**Standard Endpoints:**
```
POST   /entities          - Create new entity
GET    /entities          - List entities (with pagination)
GET    /entities/{id}     - Get entity by ID
PUT    /entities/{id}     - Update entity
DELETE /entities/{id}     - Delete entity (soft delete)
```

### HTTP Status Codes
- **200 OK** - Successful GET, PUT
- **201 Created** - Successful POST
- **400 Bad Request** - Invalid input, validation errors
- **401 Unauthorized** - Authentication failed
- **404 Not Found** - Resource not found
- **500 Internal Server Error** - Unexpected server error

## Project Structure
```
/
├── cmd/                    # Application entry point and setup
│   ├── app.go             # Main application initialization
│   ├── app_test.go        # Application tests
│   ├── properties.go      # Configuration structs
│   ├── properties_test.go # Configuration tests
│   ├── health/            # Health check handlers
│   │   ├── checker.go
│   │   └── checker_test.go
│   └── log/               # Logging configuration
│       ├── logger.go
│       └── logger_test.go
├── models/                # Domain models (GORM entities)
│   └── entity.go
├── dto/                   # Data Transfer Objects
│   └── entity_dto.go
├── repository/            # Data access layer
│   ├── entity_repository.go
│   ├── entity_repository_test.go
│   └── sqlite_test_helper.go
├── service/               # Business logic layer
│   ├── interfaces.go      # Service interfaces for testing
│   ├── entity_service.go
│   └── entity_service_test.go
├── handlers/              # HTTP handlers
│   ├── entity_handler.go
│   └── entity_handler_test.go
├── database/
│   └── migrations/        # Goose migrations
│       └── YYYYMMDDHHMMSS_initial_setup.sql
├── resources/             # Configuration files
│   ├── local.yaml
│   ├── cloud.yaml
│   └── test.yaml
├── mocks/                 # Generated mocks
│   └── README.md
├── k8s/                   # Kubernetes deployment files
│   ├── deployment.yaml
│   ├── service.yaml
│   └── README.md
├── .github/               # GitHub configuration
│   └── copilot-instructions.md
├── main.go                # Application entry point
├── main_test.go           # Main tests
├── go.mod                 # Go module definition
└── go.sum                 # Go dependencies lock file
```

## Dependencies

### Core Dependencies
```go
require (
    github.com/google/uuid v1.6.0
    github.com/gorilla/mux v1.8.1
    github.com/pressly/goose/v3 v3.26.0
    github.com/sirupsen/logrus v1.9.3
    golang.org/x/crypto v0.47.0
    gopkg.in/yaml.v3 v3.0.1
    gorm.io/driver/postgres v1.6.0
    gorm.io/gorm v1.31.1
)
```

### Testing Dependencies
```go
require (
    github.com/DATA-DOG/go-sqlmock v1.5.2
    github.com/stretchr/testify v1.11.1
    gorm.io/driver/sqlite v1.6.0
    modernc.org/sqlite v1.44.3
)
```

## Naming Conventions

### Files
- **Snake case**: `entity_service.go`, `entity_repository.go`
- **Test files**: `*_test.go`
- **Mock files**: `mock_*.go`

### Variables
- **Camel case**: `entityService`, `entityRepo`, `userID`
- **Constants**: `MaxRetries`, `DefaultPageSize`
- **Package-level vars**: `ErrNotFound`, `DefaultTimeout`

### Functions
- **Exported**: PascalCase `CreateEntity`, `GetByID`, `NewEntityService`
- **Internal**: camelCase `validateInput`, `toResponse`, `handleError`

### Interfaces
- **Interface suffix or descriptive**: `EntityRepository`, `EntityServiceInterface`
- **Mock prefix**: `MockEntityService`, `MockEntityRepository`

## Best Practices

### Always Do ✅
- Use `context.Context` for all database operations
- Implement interfaces for services and repositories
- Write unit tests for all new code (target 100% coverage)
- Log errors with context using logrus
- Validate input in service layer
- Use transactions for multi-step operations
- Return DTOs from services, not models
- Handle errors gracefully with proper wrapping
- Use proper HTTP status codes
- Document complex logic with comments
- Use UUID for entity IDs (except lookup tables with integer IDs)
- Implement soft deletes with `DeletedAt` timestamp
- Create database indexes for performance
- Use pagination for list endpoints
- Preload relationships to avoid N+1 queries

### Never Do ❌
- Return database models directly from handlers
- Perform business logic in handlers
- Ignore errors
- Use `panic` in production code (except initialization)
- Store passwords in plain text
- Skip input validation
- Use `SELECT *` in raw queries
- Forget to close resources
- Hard-code configuration values
- Skip error wrapping
- Log sensitive data (passwords, tokens, API keys)
- Use `fmt.Println` for logging (use logrus)
- Create handlers without testing
- Expose internal errors to API clients

## Security Best Practices
- Hash passwords with bcrypt (cost 10-12)
- Validate all user input
- Use parameterized queries (GORM handles this)
- Implement proper authentication/authorization
- Log security-relevant events
- Never log sensitive data (passwords, tokens)
- Use HTTPS in production
- Implement rate limiting for public endpoints
- Validate UUID format before database queries
- Sanitize error messages returned to clients

## Performance Best Practices
- Use connection pooling (configured in app.go)
- Add database indexes for frequently queried fields
- Use pagination for list endpoints (default: page_size=10)
- Preload relationships when needed (avoid N+1 queries)
- Set appropriate timeout contexts
- Configure GORM to log slow queries (>200ms)
- Use batch operations for bulk inserts
- Cache frequently accessed data when appropriate

## Application Initialization Pattern

**In `cmd/app.go`:**
```go
func Init(configName string, opener DBOpener, starter ServerStarter) error {
    // 1. Load configuration
    cfg, err := loadConfiguration(configName)
    if err != nil {
        return err
    }
    
    // 2. Configure logging
    log.Configure(cfg.Logging.Level, cfg.Logging.Format)
    
    // 3. Connect to database
    db, err := connectDatabase(cfg, opener)
    if err != nil {
        return err
    }
    
    // 4. Run migrations
    if err := runMigrations(db); err != nil {
        return err
    }
    
    // 5. Initialize repositories
    entityRepo := repository.NewPostgresEntityRepository(db)
    
    // 6. Initialize services
    entityService := service.NewEntityService(entityRepo)
    
    // 7. Initialize handlers
    entityHandler := handlers.NewEntityHandler(entityService)
    
    // 8. Setup routes
    router := createRouter(cfg, entityHandler)
    
    // 9. Start server
    return startServer(cfg, router, starter)
}
```

## Router Setup Pattern

```go
func createRouter(cfg Config, entityHandler *handlers.EntityHandler) *mux.Router {
    r := mux.NewRouter()
    
    // Health checks
    checker := &health.Checker{DB: db}
    r.HandleFunc(cfg.Server.LivenessPath, livenessHandler).Methods("GET")
    r.HandleFunc(cfg.Server.ReadinessPath, checker.Check).Methods("GET")
    
    // Entity routes
    r.HandleFunc("/entities", entityHandler.CreateEntity).Methods("POST")
    r.HandleFunc("/entities", entityHandler.ListEntities).Methods("GET")
    r.HandleFunc("/entities/{id}", entityHandler.GetEntity).Methods("GET")
    r.HandleFunc("/entities/{id}", entityHandler.UpdateEntity).Methods("PUT")
    r.HandleFunc("/entities/{id}", entityHandler.DeleteEntity).Methods("DELETE")
    
    return r
}
```

## When Creating New Microservices

### Step-by-Step Checklist

1. **Setup Project Structure**
   ```bash
   mkdir -p cmd/health cmd/log models dto repository service handlers database/migrations resources mocks k8s
   ```

2. **Initialize Go Module**
   ```bash
   go mod init github.com/yourorg/service_name
   ```

3. **Create Configuration Files**
   - `resources/local.yaml`
   - `resources/cloud.yaml`
   - `resources/test.yaml`

4. **Define Models** (`models/`)
   - Create GORM entities
   - Add TableName() methods
   - Define relationships

5. **Create Database Migrations** (`database/migrations/`)
   - Use Goose format
   - Include indexes
   - Add foreign keys

6. **Create DTOs** (`dto/`)
   - Request structures
   - Response structures
   - Error response structure

7. **Build Repository Layer** (`repository/`)
   - Define interface
   - Implement PostgreSQL repository
   - Write repository tests

8. **Create Service Layer** (`service/`)
   - Define interface in `interfaces.go`
   - Implement service
   - Add validation logic
   - Write service tests

9. **Implement Handlers** (`handlers/`)
   - Create handler struct with service interface
   - Implement HTTP methods
   - Write handler tests

10. **Update Main Application** (`cmd/app.go`)
    - Initialize repositories
    - Initialize services
    - Initialize handlers
    - Setup routes

11. **Write Tests**
    - Repository tests (with SQLite)
    - Service tests (with mocks)
    - Handler tests (with mocks)
    - Integration tests

12. **Create Documentation**
    - API documentation
    - README
    - Quick start guide

## Commands Reference

```bash
# Development
go run main.go -config=local                 # Run locally
go run main.go -config=cloud                 # Run with cloud config

# Testing
go test ./... -v                             # Run all tests
go test ./... -cover                         # Run with coverage
go test ./... -coverprofile=coverage.out     # Generate coverage file
go tool cover -html=coverage.out             # View coverage in browser
go test ./handlers/... -v                    # Test specific package
CGO_ENABLED=0 go test ./... -v               # Test without CGO

# Build
go build -o build/service .                  # Build binary
CGO_ENABLED=0 go build -o build/service .    # Build without CGO

# Dependencies
go mod tidy                                  # Clean up dependencies
go mod download                              # Download dependencies
go mod vendor                                # Vendor dependencies

# Database Migrations
goose -dir database/migrations postgres "postgresql://user:pass@localhost:5432/dbname" up
goose -dir database/migrations postgres "postgresql://user:pass@localhost:5432/dbname" down
goose -dir database/migrations postgres "postgresql://user:pass@localhost:5432/dbname" status

# Code Quality
go fmt ./...                                 # Format code
go vet ./...                                 # Vet code
golangci-lint run                            # Run linter (if installed)

# Generate Mocks
mockgen -source=repository/entity_repository.go -destination=mocks/mock_entity_repository.go
```

## Common Patterns

### Transaction Pattern
```go
func (r *PostgresEntityRepository) Create(ctx context.Context, entity *models.Entity) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        entity.ID = uuid.New()
        entity.CreatedAt = time.Now()
        entity.UpdatedAt = time.Now()
        
        if err := tx.Create(entity).Error; err != nil {
            return fmt.Errorf("failed to create entity: %w", err)
        }
        
        // Additional operations in same transaction
        
        return nil
    })
}
```

### Pagination Pattern
```go
func (r *PostgresEntityRepository) List(ctx context.Context, limit, offset int) ([]*models.Entity, int, error) {
    var entities []*models.Entity
    var total int64
    
    // Get total count
    if err := r.db.WithContext(ctx).Model(&models.Entity{}).Count(&total).Error; err != nil {
        return nil, 0, fmt.Errorf("failed to count entities: %w", err)
    }
    
    // Get paginated results
    if err := r.db.WithContext(ctx).
        Limit(limit).
        Offset(offset).
        Order("created_at DESC").
        Find(&entities).Error; err != nil {
        return nil, 0, fmt.Errorf("failed to list entities: %w", err)
    }
    
    return entities, int(total), nil
}
```

### Response Helper Pattern
```go
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, err := json.Marshal(payload)
    if err != nil {
        logrus.Errorf("Failed to marshal response: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, error string, message string) {
    respondWithJSON(w, code, dto.ErrorResponse{
        Error:   error,
        Message: message,
    })
}
```

## Testing Patterns

### Repository Test with SQLite
```go
func TestCreateEntity(t *testing.T) {
    db := setupTestDB(t)
    defer cleanupTestDB(db)
    
    repo := repository.NewPostgresEntityRepository(db)
    
    entity := &models.Entity{
        Name: "Test Entity",
    }
    
    err := repo.Create(context.Background(), entity)
    assert.NoError(t, err)
    assert.NotEqual(t, uuid.Nil, entity.ID)
}
```

### Service Test with Mocks
```go
type MockEntityRepository struct {
    mock.Mock
}

func (m *MockEntityRepository) Create(ctx context.Context, entity *models.Entity) error {
    args := m.Called(ctx, entity)
    return args.Error(0)
}

func TestEntityService_CreateEntity(t *testing.T) {
    mockRepo := new(MockEntityRepository)
    service := service.NewEntityService(mockRepo)
    
    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Entity")).Return(nil)
    
    req := &dto.CreateEntityRequest{Name: "Test"}
    response, err := service.CreateEntity(context.Background(), req)
    
    assert.NoError(t, err)
    assert.NotNil(t, response)
    mockRepo.AssertExpectations(t)
}
```

## CI/CD Considerations
- Fast test execution (<5s for unit tests)
- No external dependencies in unit tests
- Proper exit codes (0 for success)
- Environment variable support for configuration
- Docker-friendly build process
- Health check endpoints for orchestrators
- Graceful shutdown handling
- Kubernetes deployment manifests included

## Additional Resources
- GORM Documentation: https://gorm.io/docs/
- Gorilla Mux Documentation: https://github.com/gorilla/mux
- Testify Documentation: https://github.com/stretchr/testify
- Goose Migrations: https://github.com/pressly/goose
