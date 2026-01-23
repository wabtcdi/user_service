# Microservice Template Configuration

This document provides templates and quick-start code for creating new Go microservices based on the user_service architecture.

## Quick Start: New Service Checklist

### 1. Project Initialization
```bash
# Create new service directory
SERVICE_NAME="your_service"
mkdir -p $SERVICE_NAME
cd $SERVICE_NAME

# Initialize Go module
go mod init github.com/yourorg/$SERVICE_NAME

# Create directory structure
mkdir -p cmd/health cmd/log models dto repository service handlers database/migrations resources mocks k8s .github
```

### 2. Core Files to Copy

Copy these files from user_service as templates:
- `main.go` - Update module imports
- `cmd/app.go` - Main application logic
- `cmd/properties.go` - Configuration structure
- `cmd/health/checker.go` - Health check handler
- `cmd/log/logger.go` - Logging configuration
- `resources/*.yaml` - Configuration files
- `.github/copilot-instructions.md` - This file

## Entity Template Generator

### Model Template (`models/entity.go`)
```go
package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

// {{ENTITY_NAME}} represents a {{ENTITY_DESCRIPTION}}
type {{ENTITY_NAME}} struct {
    ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
    Name        string         `json:"name" gorm:"column:name;size:50;not null"`
    Description *string        `json:"description,omitempty" gorm:"column:description;type:text"`
    // Add your fields here
    CreatedAt   time.Time      `json:"created_at" gorm:"column:created_at"`
    UpdatedAt   time.Time      `json:"updated_at" gorm:"column:updated_at"`
    DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

func ({{ENTITY_NAME}}) TableName() string {
    return "{{TABLE_NAME}}"
}
```

### DTO Template (`dto/entity_dto.go`)
```go
package dto

import (
    "time"
    "github.com/google/uuid"
)

// Create{{ENTITY_NAME}}Request represents the request to create a new {{ENTITY_NAME_LOWER}}
type Create{{ENTITY_NAME}}Request struct {
    Name        string `json:"name" validate:"required,min=1,max=50"`
    Description string `json:"description,omitempty" validate:"omitempty,max=255"`
    // Add your fields here
}

// Update{{ENTITY_NAME}}Request represents the request to update an existing {{ENTITY_NAME_LOWER}}
type Update{{ENTITY_NAME}}Request struct {
    Name        string `json:"name,omitempty" validate:"omitempty,min=1,max=50"`
    Description string `json:"description,omitempty" validate:"omitempty,max=255"`
    // Add your fields here
}

// {{ENTITY_NAME}}Response represents the {{ENTITY_NAME_LOWER}} data returned in API responses
type {{ENTITY_NAME}}Response struct {
    ID          uuid.UUID `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description,omitempty"`
    // Add your fields here
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// List{{ENTITY_NAME}}sResponse represents paginated list of {{ENTITY_NAME_LOWER}}s
type List{{ENTITY_NAME}}sResponse struct {
    {{ENTITY_NAME}}s []{{ENTITY_NAME}}Response `json:"{{ENTITY_NAME_LOWER}}s"`
    Total            int                        `json:"total"`
    Page             int                        `json:"page"`
    PageSize         int                        `json:"page_size"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message,omitempty"`
}
```

### Repository Template (`repository/entity_repository.go`)
```go
package repository

import (
    "context"
    "fmt"
    "time"
    "github.com/google/uuid"
    "github.com/yourorg/{{SERVICE_NAME}}/models"
    "gorm.io/gorm"
)

// {{ENTITY_NAME}}Repository defines the interface for {{ENTITY_NAME_LOWER}} data operations
type {{ENTITY_NAME}}Repository interface {
    Create(ctx context.Context, entity *models.{{ENTITY_NAME}}) error
    GetByID(ctx context.Context, id uuid.UUID) (*models.{{ENTITY_NAME}}, error)
    GetByName(ctx context.Context, name string) (*models.{{ENTITY_NAME}}, error)
    Update(ctx context.Context, entity *models.{{ENTITY_NAME}}) error
    Delete(ctx context.Context, id uuid.UUID) error
    List(ctx context.Context, limit, offset int) ([]*models.{{ENTITY_NAME}}, int, error)
}

// Postgres{{ENTITY_NAME}}Repository implements {{ENTITY_NAME}}Repository using PostgreSQL
type Postgres{{ENTITY_NAME}}Repository struct {
    db *gorm.DB
}

