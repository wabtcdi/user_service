package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wabtcdi/user_service/dto"
	"github.com/wabtcdi/user_service/service"
)

// MockUserService is a mock implementation of the UserServiceInterface
type MockUserService struct {
	mock.Mock
}

// Ensure MockUserService implements service.UserServiceInterface
var _ service.UserServiceInterface = (*MockUserService)(nil)

func (m *MockUserService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserService) GetUser(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, id uuid.UUID, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserService) ListUsers(ctx context.Context, page, pageSize int) (*dto.ListUsersResponse, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ListUsersResponse), args.Error(1)
}

func (m *MockUserService) AuthenticateUser(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LoginResponse), args.Error(1)
}

func (m *MockUserService) AssignAccessLevels(ctx context.Context, userID uuid.UUID, req *dto.AssignAccessLevelRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *MockUserService) GetUserAccessLevels(ctx context.Context, userID uuid.UUID) ([]dto.AccessLevelResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.AccessLevelResponse), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		userID := uuid.New()
		req := &dto.CreateUserRequest{
			FirstName:   "John",
			LastName:    "Doe",
			Email:       "john.doe@example.com",
			PhoneNumber: "1234567890",
			Password:    "password123",
		}

		expectedResponse := &dto.UserResponse{
			ID:        userID,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		}

		mockService.On("CreateUser", mock.Anything, req).Return(expectedResponse, nil)

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		recorder := httptest.NewRecorder()

		handler.CreateUser(recorder, request)

		assert.Equal(t, http.StatusCreated, recorder.Code)
		var response dto.UserResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, expectedResponse.ID, response.ID)
		assert.Equal(t, expectedResponse.Email, response.Email)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		request := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte("invalid json")))
		recorder := httptest.NewRecorder()

		handler.CreateUser(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		var response dto.ErrorResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, "Invalid request body", response.Error)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		req := &dto.CreateUserRequest{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Password:  "password123",
		}

		mockService.On("CreateUser", mock.Anything, req).Return(nil, errors.New("user already exists"))

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		recorder := httptest.NewRecorder()

		handler.CreateUser(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		var response dto.ErrorResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, "Failed to create user", response.Error)
		mockService.AssertExpectations(t)
	})
}

func TestGetUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		userID := uuid.New()
		expectedResponse := &dto.UserResponse{
			ID:        userID,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		}

		mockService.On("GetUser", mock.Anything, userID).Return(expectedResponse, nil)

		request := httptest.NewRequest(http.MethodGet, "/users/"+userID.String(), nil)
		recorder := httptest.NewRecorder()

		// Setup mux to parse path variables
		router := mux.NewRouter()
		router.HandleFunc("/users/{id}", handler.GetUser)
		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		var response dto.UserResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, expectedResponse.ID, response.ID)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid User ID", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		request := httptest.NewRequest(http.MethodGet, "/users/invalid-id", nil)
		recorder := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/users/{id}", handler.GetUser)
		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		var response dto.ErrorResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, "Invalid user ID", response.Error)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		userID := uuid.New()
		mockService.On("GetUser", mock.Anything, userID).Return(nil, errors.New("user not found"))

		request := httptest.NewRequest(http.MethodGet, "/users/"+userID.String(), nil)
		recorder := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/users/{id}", handler.GetUser)
		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
		var response dto.ErrorResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, "User not found", response.Error)
		mockService.AssertExpectations(t)
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		userID := uuid.New()
		req := &dto.UpdateUserRequest{
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "jane.doe@example.com",
		}

		expectedResponse := &dto.UserResponse{
			ID:        userID,
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "jane.doe@example.com",
		}

		mockService.On("UpdateUser", mock.Anything, userID, req).Return(expectedResponse, nil)

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPut, "/users/"+userID.String(), bytes.NewReader(body))
		recorder := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/users/{id}", handler.UpdateUser)
		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		var response dto.UserResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, expectedResponse.ID, response.ID)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid User ID", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		req := &dto.UpdateUserRequest{
			FirstName: "Jane",
		}

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPut, "/users/invalid-id", bytes.NewReader(body))
		recorder := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/users/{id}", handler.UpdateUser)
		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		userID := uuid.New()
		request := httptest.NewRequest(http.MethodPut, "/users/"+userID.String(), bytes.NewReader([]byte("invalid json")))
		recorder := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/users/{id}", handler.UpdateUser)
		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		var response dto.ErrorResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, "Invalid request body", response.Error)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		userID := uuid.New()
		req := &dto.UpdateUserRequest{
			FirstName: "Jane",
		}

		mockService.On("UpdateUser", mock.Anything, userID, req).Return(nil, errors.New("update failed"))

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPut, "/users/"+userID.String(), bytes.NewReader(body))
		recorder := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/users/{id}", handler.UpdateUser)
		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		mockService.AssertExpectations(t)
	})
}

func TestDeleteUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		userID := uuid.New()
		mockService.On("DeleteUser", mock.Anything, userID).Return(nil)

		request := httptest.NewRequest(http.MethodDelete, "/users/"+userID.String(), nil)
		recorder := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/users/{id}", handler.DeleteUser)
		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		var response map[string]string
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, "User deleted successfully", response["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid User ID", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		request := httptest.NewRequest(http.MethodDelete, "/users/invalid-id", nil)
		recorder := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/users/{id}", handler.DeleteUser)
		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		userID := uuid.New()
		mockService.On("DeleteUser", mock.Anything, userID).Return(errors.New("user not found"))

		request := httptest.NewRequest(http.MethodDelete, "/users/"+userID.String(), nil)
		recorder := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/users/{id}", handler.DeleteUser)
		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
		mockService.AssertExpectations(t)
	})
}

func TestListUsers(t *testing.T) {
	t.Run("Success With Pagination", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		expectedResponse := &dto.ListUsersResponse{
			Users: []dto.UserResponse{
				{
					ID:        uuid.New(),
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@example.com",
				},
			},
			Total:    1,
			Page:     1,
			PageSize: 10,
		}

		mockService.On("ListUsers", mock.Anything, 1, 10).Return(expectedResponse, nil)

		request := httptest.NewRequest(http.MethodGet, "/users?page=1&page_size=10", nil)
		recorder := httptest.NewRecorder()

		handler.ListUsers(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		var response dto.ListUsersResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, expectedResponse.Total, response.Total)
		mockService.AssertExpectations(t)
	})

	t.Run("Success With Default Pagination", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		expectedResponse := &dto.ListUsersResponse{
			Users:    []dto.UserResponse{},
			Total:    0,
			Page:     1,
			PageSize: 10,
		}

		mockService.On("ListUsers", mock.Anything, 1, 10).Return(expectedResponse, nil)

		request := httptest.NewRequest(http.MethodGet, "/users", nil)
		recorder := httptest.NewRecorder()

		handler.ListUsers(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		mockService.On("ListUsers", mock.Anything, 1, 10).Return(nil, errors.New("database error"))

		request := httptest.NewRequest(http.MethodGet, "/users", nil)
		recorder := httptest.NewRecorder()

		handler.ListUsers(recorder, request)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		var response dto.ErrorResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, "Failed to list users", response.Error)
		mockService.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		req := &dto.LoginRequest{
			Email:    "john.doe@example.com",
			Password: "password123",
		}

		expectedResponse := &dto.LoginResponse{
			User: dto.UserResponse{
				ID:        uuid.New(),
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
			},
			Message: "Login successful",
		}

		mockService.On("AuthenticateUser", mock.Anything, req).Return(expectedResponse, nil)

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		recorder := httptest.NewRecorder()

		handler.Login(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		var response dto.LoginResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, expectedResponse.Message, response.Message)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		request := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte("invalid json")))
		recorder := httptest.NewRecorder()

		handler.Login(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		var response dto.ErrorResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, "Invalid request body", response.Error)
	})

	t.Run("Authentication Failed", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		req := &dto.LoginRequest{
			Email:    "john.doe@example.com",
			Password: "wrongpassword",
		}

		mockService.On("AuthenticateUser", mock.Anything, req).Return(nil, errors.New("invalid credentials"))

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		recorder := httptest.NewRecorder()

		handler.Login(recorder, request)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		var response dto.ErrorResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, "Authentication failed", response.Error)
		mockService.AssertExpectations(t)
	})
}

