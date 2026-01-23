# GitHub Copilot Configuration - Quick Reference Card

## 📁 Configuration Files Location
```
user_service/.github/
├── README.md                    - Start here! Navigation and overview
├── copilot-instructions.md      - Primary config (799 lines)
├── service-templates.md         - Code templates (645 lines)
└── architecture-patterns.md     - Deep dive (754 lines)
```

## 🚀 Quick Actions

### Create New Microservice
```bash
cp -r user_service new_service
cd new_service
# Update go.mod module name
# Follow .github/service-templates.md
```

### Add New Entity
```bash
# 1. Open .github/service-templates.md
# 2. Copy template (Model, DTO, Repository, Service, Handler)
# 3. Replace: {{ENTITY_NAME}}, {{ENTITY_NAME_LOWER}}, {{TABLE_NAME}}
# 4. Create files in: models/, dto/, repository/, service/, handlers/
# 5. Write tests for each layer
# 6. Wire up in cmd/app.go
```

### Use with Copilot
```go
// Just write descriptive comments, Copilot will reference config files:

// Create a Product entity following the standard pattern

// Implement ProductRepository with PostgreSQL following the pattern

// Write comprehensive tests for CreateProduct including error cases
```

## 📚 File Purpose Summary

| File | When to Use |
|------|-------------|
| **README.md** | First time reading, navigation |
| **copilot-instructions.md** | Writing code, checking standards |
| **service-templates.md** | Creating new entities/services |
| **architecture-patterns.md** | Understanding design decisions |

## 🎯 Key Patterns

### Architecture Layers
```
Handlers (HTTP) → Services (Logic) → Repositories (Data) → Database
     ↓                 ↓                    ↓
   DTOs            Business             Models
                    Logic
```

### Always Use Interfaces
```go
// Service layer
type UserServiceInterface interface { ... }
var _ UserServiceInterface = (*UserService)(nil)

// Handlers use interfaces
type UserHandler struct {
    service UserServiceInterface // Not *UserService!
}
```

### Testing Pattern
```go
// Repository: Use SQLite in-memory
// Service: Mock repository
// Handler: Mock service
// Target: 100% coverage
```

## ✅ Code Quality Checklist

- [ ] All errors wrapped with context: `fmt.Errorf("failed: %w", err)`
- [ ] All functions accept `context.Context`
- [ ] All services/repos use interfaces
- [ ] All code has unit tests
- [ ] All passwords hashed with bcrypt
- [ ] All inputs validated in service layer
- [ ] All sensitive data protected
- [ ] All HTTP responses are JSON
- [ ] All endpoints return proper status codes
- [ ] All logs use logrus (not fmt.Println)

## 🔧 Commands

```bash
# Run tests
go test ./... -v

# With coverage
go test ./... -cover -coverprofile=coverage.out

# View coverage
go tool cover -html=coverage.out

# Without CGO (Windows)
CGO_ENABLED=0 go test ./... -v

# Build
go build -o build/service .

# Run
./build/service -config=local

# Migrations
goose -dir database/migrations postgres "connection_string" up
```

## 🛠️ Tech Stack

- Go 1.25+
- Gorilla Mux (routing)
- GORM (ORM)
- PostgreSQL (database)
- Goose (migrations)
- testify (testing)
- logrus (logging)
- YAML (config)

## 📊 File Statistics

- **Total Lines**: 2,550
- **Code Examples**: 50+
- **Patterns**: 15+
- **Templates**: 8
- **Size**: ~80 KB

## 💡 Pro Tips

1. **Always start with** `.github/README.md`
2. **Copy templates** from `service-templates.md`
3. **Reference patterns** in `architecture-patterns.md`
4. **Check standards** in `copilot-instructions.md`
5. **Write descriptive comments** for Copilot
6. **Follow naming conventions** consistently
7. **Test everything** (100% coverage target)
8. **Use interfaces** for all services/repos

## 🎯 Layer Responsibilities

### Models
- GORM entities only
- No business logic
- TableName() method

### DTOs
- API contract
- JSON tags
- Validation tags
- No logic

### Repository
- Interface + Implementation
- CRUD operations
- Return models
- Wrap errors

### Service
- Interface + Implementation
- Business logic
- Validation
- Return DTOs
- Orchestrate repos

### Handlers
- HTTP only
- Parse requests
- Call services
- Format responses
- Set status codes

## 🔒 Security Essentials

```go
// Hash passwords
hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// Validate input
if req.Name == "" {
    return fmt.Errorf("name is required")
}

// Never log sensitive data
logrus.Infof("User created: ID=%s", user.ID) // ✅
logrus.Infof("User: %+v", user)              // ❌
```

## 📈 Testing Coverage

| Layer | Method | Target |
|-------|--------|--------|
| Repository | SQLite in-memory | 100% |
| Service | Mock repository | 100% |
| Handler | Mock service | 100% |
| **Overall** | - | **100%** |

## 🎓 Learning Path

1. Read `.github/README.md` (10 mins)
2. Skim `copilot-instructions.md` (20 mins)
3. Try creating entity with `service-templates.md` (30 mins)
4. Deep dive `architecture-patterns.md` (60 mins)
5. Practice with real entity (120 mins)

**Total**: ~4 hours to mastery

## 🌟 Success Criteria

- [ ] Can create new entity in < 1 hour
- [ ] Can explain architecture layers
- [ ] Can write tests with 100% coverage
- [ ] Can use Copilot effectively
- [ ] Code follows all standards

## 📞 Help

1. Check `.github/README.md`
2. Search configuration files
3. Look at user_service examples
4. Review project documentation

---

**Quick Start**: Open `.github/README.md` and follow the guide!

**Status**: ✅ Ready to use
**Version**: 1.0
**Last Updated**: January 22, 2026