// NewPostgres{{ENTITY_NAME}}Repository creates a new Postgres{{ENTITY_NAME}}Repository
func NewPostgres{{ENTITY_NAME}}Repository(db *gorm.DB) *Postgres{{ENTITY_NAME}}Repository {
    return &Postgres{{ENTITY_NAME}}Repository{db: db}
}

// Create creates a new {{ENTITY_NAME_LOWER}} in the database
func (r *Postgres{{ENTITY_NAME}}Repository) Create(ctx context.Context, entity *models.{{ENTITY_NAME}}) error {
    entity.ID = uuid.New()
    entity.CreatedAt = time.Now()
    entity.UpdatedAt = time.Now()
    
    if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
        return fmt.Errorf("failed to create {{ENTITY_NAME_LOWER}}: %w", err)
    }
    
    return nil
}

// GetByID retrieves a {{ENTITY_NAME_LOWER}} by ID
func (r *Postgres{{ENTITY_NAME}}Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.{{ENTITY_NAME}}, error) {
    entity := &models.{{ENTITY_NAME}}{}
    err := r.db.WithContext(ctx).Where("id = ?", id).First(entity).Error
    if err == gorm.ErrRecordNotFound {
        return nil, fmt.Errorf("{{ENTITY_NAME_LOWER}} not found")
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get {{ENTITY_NAME_LOWER}}: %w", err)
    }
    return entity, nil
}

// GetByName retrieves a {{ENTITY_NAME_LOWER}} by name
func (r *Postgres{{ENTITY_NAME}}Repository) GetByName(ctx context.Context, name string) (*models.{{ENTITY_NAME}}, error) {
    entity := &models.{{ENTITY_NAME}}{}
    err := r.db.WithContext(ctx).Where("name = ?", name).First(entity).Error
    if err == gorm.ErrRecordNotFound {
        return nil, fmt.Errorf("{{ENTITY_NAME_LOWER}} not found")
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get {{ENTITY_NAME_LOWER}}: %w", err)
    }
    return entity, nil
}

// Update updates an existing {{ENTITY_NAME_LOWER}}
func (r *Postgres{{ENTITY_NAME}}Repository) Update(ctx context.Context, entity *models.{{ENTITY_NAME}}) error {
    entity.UpdatedAt = time.Now()
    
    result := r.db.WithContext(ctx).Model(entity).Updates(entity)
    if result.Error != nil {
        return fmt.Errorf("failed to update {{ENTITY_NAME_LOWER}}: %w", result.Error)
    }
    if result.RowsAffected == 0 {
        return fmt.Errorf("{{ENTITY_NAME_LOWER}} not found")
    }
    
    return nil
}

// Delete soft deletes a {{ENTITY_NAME_LOWER}}
func (r *Postgres{{ENTITY_NAME}}Repository) Delete(ctx context.Context, id uuid.UUID) error {
    result := r.db.WithContext(ctx).Delete(&models.{{ENTITY_NAME}}{}, "id = ?", id)
    if result.Error != nil {
        return fmt.Errorf("failed to delete {{ENTITY_NAME_LOWER}}: %w", result.Error)
    }
    if result.RowsAffected == 0 {
        return fmt.Errorf("{{ENTITY_NAME_LOWER}} not found")
    }
    
    return nil
}

// List retrieves a paginated list of {{ENTITY_NAME_LOWER}}s
func (r *Postgres{{ENTITY_NAME}}Repository) List(ctx context.Context, limit, offset int) ([]*models.{{ENTITY_NAME}}, int, error) {
    var entities []*models.{{ENTITY_NAME}}
    var total int64
    
    // Get total count
    if err := r.db.WithContext(ctx).Model(&models.{{ENTITY_NAME}}{}).Count(&total).Error; err != nil {
        return nil, 0, fmt.Errorf("failed to count {{ENTITY_NAME_LOWER}}s: %w", err)
    }
    
    // Get paginated results
    if err := r.db.WithContext(ctx).
        Limit(limit).
        Offset(offset).
        Order("created_at DESC").
        Find(&entities).Error; err != nil {
        return nil, 0, fmt.Errorf("failed to list {{ENTITY_NAME_LOWER}}s: %w", err)
    }
    
    return entities, int(total), nil
}
```

### Service Template (`service/entity_service.go`)
```go
package service

import (
    "context"
    "fmt"
    "github.com/google/uuid"
    "github.com/yourorg/{{SERVICE_NAME}}/dto"
    "github.com/yourorg/{{SERVICE_NAME}}/models"
    "github.com/yourorg/{{SERVICE_NAME}}/repository"
)

