# GitHub Copilot Configuration for Go Microservices

This directory contains comprehensive GitHub Copilot configuration and templates for building consistent, high-quality Go microservices based on the user_service architecture.

## 📁 Files in This Directory

### 1. `copilot-instructions.md`
**Primary Copilot configuration file** containing:
- Project context and tech stack
- Architecture layer definitions
- Code standards and conventions
- Testing patterns and strategies
- Best practices and anti-patterns
- Security guidelines
- Performance optimization tips
- Complete examples for each layer

**Use this for**: General coding guidance, standards, and patterns

### 2. `service-templates.md`
**Ready-to-use code templates** including:
- Complete entity templates (Model, DTO, Repository, Service, Handler)
- Database migration templates
- Configuration file templates
- Step-by-step new service creation guide
- Variable replacement guide

**Use this for**: Quickly bootstrapping new entities or services

### 3. `architecture-patterns.md`
**In-depth architectural guidance** covering:
- Clean architecture principles
- Layer responsibilities and boundaries
- Dependency injection patterns
- Testing strategies by layer
- Error handling patterns
- Database patterns (transactions, pagination, soft deletes)
- API design patterns
- Security patterns
- Real-world examples

**Use this for**: Understanding design decisions and architectural patterns

## 🚀 Quick Start

### For New Microservice

1. **Copy the user_service structure**:
   ```bash
   cp -r user_service new_service
   cd new_service
   ```

2. **Update module name**:
   ```bash
   # In go.mod, change:
   module github.com/wabtcdi/user_service
   # To:
   module github.com/yourorg/new_service
   ```

3. **Follow checklist** in `service-templates.md`:
   - Create models
   - Create migrations
   - Build repository layer
   - Build service layer
   - Build handler layer
   - Write tests (100% coverage target)

### For New Entity in Existing Service

1. **Use templates** from `service-templates.md`
2. **Replace placeholders**:
   - `{{ENTITY_NAME}}` → Your entity name (e.g., "Product")
   - `{{ENTITY_NAME_LOWER}}` → Lowercase name (e.g., "product")
   - `{{TABLE_NAME}}` → Database table name (e.g., "products")
   - `{{SERVICE_NAME}}` → Service name (e.g., "product_service")

3. **Follow the pattern**:
   - Create model in `models/`
   - Create DTOs in `dto/`
   - Create repository with interface in `repository/`
   - Create service with interface in `service/`
   - Create handler in `handlers/`
   - Write comprehensive tests for each layer

## 📖 How to Use These Files

### When Writing New Code

**Copilot will automatically reference these files** when you:
- Create new files in standard locations
- Write comments asking for specific patterns
- Use naming conventions from the templates

**Example prompts that work well**:
```go
// Create a Product entity following the standard pattern

// Implement ProductRepository interface with PostgreSQL

// Write handler for CreateProduct following the pattern

// Add comprehensive tests for ProductService
```

### Reference the Documentation

When in doubt:
1. Check `copilot-instructions.md` for standards
2. Check `service-templates.md` for templates
3. Check `architecture-patterns.md` for deep explanations

## 🎯 Key Principles

### 1. Clean Architecture
```
Handlers → Services → Repositories → Database
   ↓          ↓           ↓
  DTOs     Business    Models
           Logic
```

### 2. Interface-Based Design
```go
// Define interface
type EntityServiceInterface interface {
    CreateEntity(ctx context.Context, req *dto.CreateEntityRequest) (*dto.EntityResponse, error)
}

// Use interface in handlers
type EntityHandler struct {
    service EntityServiceInterface // Not concrete type!
}
```

### 3. Test Everything
```go
// Repository: Use SQLite in-memory
// Service: Mock repository
// Handler: Mock service
// Target: 100% coverage
```

### 4. Error Handling
```go
// Always wrap errors with context
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}
```

## 🛠️ Technology Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.25+ |
| Router | Gorilla Mux |
| ORM | GORM |
| Database | PostgreSQL |
| Migrations | Goose |
| Testing | testify |
| Logging | logrus |
| Config | YAML |
| UUID | google/uuid |
| Crypto | golang.org/x/crypto/bcrypt |

## 📋 Standard Project Structure

