# User Service API - Quick Reference

## Base URL
```
Local: http://localhost:8080
```

## Health Checks

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Liveness check |
| `/ready` | GET | Readiness check |

## Access Levels

| Endpoint | Method | Description | Request Body |
|----------|--------|-------------|--------------|
| `/access-levels` | POST | Create access level | `{"name": "Admin", "description": "..."}` |
| `/access-levels` | GET | List all access levels | - |
| `/access-levels/{id}` | GET | Get access level by ID | - |

## Users

| Endpoint | Method | Description | Request Body |
|----------|--------|-------------|--------------|
| `/users` | POST | Create user | See below |
| `/users` | GET | List users (paginated) | Query: `?page=1&page_size=10` |
| `/users/{id}` | GET | Get user by ID | - |
| `/users/{id}` | PUT | Update user | See below |
| `/users/{id}` | DELETE | Delete user | - |

### Create User Body
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "phone_number": "+1234567890",  // optional
  "password": "SecurePass123!"     // min 8 chars
}
```

### Update User Body
All fields are optional:
```json
{
  "first_name": "John",
  "last_name": "Updated",
  "email": "new.email@example.com",
  "phone_number": "+9876543210"
}
```

## User Access Levels

| Endpoint | Method | Description | Request Body |
|----------|--------|-------------|--------------|
| `/users/{id}/access-levels` | POST | Assign access levels | `{"access_level_ids": [1, 2]}` |
| `/users/{id}/access-levels` | GET | Get user's access levels | - |

## Authentication

| Endpoint | Method | Description | Request Body |
|----------|--------|-------------|--------------|
| `/auth/login` | POST | Login | `{"email": "...", "password": "..."}` |

## Response Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request (validation error) |
| 401 | Unauthorized (auth failed) |
| 404 | Not Found |
| 500 | Internal Server Error |

## Common Error Response Format
```json
{
  "error": "Error description",
  "message": "Detailed error message"
}
```

## Validation Rules

- **Email**: Must be valid email format, unique
- **Password**: Minimum 8 characters
- **First/Last Name**: 1-50 characters
- **Phone Number**: Optional, max 20 characters
- **User ID**: Must be valid UUID v4
- **Access Level ID**: Must be valid integer

## Example cURL Commands

### Create User
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "password": "SecurePass123!"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "SecurePass123!"
  }'
```

### List Users
```bash
curl -X GET "http://localhost:8080/users?page=1&page_size=10"
```

### Get User
```bash
curl -X GET http://localhost:8080/users/{user-id}
```

### Create Access Level
```bash
curl -X POST http://localhost:8080/access-levels \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Admin",
    "description": "Administrator access"
  }'
```

### Assign Access Levels
```bash
curl -X POST http://localhost:8080/users/{user-id}/access-levels \
  -H "Content-Type: application/json" \
  -d '{
    "access_level_ids": [1, 2]
  }'
```

## Testing Flow

1. ✅ Create Access Levels (Admin, User)
2. ✅ Create User(s)
3. ✅ Assign Access Levels to User
4. ✅ Login with User credentials
5. ✅ List/Get/Update Users
6. ✅ Delete User (cleanup)

## Tips

- Save user IDs and access level IDs from creation responses
- Use pagination for large user lists
- Test error cases with invalid data
- Check health endpoints before running tests
- Clean up test data when done
