#!/bin/bash

# User Service API Test Script
# This script tests all the implemented API endpoints

BASE_URL="http://localhost:8080"

echo "========================================="
echo "User Service API Test Script"
echo "========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Function to test endpoint
test_endpoint() {
    local name=$1
    local method=$2
    local endpoint=$3
    local data=$4

    echo -e "${BLUE}Testing: $name${NC}"

    if [ -z "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data")
    fi

    http_code=$(echo "$response" | tail -n 1)
    body=$(echo "$response" | head -n -1)

    echo "HTTP Status: $http_code"
    echo "Response: $body"
    echo ""

    if [ $http_code -ge 200 ] && [ $http_code -lt 300 ]; then
        echo -e "${GREEN}✓ PASSED${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "${RED}✗ FAILED${NC}"
        ((TESTS_FAILED++))
    fi
    echo "========================================="
    echo ""
}

# Wait for service to be ready
echo "Checking if service is running..."
if ! curl -s "$BASE_URL/health" > /dev/null; then
    echo -e "${RED}Service is not running at $BASE_URL${NC}"
    echo "Please start the service first with: ./build/user_service.exe --config=local"
    exit 1
fi
echo -e "${GREEN}Service is running!${NC}"
echo ""

# Test 1: Health Check
test_endpoint "Health Check" "GET" "/health"

# Test 2: Readiness Check
test_endpoint "Readiness Check" "GET" "/ready"

# Test 3: Create Access Level - Admin
test_endpoint "Create Access Level - Admin" "POST" "/access-levels" \
'{
  "name": "admin",
  "description": "Administrator access with full permissions"
}'

# Test 4: Create Access Level - User
test_endpoint "Create Access Level - User" "POST" "/access-levels" \
'{
  "name": "user",
  "description": "Standard user access"
}'

# Test 5: List Access Levels
test_endpoint "List Access Levels" "GET" "/access-levels"

# Test 6: Create User - John Doe
test_endpoint "Create User - John Doe" "POST" "/users" \
'{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "phone_number": "+1234567890",
  "password": "securepass123"
}'

# Extract user ID from response (you may need to adjust this based on actual response)
USER_ID=$(curl -s -X POST "$BASE_URL/users" \
    -H "Content-Type: application/json" \
    -d '{"first_name":"Jane","last_name":"Smith","email":"jane.smith@example.com","password":"password123"}' | \
    grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

# Test 7: Create Another User
test_endpoint "Create User - Jane Smith" "POST" "/users" \
'{
  "first_name": "Jane",
  "last_name": "Smith",
  "email": "jane.smith2@example.com",
  "phone_number": "+9876543210",
  "password": "password123"
}'

# Test 8: Login
test_endpoint "Login with Valid Credentials" "POST" "/auth/login" \
'{
  "email": "john.doe@example.com",
  "password": "securepass123"
}'

# Test 9: List Users
test_endpoint "List Users (Page 1)" "GET" "/users?page=1&page_size=10"

# Test 10: Get User by ID (if we have USER_ID)
if [ ! -z "$USER_ID" ]; then
    test_endpoint "Get User by ID" "GET" "/users/$USER_ID"

    # Test 11: Update User
    test_endpoint "Update User" "PUT" "/users/$USER_ID" \
    '{
      "first_name": "Janet",
      "last_name": "Smithson"
    }'

    # Test 12: Assign Access Levels
    test_endpoint "Assign Access Levels to User" "POST" "/users/$USER_ID/access-levels" \
    '{
      "access_level_ids": [1, 2]
    }'

    # Test 13: Get User Access Levels
    test_endpoint "Get User Access Levels" "GET" "/users/$USER_ID/access-levels"
else
    echo -e "${RED}Could not extract user ID, skipping user-specific tests${NC}"
fi

# Test 14: Test Invalid Login
test_endpoint "Login with Invalid Credentials (should fail)" "POST" "/auth/login" \
'{
  "email": "john.doe@example.com",
  "password": "wrongpassword"
}'

# Summary
echo ""
echo "========================================="
echo "Test Summary"
echo "========================================="
echo -e "Tests Passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "Tests Failed: ${RED}$TESTS_FAILED${NC}"
echo "========================================="

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed!${NC}"
    exit 1
fi