func TestAssignAccessLevels(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		userID := uuid.New()
		req := &dto.AssignAccessLevelRequest{
			AccessLevelIDs: []int{1, 2, 3},
		}

		mockService.On("AssignAccessLevels", mock.Anything, userID, req).Return(nil)

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/users/"+userID.String()+"/access-levels", bytes.NewReader(body))
		recorder := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/users/{id}/access-levels", handler.AssignAccessLevels)
		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		var response map[string]string
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, "Access levels assigned successfully", response["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid User ID", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		req := &dto.AssignAccessLevelRequest{
			AccessLevelIDs: []int{1, 2, 3},
		}

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/users/invalid-id/access-levels", bytes.NewReader(body))
		recorder := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/users/{id}/access-levels", handler.AssignAccessLevels)
		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		userID := uuid.New()
		request := httptest.NewRequest(http.MethodPost, "/users/"+userID.String()+"/access-levels", bytes.NewReader([]byte("invalid json")))
		recorder := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/users/{id}/access-levels", handler.AssignAccessLevels)
		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		userID := uuid.New()
		req := &dto.AssignAccessLevelRequest{
			AccessLevelIDs: []int{1, 2, 3},
		}

		mockService.On("AssignAccessLevels", mock.Anything, userID, req).Return(errors.New("assignment failed"))

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/users/"+userID.String()+"/access-levels", bytes.NewReader(body))
		recorder := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/users/{id}/access-levels", handler.AssignAccessLevels)
		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetUserAccessLevels(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		userID := uuid.New()
		expectedResponse := []dto.AccessLevelResponse{
			{ID: 1, Name: "Admin", Description: "Administrator access"},
			{ID: 2, Name: "User", Description: "User access"},
		}

		mockService.On("GetUserAccessLevels", mock.Anything, userID).Return(expectedResponse, nil)

		request := httptest.NewRequest(http.MethodGet, "/users/"+userID.String()+"/access-levels", nil)
		recorder := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/users/{id}/access-levels", handler.GetUserAccessLevels)
		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		var response []dto.AccessLevelResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, len(expectedResponse), len(response))
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid User ID", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		request := httptest.NewRequest(http.MethodGet, "/users/invalid-id/access-levels", nil)
		recorder := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/users/{id}/access-levels", handler.GetUserAccessLevels)
		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		userID := uuid.New()
		mockService.On("GetUserAccessLevels", mock.Anything, userID).Return(nil, errors.New("user not found"))

		request := httptest.NewRequest(http.MethodGet, "/users/"+userID.String()+"/access-levels", nil)
		recorder := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/users/{id}/access-levels", handler.GetUserAccessLevels)
		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
		mockService.AssertExpectations(t)
	})
}

func TestRespondWithJSON(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		payload := map[string]string{"message": "success"}

		respondWithJSON(recorder, http.StatusOK, payload)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

		var response map[string]string
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, "success", response["message"])
	})

	t.Run("Marshal Error", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		// Create a channel which cannot be marshaled to JSON
		payload := make(chan int)

		respondWithJSON(recorder, http.StatusOK, payload)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}

func TestRespondWithError(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		respondWithError(recorder, http.StatusBadRequest, "Test error", "Error message")

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

		var response dto.ErrorResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, "Test error", response.Error)
		assert.Equal(t, "Error message", response.Message)
	})
}