// {{ENTITY_NAME}}Service handles business logic for {{ENTITY_NAME_LOWER}}s
type {{ENTITY_NAME}}Service struct {
    repo repository.{{ENTITY_NAME}}Repository
}

// New{{ENTITY_NAME}}Service creates a new {{ENTITY_NAME}}Service
func New{{ENTITY_NAME}}Service(repo repository.{{ENTITY_NAME}}Repository) *{{ENTITY_NAME}}Service {
    return &{{ENTITY_NAME}}Service{repo: repo}
}

// Create{{ENTITY_NAME}} creates a new {{ENTITY_NAME_LOWER}}
func (s *{{ENTITY_NAME}}Service) Create{{ENTITY_NAME}}(ctx context.Context, req *dto.Create{{ENTITY_NAME}}Request) (*dto.{{ENTITY_NAME}}Response, error) {
    // Validate input
    if req.Name == "" {
        return nil, fmt.Errorf("name is required")
    }
    
    // Check if {{ENTITY_NAME_LOWER}} already exists
    existing, _ := s.repo.GetByName(ctx, req.Name)
    if existing != nil {
        return nil, fmt.Errorf("{{ENTITY_NAME_LOWER}} with name %s already exists", req.Name)
    }
    
    // Create model
    entity := &models.{{ENTITY_NAME}}{
        Name: req.Name,
    }
    if req.Description != "" {
        entity.Description = &req.Description
    }
    
    // Save to database
    if err := s.repo.Create(ctx, entity); err != nil {
        return nil, fmt.Errorf("failed to create {{ENTITY_NAME_LOWER}}: %w", err)
    }
    
    return s.to{{ENTITY_NAME}}Response(entity), nil
}

// Get{{ENTITY_NAME}} retrieves a {{ENTITY_NAME_LOWER}} by ID
func (s *{{ENTITY_NAME}}Service) Get{{ENTITY_NAME}}(ctx context.Context, id uuid.UUID) (*dto.{{ENTITY_NAME}}Response, error) {
    entity, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    return s.to{{ENTITY_NAME}}Response(entity), nil
}

// Update{{ENTITY_NAME}} updates an existing {{ENTITY_NAME_LOWER}}
func (s *{{ENTITY_NAME}}Service) Update{{ENTITY_NAME}}(ctx context.Context, id uuid.UUID, req *dto.Update{{ENTITY_NAME}}Request) (*dto.{{ENTITY_NAME}}Response, error) {
    // Get existing entity
    entity, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Update fields if provided
    if req.Name != "" {
        // Check if new name is already taken
        existing, _ := s.repo.GetByName(ctx, req.Name)
        if existing != nil && existing.ID != id {
            return nil, fmt.Errorf("name %s is already taken", req.Name)
        }
        entity.Name = req.Name
    }
    if req.Description != "" {
        entity.Description = &req.Description
    }
    
    // Save to database
    if err := s.repo.Update(ctx, entity); err != nil {
        return nil, fmt.Errorf("failed to update {{ENTITY_NAME_LOWER}}: %w", err)
    }
    
    return s.to{{ENTITY_NAME}}Response(entity), nil
}

// Delete{{ENTITY_NAME}} deletes a {{ENTITY_NAME_LOWER}}
func (s *{{ENTITY_NAME}}Service) Delete{{ENTITY_NAME}}(ctx context.Context, id uuid.UUID) error {
    return s.repo.Delete(ctx, id)
}

// List{{ENTITY_NAME}}s retrieves a paginated list of {{ENTITY_NAME_LOWER}}s
func (s *{{ENTITY_NAME}}Service) List{{ENTITY_NAME}}s(ctx context.Context, page, pageSize int) (*dto.List{{ENTITY_NAME}}sResponse, error) {
    if page < 1 {
        page = 1
    }
    if pageSize < 1 {
        pageSize = 10
    }
    
    offset := (page - 1) * pageSize
    entities, total, err := s.repo.List(ctx, pageSize, offset)
    if err != nil {
        return nil, err
    }
    
    responses := make([]dto.{{ENTITY_NAME}}Response, len(entities))
    for i, entity := range entities {
        responses[i] = *s.to{{ENTITY_NAME}}Response(entity)
    }
    
    return &dto.List{{ENTITY_NAME}}sResponse{
        {{ENTITY_NAME}}s: responses,
        Total:            total,
        Page:             page,
        PageSize:         pageSize,
    }, nil
}

