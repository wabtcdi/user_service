# User Service - Quick Start Guide

## Project Structure

```
user_service/
├── models/              # Domain models (User, AccessLevel, etc.)
├── dto/                 # Data Transfer Objects for API requests/responses
├── repository/          # Database access layer (PostgreSQL)
├── service/             # Business logic layer
├── handlers/            # HTTP handlers (REST API endpoints)
├── cmd/                 # Application initialization and configuration
├── database/migrations/ # Database migrations (Goose)
├── resources/           # Configuration files (local.yaml, cloud.yaml)
├── build/               # Compiled binaries
└── API_DOCUMENTATION.md # Complete API documentation
```

## Quick Start

### 1. Prerequisites
```bash
# Ensure PostgreSQL is running
# Default connection: localhost:30432 (Minikube)
# Database: userdb
# User: postgres
# Password: postgres123
```

### 2. Build the Service
```bash
cd C:\Users\19802\GolandProjects\user_service
go build -o build/user_service.exe
```

### 3. Run the Service
```bash
# Run with local configuration
./build/user_service.exe --config=local

# Or run directly with go
go run main.go --config=local
```

### 4. Test the Service
```bash
# Make the test script executable
chmod +x test_api.sh

# Run all API tests
./test_api.sh
```

## API Endpoints Quick Reference

### Users
- `POST /users` - Create user
- `GET /users` - List users (with pagination)
- `GET /users/{id}` - Get user by ID
- `PUT /users/{id}` - Update user
- `DELETE /users/{id}` - Delete user

### Authentication
- `POST /auth/login` - Login

### Access Levels
- `POST /access-levels` - Create access level
- `GET /access-levels` - List all access levels
- `GET /access-levels/{id}` - Get access level
- `POST /users/{id}/access-levels` - Assign access levels to user
- `GET /users/{id}/access-levels` - Get user's access levels

### Health Checks
- `GET /health` - Liveness check
- `GET /ready` - Readiness check

## Example Usage

### Create a User
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### List Users
```bash
curl "http://localhost:8080/users?page=1&page_size=10"
```

### Create Access Level
```bash
curl -X POST http://localhost:8080/access-levels \
  -H "Content-Type: application/json" \
  -d '{
    "name": "admin",
    "description": "Administrator access"
  }'
```

## Configuration

Edit `resources/local.yaml` to change:
- Server host/port
- Database connection details
- Logging level (debug, info, warn, error)
- Resource limits

## Database Migrations

Migrations are automatically run on service startup using Goose.

Location: `database/migrations/20260115015058_initial_setup.sql`

Tables created:
- `users` - User information
- `user_authentications` - Password hashes
- `access_levels` - Access level definitions
- `user_access_levels` - User-to-access-level mappings

## Troubleshooting

### Service won't start
1. Check PostgreSQL is running: `psql -h 192.168.59.102 -p 30432 -U postgres -d userdb`
2. Check configuration in `resources/local.yaml`
3. Check logs for specific error messages

### Database connection failed
1. Verify PostgreSQL connection details
2. Check network connectivity to database
3. Ensure database `userdb` exists
4. Verify user credentials

### API returns 500 errors
1. Check service logs for stack traces
2. Verify database migrations ran successfully
3. Check database constraints are satisfied

## Development

### Adding New Endpoints
1. Add DTO in `dto/user_dto.go`
2. Add service method in `service/user_service.go`
3. Add handler method in `handlers/user_handler.go`
4. Register route in `cmd/app.go` createRouter()

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./service/...
```

## Security Notes

- Passwords are hashed with bcrypt (cost 10)
- All operations support soft delete
- Email addresses must be unique
- Minimum password length: 8 characters

## Performance

- Pagination: Default 10 items, max 100 per page
- Database indexes on:
  - users.email
  - user_authentications.user_id
  - user_access_levels.user_id
  - user_access_levels.access_level_id

## Documentation

- **API_DOCUMENTATION.md** - Complete API reference with examples
- **IMPLEMENTATION_SUMMARY.md** - Technical implementation details
- **QUICK_START.md** - This file

## Support

For issues or questions:
1. Check the API_DOCUMENTATION.md
2. Review IMPLEMENTATION_SUMMARY.md for architecture details
3. Check service logs for errors
