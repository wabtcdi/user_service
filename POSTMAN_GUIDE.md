# Postman Collection Guide

This guide explains how to use the Postman collection to test all endpoints in the User Service API.

## Import the Collection

1. Open Postman
2. Click **Import** button (top left)
3. Select **File** tab
4. Choose the `postman_collection.json` file from this directory
5. Click **Import**

## Environment Variables

The collection uses the following variables:

### Collection Variable
- **base_url**: `http://localhost:8080` (default server URL)

### Environment Variables (Auto-set by tests)
- **user_id**: Automatically set when a user is created
- **access_level_id**: Automatically set when an access level is created
- **access_level_id_user**: Automatically set when a second access level is created

You can modify the `base_url` in the collection variables if your server runs on a different host/port.

## Testing Workflow

### Recommended Testing Order

Follow this sequence for a complete test flow:

#### 1. Health Checks
- **Liveness Check** - Verify service is running
- **Readiness Check** - Verify service is ready

#### 2. Access Levels (Setup)
- **Create Access Level** (Admin) - Creates "Admin" access level
- **Create Access Level - User** - Creates "User" access level
- **List Access Levels** - View all created access levels
- **Get Access Level** - Get specific access level by ID

#### 3. Users
- **Create User** - Creates a user with phone number
- **Create User - Without Phone** - Creates a user without phone
- **Get User** - Retrieve user by ID
- **List Users** - Get paginated list of all users
- **Update User** - Update user information
- **Assign Access Levels** - Assign both access levels to the user
- **Get User Access Levels** - View user's assigned access levels

#### 4. Authentication
- **Login - Success** - Test successful authentication
- **Login - Invalid Credentials** - Test failed authentication

#### 5. User Cleanup (Optional)
- **Delete User** - Soft delete the user

#### 6. Error Cases
Test various error scenarios:
- **Create User - Invalid Email** - Test email validation
- **Create User - Short Password** - Test password length validation
- **Get User - Invalid UUID** - Test UUID format validation
- **Get User - Not Found** - Test 404 error handling
- **Assign Access Levels - Empty Array** - Test array validation

## Folder Structure

The collection is organized into logical folders:

### 1. Health Checks
- Liveness and readiness probes for Kubernetes

### 2. Access Levels
- CRUD operations for access levels
- List all access levels

### 3. Users
- Full user lifecycle: Create, Read, Update, Delete
- User listing with pagination
- Access level assignment and retrieval

### 4. Authentication
- Login endpoint with success and failure scenarios

### 5. Error Cases
- Comprehensive error handling tests
- Input validation tests

## Request Examples

### Create User
```json
POST /users
{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "phone_number": "+1234567890",
    "password": "SecurePass123!"
}
```

### Update User
```json
PUT /users/{id}
{
    "first_name": "John",
    "last_name": "Updated",
    "phone_number": "+9876543210"
}
```

### Login
```json
POST /auth/login
{
    "email": "john.doe@example.com",
    "password": "SecurePass123!"
}
```

### Assign Access Levels
```json
POST /users/{id}/access-levels
{
    "access_level_ids": [1, 2]
}
```

### Create Access Level
```json
POST /access-levels
{
    "name": "Admin",
    "description": "Administrator access level with full permissions"
}
```

### List Users with Pagination
```
GET /users?page=1&page_size=10
```

## Test Scripts

The collection includes automated test scripts that:

1. **Extract IDs**: Automatically save user and access level IDs from responses
2. **Validate Status Codes**: Check for expected HTTP status codes
3. **Verify Response Structure**: Ensure response contains required fields

These scripts enable chaining requests - IDs from creation requests are automatically used in subsequent requests.

## Response Validation

### Successful Responses

#### User Created (201)
```json
{
    "id": "uuid-here",
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "phone_number": "+1234567890",
    "created_at": "2026-01-17T10:00:00Z",
    "updated_at": "2026-01-17T10:00:00Z"
}
```

#### Login Success (200)
```json
{
    "user": {
        "id": "uuid-here",
        "first_name": "John",
        "last_name": "Doe",
        "email": "john.doe@example.com",
        "access_levels": [...],
        "created_at": "2026-01-17T10:00:00Z",
        "updated_at": "2026-01-17T10:00:00Z"
    },
    "message": "Login successful"
}
```

#### List Users (200)
```json
{
    "users": [...],
    "total": 2,
    "page": 1,
    "page_size": 10
}
```

### Error Responses

#### Bad Request (400)
```json
{
    "error": "Invalid request body",
    "message": "email: invalid format"
}
```

#### Unauthorized (401)
```json
{
    "error": "Authentication failed",
    "message": "invalid credentials"
}
```

#### Not Found (404)
```json
{
    "error": "User not found",
    "message": "no user found with id: ..."
}
```

## Tips

1. **Run in Sequence**: Execute requests in the recommended order for best results
2. **Check Variables**: After creating users/access levels, verify the environment variables are set
3. **Database State**: The service uses PostgreSQL - ensure it's running and accessible
4. **Pagination**: Test different page numbers and page sizes with the List Users endpoint
5. **Validation**: Use the Error Cases folder to understand input validation rules
6. **Cleanup**: Delete test users when done to keep the database clean

## Troubleshooting

### Connection Refused
- Ensure the service is running on `localhost:8080`
- Check if the database is accessible
- Verify firewall settings

### 404 Not Found on All Endpoints
- Check the base_url in collection variables
- Ensure the service started successfully
- Review server logs for errors

### Variable Not Set
- Run the creation requests first (Create User, Create Access Level)
- Check the test scripts executed successfully
- Manually set variables if needed in the environment tab

### Database Errors
- Ensure PostgreSQL is running
- Check database connection settings in `resources/local.yaml`
- Verify migrations ran successfully

## Advanced Usage

### Running the Entire Collection

Use Postman's Collection Runner:
1. Click on the collection name
2. Click **Run** button
3. Select the requests you want to run
4. Click **Run User Service API**

This will execute all requests in sequence and show a summary of results.

### Exporting Results

After running the collection:
1. Click **Export Results**
2. Save as JSON or HTML
3. Share with your team or include in reports

## API Documentation

For detailed API documentation, see:
- `API_DOCUMENTATION.md` - Complete API reference
- `README_COMPLETE.md` - Service overview and setup guide
- `QUICK_START.md` - Quick start instructions

## Notes

- All passwords must be at least 8 characters long
- Email addresses must be valid and unique
- UUIDs must be valid v4 format
- Phone numbers are optional but should not exceed 20 characters
- Access level IDs must exist before assignment
- Deleted users/access levels are soft-deleted (deleted_at is set)