// to{{ENTITY_NAME}}Response converts a model to a response DTO
func (s *{{ENTITY_NAME}}Service) to{{ENTITY_NAME}}Response(entity *models.{{ENTITY_NAME}}) *dto.{{ENTITY_NAME}}Response {
    response := &dto.{{ENTITY_NAME}}Response{
        ID:        entity.ID,
        Name:      entity.Name,
        CreatedAt: entity.CreatedAt,
        UpdatedAt: entity.UpdatedAt,
    }
    
    if entity.Description != nil {
        response.Description = *entity.Description
    }
    
    return response
}
```

### Service Interface Template (`service/interfaces.go`)
```go
package service

import (
    "context"
    "github.com/google/uuid"
    "github.com/yourorg/{{SERVICE_NAME}}/dto"
)

// {{ENTITY_NAME}}ServiceInterface defines the interface for {{ENTITY_NAME_LOWER}} service operations
type {{ENTITY_NAME}}ServiceInterface interface {
    Create{{ENTITY_NAME}}(ctx context.Context, req *dto.Create{{ENTITY_NAME}}Request) (*dto.{{ENTITY_NAME}}Response, error)
    Get{{ENTITY_NAME}}(ctx context.Context, id uuid.UUID) (*dto.{{ENTITY_NAME}}Response, error)
    Update{{ENTITY_NAME}}(ctx context.Context, id uuid.UUID, req *dto.Update{{ENTITY_NAME}}Request) (*dto.{{ENTITY_NAME}}Response, error)
    Delete{{ENTITY_NAME}}(ctx context.Context, id uuid.UUID) error
    List{{ENTITY_NAME}}s(ctx context.Context, page, pageSize int) (*dto.List{{ENTITY_NAME}}sResponse, error)
}

// Ensure {{ENTITY_NAME}}Service implements {{ENTITY_NAME}}ServiceInterface
var _ {{ENTITY_NAME}}ServiceInterface = (*{{ENTITY_NAME}}Service)(nil)
```

### Handler Template (`handlers/entity_handler.go`)
```go
package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "github.com/sirupsen/logrus"
    "github.com/yourorg/{{SERVICE_NAME}}/dto"
    "github.com/yourorg/{{SERVICE_NAME}}/service"
)

// {{ENTITY_NAME}}Handler handles HTTP requests for {{ENTITY_NAME_LOWER}}s
type {{ENTITY_NAME}}Handler struct {
    service service.{{ENTITY_NAME}}ServiceInterface
}

// New{{ENTITY_NAME}}Handler creates a new {{ENTITY_NAME}}Handler
func New{{ENTITY_NAME}}Handler(service service.{{ENTITY_NAME}}ServiceInterface) *{{ENTITY_NAME}}Handler {
    return &{{ENTITY_NAME}}Handler{service: service}
}

// Create{{ENTITY_NAME}} handles POST /{{ENTITY_NAME_LOWER}}s
func (h *{{ENTITY_NAME}}Handler) Create{{ENTITY_NAME}}(w http.ResponseWriter, r *http.Request) {
    var req dto.Create{{ENTITY_NAME}}Request
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        logrus.Errorf("Failed to decode request: %v", err)
        respondWithError(w, http.StatusBadRequest, "Invalid request body", err.Error())
        return
    }
    
    entity, err := h.service.Create{{ENTITY_NAME}}(r.Context(), &req)
    if err != nil {
        logrus.Errorf("Failed to create {{ENTITY_NAME_LOWER}}: %v", err)
        respondWithError(w, http.StatusBadRequest, "Failed to create {{ENTITY_NAME_LOWER}}", err.Error())
        return
    }
    
    respondWithJSON(w, http.StatusCreated, entity)
}

// Get{{ENTITY_NAME}} handles GET /{{ENTITY_NAME_LOWER}}s/{id}
func (h *{{ENTITY_NAME}}Handler) Get{{ENTITY_NAME}}(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := uuid.Parse(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid {{ENTITY_NAME_LOWER}} ID", err.Error())
        return
    }
    
    entity, err := h.service.Get{{ENTITY_NAME}}(r.Context(), id)
    if err != nil {
        logrus.Errorf("Failed to get {{ENTITY_NAME_LOWER}}: %v", err)
        respondWithError(w, http.StatusNotFound, "{{ENTITY_NAME}} not found", err.Error())
        return
    }
    
    respondWithJSON(w, http.StatusOK, entity)
}