```
service_name/
├── .github/
│   ├── copilot-instructions.md
│   ├── service-templates.md
│   └── architecture-patterns.md
├── cmd/
│   ├── app.go
│   ├── properties.go
│   ├── health/
│   └── log/
├── models/
├── dto/
├── repository/
├── service/
│   └── interfaces.go
├── handlers/
├── database/
│   └── migrations/
├── resources/
│   ├── local.yaml
│   ├── cloud.yaml
│   └── test.yaml
├── mocks/
├── k8s/
├── main.go
└── go.mod
```

## ✅ Checklist for New Services

- [ ] Project structure created
- [ ] Go module initialized
- [ ] Configuration files created
- [ ] Database migrations created
- [ ] Models defined with GORM tags
- [ ] DTOs created for requests/responses
- [ ] Repository interface and implementation
- [ ] Service interface and implementation
- [ ] Handlers created
- [ ] Routes registered in app.go
- [ ] Unit tests written (100% coverage)
- [ ] Integration tests written
- [ ] API documentation created
- [ ] README updated
- [ ] Health checks implemented
- [ ] Logging configured
- [ ] Error handling implemented

## 🧪 Testing Standards

### Coverage Targets
- Repository: 100%
- Service: 100%
- Handlers: 100%
- Overall: 100%

### Test Commands
```bash
# Run all tests
go test ./... -v

# With coverage
go test ./... -cover -coverprofile=coverage.out

# View coverage
go tool cover -html=coverage.out

# Without CGO (Windows)
CGO_ENABLED=0 go test ./... -v
```

### Test Patterns
- Use subtests with `t.Run()`
- Test success, error, and edge cases
- Mock all external dependencies
- Use testify/assert and testify/mock
- Name tests descriptively

## 🔒 Security Standards

- ✅ Hash passwords with bcrypt
- ✅ Validate all user input
- ✅ Use parameterized queries (GORM)
- ✅ Never log sensitive data
- ✅ Implement soft deletes
- ✅ Use UUIDs for entity IDs
- ✅ Return sanitized error messages to clients

## 🚀 Performance Best Practices

- ✅ Use connection pooling
- ✅ Create database indexes
- ✅ Implement pagination
- ✅ Preload relationships (avoid N+1)
- ✅ Set query timeouts
- ✅ Log slow queries (>200ms)
- ✅ Use transactions for multi-step ops

## 📚 Additional Resources

### Documentation Files (in parent directory)
- `API_DOCUMENTATION.md` - Complete API reference
- `TESTING_GUIDE.md` - Comprehensive testing guide
- `QUICK_START.md` - Quick start for developers
- `README_COMPLETE.md` - Project overview
- `HANDLERS_TESTS_SUMMARY.md` - Handler testing examples

### External Resources
- [GORM Documentation](https://gorm.io/docs/)
- [Gorilla Mux](https://github.com/gorilla/mux)
- [Testify](https://github.com/stretchr/testify)
- [Goose Migrations](https://github.com/pressly/goose)
- [Logrus](https://github.com/sirupsen/logrus)

## 💡 Tips for Effective Copilot Usage

### 1. Use Descriptive Comments
```go
// Create a new user with email validation and password hashing
// following the standard service pattern
func (s *UserService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
```

### 2. Reference Patterns
```go
// Implement GetByID following the repository pattern from user_repository.go
```

### 3. Request Tests
```go
// Write comprehensive tests for CreateUser including success, validation errors,
// and duplicate email cases
```

### 4. Specify Standards
```go
// Add logging with logrus, wrap errors with context, return proper status codes
```

## 🔄 Keeping Configuration Updated

When you make improvements to the architecture:

1. Update `copilot-instructions.md` with new patterns
2. Add new templates to `service-templates.md`
3. Document architectural decisions in `architecture-patterns.md`
4. Update this README

## 🤝 Contributing

When adding new patterns or improving existing ones:

1. Test the pattern in a real service first
2. Document it comprehensively
3. Add examples
4. Update related sections
5. Keep consistency across all three files

## 📞 Support

For questions or issues:

1. Check this README
2. Review the three configuration files
3. Look at existing code in user_service
4. Check project documentation files

## 🎉 Benefits of This Configuration

✅ **Consistency**: All services follow the same patterns
✅ **Quality**: Built-in best practices and standards
✅ **Speed**: Templates accelerate development
✅ **Testability**: 100% coverage is achievable
✅ **Maintainability**: Clear architecture and documentation
✅ **Learning**: New developers can reference examples
✅ **Copilot Enhancement**: AI assistant has full context

---

**Last Updated**: January 22, 2026
**Version**: 1.0
**Based On**: user_service v1.0 (100% test coverage achieved)
