# User Service API Documentation

This service provides RESTful APIs for user management, authentication, and access level management.

## Table of Contents
- [Getting Started](#getting-started)
- [API Endpoints](#api-endpoints)
  - [User Management](#user-management)
  - [Authentication](#authentication)
  - [Access Levels](#access-levels)
- [Data Models](#data-models)
- [Error Handling](#error-handling)

## Getting Started

### Prerequisites
- Go 1.25+
- PostgreSQL database
- All dependencies installed via `go mod download`

### Running the Service
```bash
# Run with local configuration
go run main.go --config=local

# Build and run
go build -o build/user_service.exe
./build/user_service.exe --config=local
```

The service will start on `http://0.0.0.0:8080` (configurable in `resources/local.yaml`)

## API Endpoints

### User Management

#### Create User
Create a new user with authentication credentials.

**Endpoint:** `POST /users`

**Request Body:**
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "phone_number": "+1234567890",
  "password": "securepassword123"
}
```

**Response:** `201 Created`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "phone_number": "+1234567890",
  "access_levels": [],
  "created_at": "2026-01-17T10:30:00Z",
  "updated_at": "2026-01-17T10:30:00Z"
}
```

**Validation Rules:**
- `first_name`: Required, 1-50 characters
- `last_name`: Required, 1-50 characters
- `email`: Required, valid email format, max 255 characters
- `phone_number`: Optional, max 20 characters
- `password`: Required, minimum 8 characters

---

#### Get User by ID
Retrieve a specific user by their ID.

**Endpoint:** `GET /users/{id}`

**Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "phone_number": "+1234567890",
  "access_levels": [
    {
      "id": 1,
      "name": "admin",
      "description": "Administrator access"
    }
  ],
  "created_at": "2026-01-17T10:30:00Z",
  "updated_at": "2026-01-17T10:30:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid user ID format
- `404 Not Found`: User not found

---

#### Update User
Update an existing user's information.

**Endpoint:** `PUT /users/{id}`

**Request Body:** (all fields optional)
```json
{
  "first_name": "Jane",
  "last_name": "Smith",
  "email": "jane.smith@example.com",
  "phone_number": "+9876543210"
}
```

**Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "first_name": "Jane",
  "last_name": "Smith",
  "email": "jane.smith@example.com",
  "phone_number": "+9876543210",
  "access_levels": [],
  "created_at": "2026-01-17T10:30:00Z",
  "updated_at": "2026-01-17T11:45:00Z"
}
```

---

#### Delete User (Soft Delete)
Soft delete a user by setting their `deleted_at` timestamp.

**Endpoint:** `DELETE /users/{id}`

**Response:** `200 OK`
```json
{
  "message": "User deleted successfully"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid user ID format
- `404 Not Found`: User not found

---

#### List Users
Retrieve a paginated list of all users.

**Endpoint:** `GET /users`

**Query Parameters:**
- `page`: Page number (default: 1)
- `page_size`: Number of users per page (default: 10, max: 100)

**Example:** `GET /users?page=1&page_size=20`

**Response:** `200 OK`
```json
{
  "users": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "first_name": "John",
      "last_name": "Doe",
      "email": "john.doe@example.com",
      "phone_number": "+1234567890",
      "access_levels": [],
      "created_at": "2026-01-17T10:30:00Z",
      "updated_at": "2026-01-17T10:30:00Z"
    }
  ],
  "total": 45,
  "page": 1,
  "page_size": 20
}
```

---

### Authentication

#### Login
Authenticate a user with email and password.

**Endpoint:** `POST /auth/login`

**Request Body:**
```json
{
  "email": "john.doe@example.com",
  "password": "securepassword123"
}
```

**Response:** `200 OK`
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "phone_number": "+1234567890",
    "access_levels": [
      {
        "id": 1,
        "name": "admin",
        "description": "Administrator access"
      }
    ],
    "created_at": "2026-01-17T10:30:00Z",
    "updated_at": "2026-01-17T10:30:00Z"
  },
  "message": "Login successful"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Invalid email or password

---

### Access Levels

#### Create Access Level
Create a new access level.

**Endpoint:** `POST /access-levels`

**Request Body:**
```json
{
  "name": "admin",
  "description": "Administrator access with full permissions"
}
```

**Response:** `201 Created`
```json
{
  "id": 1,
  "name": "admin",
  "description": "Administrator access with full permissions"
}
```

---

#### Get Access Level by ID
Retrieve a specific access level.

**Endpoint:** `GET /access-levels/{id}`

**Response:** `200 OK`
```json
{
  "id": 1,
  "name": "admin",
  "description": "Administrator access with full permissions"
}
```

---

#### List All Access Levels
Retrieve all available access levels.

**Endpoint:** `GET /access-levels`

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "name": "admin",
    "description": "Administrator access"
  },
  {
    "id": 2,
    "name": "user",
    "description": "Standard user access"
  }
]
```

---

#### Assign Access Levels to User
Assign one or more access levels to a user.

**Endpoint:** `POST /users/{id}/access-levels`

**Request Body:**
```json
{
  "access_level_ids": [1, 2, 3]
}
```

**Response:** `200 OK`
```json
{
  "message": "Access levels assigned successfully"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid user ID or access level IDs
- `404 Not Found`: User or access level not found

---

#### Get User Access Levels
Retrieve all access levels assigned to a specific user.

**Endpoint:** `GET /users/{id}/access-levels`

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "name": "admin",
    "description": "Administrator access"
  },
  {
    "id": 2,
    "name": "user",
    "description": "Standard user access"
  }
]
```

---

## Data Models

### User
```go
{
  "id": "UUID",
  "first_name": "string",
  "last_name": "string",
  "email": "string",
  "phone_number": "string (optional)",
  "access_levels": "array of AccessLevel (optional)",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### AccessLevel
```go
{
  "id": "integer",
  "name": "string",
  "description": "string (optional)"
}
```

---

## Error Handling

All error responses follow this format:

```json
{
  "error": "Error type",
  "message": "Detailed error message"
}
```

### Common HTTP Status Codes
- `200 OK`: Request succeeded
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Authentication failed
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

---

## Database Schema

The service uses the following PostgreSQL tables:

### users
- `id` (UUID, primary key)
- `first_name` (VARCHAR(50), required)
- `last_name` (VARCHAR(50), required)
- `email` (VARCHAR(255), unique, required)
- `phone_number` (VARCHAR(20), optional)
- `created_at` (TIMESTAMPTZ)
- `updated_at` (TIMESTAMPTZ)
- `deleted_at` (TIMESTAMPTZ, nullable - for soft deletes)

### user_authentications
- `id` (UUID, primary key)
- `user_id` (UUID, foreign key to users)
- `password_hash` (VARCHAR(255), bcrypt hashed)
- `created_at` (TIMESTAMPTZ)
- `updated_at` (TIMESTAMPTZ)
- `deleted_at` (TIMESTAMPTZ, nullable)

### access_levels
- `id` (SERIAL, primary key)
- `name` (VARCHAR(50), unique, required)
- `description` (TEXT, optional)
- `created_at` (TIMESTAMPTZ)
- `updated_at` (TIMESTAMPTZ)
- `deleted_at` (TIMESTAMPTZ, nullable)

### user_access_levels
- `user_id` (UUID, foreign key to users)
- `access_level_id` (INTEGER, foreign key to access_levels)
- `created_at` (TIMESTAMPTZ)
- `updated_at` (TIMESTAMPTZ)
- `deleted_at` (TIMESTAMPTZ, nullable)
- Primary key: (user_id, access_level_id)

---

## Security Notes

1. **Password Storage**: Passwords are hashed using bcrypt with default cost before storage
2. **Soft Deletes**: Users are soft-deleted (deleted_at is set) rather than permanently removed
3. **Email Uniqueness**: Email addresses must be unique across all active users
4. **Input Validation**: All inputs are validated before processing

---

## Testing with cURL

### Create a User
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "phone_number": "+1234567890",
    "password": "securepass123"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "securepass123"
  }'
```

### List Users
```bash
curl -X GET "http://localhost:8080/users?page=1&page_size=10"
```

### Get User by ID
```bash
curl -X GET http://localhost:8080/users/{user-id}
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

### Assign Access Levels to User
```bash
curl -X POST http://localhost:8080/users/{user-id}/access-levels \
  -H "Content-Type: application/json" \
  -d '{
    "access_level_ids": [1, 2]
  }'
```
