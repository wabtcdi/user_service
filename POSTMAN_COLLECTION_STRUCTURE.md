# Postman Collection Structure

## 📁 Collection: User Service API

```
User Service API
│
├── 📁 Health Checks
│   ├── GET  /health                    (Liveness Check)
│   └── GET  /ready                     (Readiness Check)
│
├── 📁 Access Levels
│   ├── POST /access-levels             (Create Access Level) ⭐
│   ├── POST /access-levels             (Create Access Level - User) ⭐
│   ├── GET  /access-levels/{id}        (Get Access Level)
│   └── GET  /access-levels             (List Access Levels)
│
├── 📁 Users
│   ├── POST   /users                   (Create User) ⭐
│   ├── POST   /users                   (Create User - Without Phone)
│   ├── GET    /users/{id}              (Get User)
│   ├── PUT    /users/{id}              (Update User)
│   ├── GET    /users?page&page_size    (List Users)
│   ├── POST   /users/{id}/access-levels (Assign Access Levels)
│   ├── GET    /users/{id}/access-levels (Get User Access Levels)
│   └── DELETE /users/{id}              (Delete User)
│
├── 📁 Authentication
│   ├── POST /auth/login                (Login - Success)
│   └── POST /auth/login                (Login - Invalid Credentials)
│
└── 📁 Error Cases
    ├── POST /users                     (Create User - Invalid Email)
    ├── POST /users                     (Create User - Short Password)
    ├── GET  /users/invalid-uuid        (Get User - Invalid UUID)
    ├── GET  /users/00000...000         (Get User - Not Found)
    └── POST /users/{id}/access-levels  (Assign Access Levels - Empty Array)
```

⭐ = Requests with automated test scripts that extract and save IDs

## 🔄 Request Flow

### Typical Testing Sequence

```
1. Health Checks
   └─> Verify service is running and database is ready

2. Create Access Levels
   ├─> POST /access-levels (Admin)
   │   └─> Saves: access_level_id
   └─> POST /access-levels (User)
       └─> Saves: access_level_id_user

3. Create Users
   ├─> POST /users (with phone)
   │   └─> Saves: user_id
   └─> POST /users (without phone)

4. Assign Access Levels
   └─> POST /users/{user_id}/access-levels
       └─> Uses: user_id, access_level_id, access_level_id_user

5. Verify Assignment
   └─> GET /users/{user_id}/access-levels
       └─> Uses: user_id

6. Authentication
   └─> POST /auth/login
       └─> Uses credentials from user creation

7. User Operations
   ├─> GET /users/{user_id}
   ├─> PUT /users/{user_id}
   ├─> GET /users?page=1&page_size=10
   └─> DELETE /users/{user_id}

8. Error Testing
   └─> Run all error cases to verify validation
```

## 📊 Request Details

### HTTP Methods Distribution
- **GET**: 8 requests (read operations)
- **POST**: 10 requests (create operations)
- **PUT**: 1 request (update operations)
- **DELETE**: 1 request (delete operations)
- **Total**: 21 requests

### Status Codes Expected
- **200 OK**: 12 requests
- **201 Created**: 4 requests
- **400 Bad Request**: 5 error cases
- **401 Unauthorized**: 1 error case
- **404 Not Found**: 2 error cases

### Request Body Types
- **JSON**: 13 requests with body
- **No Body**: 8 GET/DELETE requests

### Response Types
- **JSON**: All responses
- **Structure**: Consistent error format

## 🔧 Environment Variables

### Collection Variables
```
base_url = http://localhost:8080
```

### Dynamic Variables (Auto-populated)
```
user_id                 = (UUID from user creation)
access_level_id         = (ID from Admin access level)
access_level_id_user    = (ID from User access level)
```

### Usage in Requests
```
{{base_url}}/users/{{user_id}}
{{base_url}}/access-levels/{{access_level_id}}
```

## 📝 Request Examples

### 1. Create Access Level
```http
POST {{base_url}}/access-levels
Content-Type: application/json

{
    "name": "Admin",
    "description": "Administrator access level with full permissions"
}
```
**Response**: 201 Created + ID saved

### 2. Create User
```http
POST {{base_url}}/users
Content-Type: application/json

{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "phone_number": "+1234567890",
    "password": "SecurePass123!"
}
```
**Response**: 201 Created + user_id saved

### 3. Assign Access Levels
```http
POST {{base_url}}/users/{{user_id}}/access-levels
Content-Type: application/json

{
    "access_level_ids": [{{access_level_id}}, {{access_level_id_user}}]
}
```
**Response**: 200 OK

### 4. Login
```http
POST {{base_url}}/auth/login
Content-Type: application/json

{
    "email": "john.doe@example.com",
    "password": "SecurePass123!"
}
```
**Response**: 200 OK with user data

## 🧪 Test Scripts

### User Creation Script
```javascript
if (pm.response.code === 201) {
    var jsonData = pm.response.json();
    pm.environment.set("user_id", jsonData.id);
    pm.test("User created successfully", function () {
        pm.response.to.have.status(201);
    });
    pm.test("Response has user id", function () {
        pm.expect(jsonData).to.have.property('id');
    });
}
```

### Access Level Creation Script
```javascript
if (pm.response.code === 201) {
    var jsonData = pm.response.json();
    pm.environment.set("access_level_id", jsonData.id);
}
```

## 🎯 Testing Goals

### Functional Testing
- ✅ Verify all endpoints work correctly
- ✅ Test request/response formats
- ✅ Validate business logic
- ✅ Check data persistence

### Validation Testing
- ✅ Email format validation
- ✅ Password length validation
- ✅ UUID format validation
- ✅ Required field validation
- ✅ Array validation

### Error Handling
- ✅ Invalid input handling
- ✅ Not found scenarios
- ✅ Authentication failures
- ✅ Proper error responses

### Integration Testing
- ✅ Multi-step workflows
- ✅ Data relationships
- ✅ Access level assignments
- ✅ Authentication flow

## 📦 Files Summary

| File | Size | Purpose |
|------|------|---------|
| `postman_collection.json` | ~14KB | Main collection with all requests |
| `postman_environment_local.json` | ~700B | Local dev environment |
| `postman_environment_cloud.json` | ~700B | Cloud environment template |
| `POSTMAN_GUIDE.md` | ~8KB | Comprehensive usage guide |
| `POSTMAN_TESTING_SUMMARY.md` | ~8KB | Overview and best practices |
| `API_QUICK_REFERENCE.md` | ~4KB | Quick reference card |

**Total**: 6 files, ~36KB

## 🚀 Quick Start Commands

### Import Collection
1. Open Postman
2. File → Import
3. Select `postman_collection.json`

### Import Environment
1. File → Import
2. Select `postman_environment_local.json`
3. Select environment from dropdown

### Run Collection
1. Click collection name
2. Click **Run**
3. Select requests to run
4. Click **Run User Service API**

## 💡 Tips

1. **Run in order**: Follow the recommended sequence for best results
2. **Check variables**: Verify IDs are saved after creation
3. **Use environments**: Switch between local/cloud easily
4. **Export results**: Save test run reports
5. **Share with team**: Export and share collection/environments
6. **Document changes**: Add notes to requests as needed

## 🎊 Ready to Test!

Import the collection and start testing your User Service API immediately!

---

**Collection Version**: 1.0  
**Last Updated**: January 17, 2026  
**Compatible With**: Postman 10.0+