// Update{{ENTITY_NAME}} handles PUT /{{ENTITY_NAME_LOWER}}s/{id}
func (h *{{ENTITY_NAME}}Handler) Update{{ENTITY_NAME}}(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := uuid.Parse(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid {{ENTITY_NAME_LOWER}} ID", err.Error())
        return
    }
    
    var req dto.Update{{ENTITY_NAME}}Request
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        logrus.Errorf("Failed to decode request: %v", err)
        respondWithError(w, http.StatusBadRequest, "Invalid request body", err.Error())
        return
    }
    
    entity, err := h.service.Update{{ENTITY_NAME}}(r.Context(), id, &req)
    if err != nil {
        logrus.Errorf("Failed to update {{ENTITY_NAME_LOWER}}: %v", err)
        respondWithError(w, http.StatusBadRequest, "Failed to update {{ENTITY_NAME_LOWER}}", err.Error())
        return
    }
    
    respondWithJSON(w, http.StatusOK, entity)
}

// Delete{{ENTITY_NAME}} handles DELETE /{{ENTITY_NAME_LOWER}}s/{id}
func (h *{{ENTITY_NAME}}Handler) Delete{{ENTITY_NAME}}(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := uuid.Parse(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid {{ENTITY_NAME_LOWER}} ID", err.Error())
        return
    }
    
    err = h.service.Delete{{ENTITY_NAME}}(r.Context(), id)
    if err != nil {
        logrus.Errorf("Failed to delete {{ENTITY_NAME_LOWER}}: %v", err)
        respondWithError(w, http.StatusNotFound, "{{ENTITY_NAME}} not found", err.Error())
        return
    }
    
    respondWithJSON(w, http.StatusOK, map[string]string{"message": "{{ENTITY_NAME}} deleted successfully"})
}

// List{{ENTITY_NAME}}s handles GET /{{ENTITY_NAME_LOWER}}s
func (h *{{ENTITY_NAME}}Handler) List{{ENTITY_NAME}}s(w http.ResponseWriter, r *http.Request) {
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
    
    if page < 1 {
        page = 1
    }
    if pageSize < 1 {
        pageSize = 10
    }
    
    response, err := h.service.List{{ENTITY_NAME}}s(r.Context(), page, pageSize)
    if err != nil {
        logrus.Errorf("Failed to list {{ENTITY_NAME_LOWER}}s: %v", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to list {{ENTITY_NAME_LOWER}}s", err.Error())
        return
    }
    
    respondWithJSON(w, http.StatusOK, response)
}

// Helper functions
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

### Migration Template (`database/migrations/YYYYMMDDHHMMSS_create_entities.sql`)
```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE {{TABLE_NAME}} (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Indexes for performance
CREATE INDEX idx_{{TABLE_NAME}}_name ON {{TABLE_NAME}}(name);
CREATE INDEX idx_{{TABLE_NAME}}_deleted_at ON {{TABLE_NAME}}(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS {{TABLE_NAME}};
-- +goose StatementEnd
```

## Variable Replacements

When using these templates, replace the following placeholders:

- `{{SERVICE_NAME}}` - Service name (e.g., "product_service", "order_service")
- `{{ENTITY_NAME}}` - Entity name in PascalCase (e.g., "Product", "Order")
- `{{ENTITY_NAME_LOWER}}` - Entity name in lowercase (e.g., "product", "order")
- `{{TABLE_NAME}}` - Database table name in snake_case (e.g., "products", "orders")
- `{{ENTITY_DESCRIPTION}}` - Brief description of the entity

## Example: Creating a Product Service

```bash
# 1. Replace placeholders
SERVICE_NAME="product_service"
ENTITY_NAME="Product"
ENTITY_NAME_LOWER="product"
TABLE_NAME="products"

# 2. Create files using templates above
# 3. Update imports in main.go and cmd/app.go
# 4. Run go mod tidy
# 5. Create and run migrations
# 6. Write tests
# 7. Test endpoints
```

## Configuration Template (`resources/local.yaml`)
```yaml
server:
  host: 0.0.0.0
  port: 8080
  readinessPath: /ready
  livenessPath: /health

resources:
  memory: 512Mi
  cpu: 500m
  storage: 1Gi
  threads: 10

logging:
  level: info
  format: json

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  name: {{SERVICE_NAME}}_db
```

## Testing Template

See `.github/copilot-instructions.md` for comprehensive testing patterns.

Quick test command:
```bash
# Run all tests with coverage
CGO_ENABLED=0 go test ./... -cover -coverprofile=coverage.out

# View coverage
go tool cover -html=coverage.out
```

## Next Steps

1. Copy templates to your new service
2. Replace all placeholders
3. Run `go mod tidy`
4. Create database migrations
5. Write tests for each layer
6. Document your APIs
7. Deploy and monitor
