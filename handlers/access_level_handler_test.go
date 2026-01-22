package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wabtcdi/user_service/dto"
	"github.com/wabtcdi/user_service/service"
)

// MockAccessLevelService is a mock implementation of the AccessLevelServiceInterface
type MockAccessLevelService struct {
	mock.Mock
}

// Ensure MockAccessLevelService implements service.AccessLevelServiceInterface
var _ service.AccessLevelServiceInterface = (*MockAccessLevelService)(nil)

func (m *MockAccessLevelService) CreateAccessLevel(ctx context.Context, req *dto.CreateAccessLevelRequest) (*dto.AccessLevelResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AccessLevelResponse), args.Error(1)
}

func (m *MockAccessLevelService) GetAccessLevel(ctx context.Context, id int) (*dto.AccessLevelResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AccessLevelResponse), args.Error(1)
}

func (m *MockAccessLevelService) ListAccessLevels(ctx context.Context) ([]dto.AccessLevelResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.AccessLevelResponse), args.Error(1)
}

func TestCreateAccessLevel(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockAccessLevelService)
		handler := NewAccessLevelHandler(mockService)

		req := &dto.CreateAccessLevelRequest{
			Name:        "Admin",
			Description: "Administrator access level",
		}

		expectedResponse := &dto.AccessLevelResponse{
			ID:          1,
			Name:        "Admin",
			Description: "Administrator access level",
		}

		mockService.On("CreateAccessLevel", mock.Anything, req).Return(expectedResponse, nil)

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/access-levels", bytes.NewReader(body))
		recorder := httptest.NewRecorder()

		handler.CreateAccessLevel(recorder, request)

		assert.Equal(t, http.StatusCreated, recorder.Code)
		var response dto.AccessLevelResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, expectedResponse.ID, response.ID)
		assert.Equal(t, expectedResponse.Name, response.Name)
		assert.Equal(t, expectedResponse.Description, response.Description)
		mockService.AssertExpectations(t)
	})

	t.Run("Success Without Description", func(t *testing.T) {
		mockService := new(MockAccessLevelService)
		handler := NewAccessLevelHandler(mockService)

		req := &dto.CreateAccessLevelRequest{
			Name: "User",
		}

		expectedResponse := &dto.AccessLevelResponse{
			ID:   2,
			Name: "User",
		}

		mockService.On("CreateAccessLevel", mock.Anything, req).Return(expectedResponse, nil)

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/access-levels", bytes.NewReader(body))
		recorder := httptest.NewRecorder()

		handler.CreateAccessLevel(recorder, request)

		assert.Equal(t, http.StatusCreated, recorder.Code)
		var response dto.AccessLevelResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, expectedResponse.ID, response.ID)
		assert.Equal(t, expectedResponse.Name, response.Name)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		mockService := new(MockAccessLevelService)
		handler := NewAccessLevelHandler(mockService)

		request := httptest.NewRequest(http.MethodPost, "/access-levels", bytes.NewReader([]byte("invalid json")))
		recorder := httptest.NewRecorder()

		handler.CreateAccessLevel(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		var response dto.ErrorResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, "Invalid request body", response.Error)
	})

	t.Run("Service Error - Duplicate Name", func(t *testing.T) {
		mockService := new(MockAccessLevelService)
		handler := NewAccessLevelHandler(mockService)

		req := &dto.CreateAccessLevelRequest{
			Name:        "Admin",
			Description: "Administrator access level",
		}

		mockService.On("CreateAccessLevel", mock.Anything, req).Return(nil, errors.New("access level with name Admin already exists"))

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/access-levels", bytes.NewReader(body))
		recorder := httptest.NewRecorder()

		handler.CreateAccessLevel(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		var response dto.ErrorResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, "Failed to create access level", response.Error)
		assert.Contains(t, response.Message, "already exists")
		mockService.AssertExpectations(t)
	})

	t.Run("Service Error - Generic", func(t *testing.T) {
		mockService := new(MockAccessLevelService)
		handler := NewAccessLevelHandler(mockService)

		req := &dto.CreateAccessLevelRequest{
			Name: "TestLevel",
		}

		mockService.On("CreateAccessLevel", mock.Anything, req).Return(nil, errors.New("database error"))

		body, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/access-levels", bytes.NewReader(body))
		recorder := httptest.NewRecorder()

		handler.CreateAccessLevel(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		var response dto.ErrorResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, "Failed to create access level", response.Error)
		mockService.AssertExpectations(t)
	})
}

func TestGetAccessLevel(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockAccessLevelService)
		handler := NewAccessLevelHandler(mockService)

		accessLevelID := 1
		expectedResponse := &dto.AccessLevelResponse{
			ID:          accessLevelID,
			Name:        "Admin",
			Description: "Administrator access level",
		}

		mockService.On("GetAccessLevel", mock.Anything, accessLevelID).Return(expectedResponse, nil)

		request := httptest.NewRequest(http.MethodGet, "/access-levels/1", nil)
		recorder := httptest.NewRecorder()

		// Setup mux to parse path variables
		router := mux.NewRouter()
		router.HandleFunc("/access-levels/{id}", handler.GetAccessLevel)
		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		var response dto.AccessLevelResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, expectedResponse.ID, response.ID)
		assert.Equal(t, expectedResponse.Name, response.Name)
		assert.Equal(t, expectedResponse.Description, response.Description)
		mockService.AssertExpectations(t)
	})

	t.Run("Success With Empty Description", func(t *testing.T) {
		mockService := new(MockAccessLevelService)
		handler := NewAccessLevelHandler(mockService)

		accessLevelID := 2
		expectedResponse := &dto.AccessLevelResponse{
			ID:   accessLevelID,
			Name: "User",
		}

		mockService.On("GetAccessLevel", mock.Anything, accessLevelID).Return(expectedResponse, nil)

		request := httptest.NewRequest(http.MethodGet, "/access-levels/2", nil)
		recorder := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/access-levels/{id}", handler.GetAccessLevel)
		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		var response dto.AccessLevelResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, expectedResponse.ID, response.ID)
		assert.Equal(t, expectedResponse.Name, response.Name)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Access Level ID - Not a Number", func(t *testing.T) {
		mockService := new(MockAccessLevelService)
		handler := NewAccessLevelHandler(mockService)

		request := httptest.NewRequest(http.MethodGet, "/access-levels/invalid", nil)
		recorder := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/access-levels/{id}", handler.GetAccessLevel)
		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		var response dto.ErrorResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, "Invalid access level ID", response.Error)
	})

	t.Run("Invalid Access Level ID - Negative Number", func(t *testing.T) {
		mockService := new(MockAccessLevelService)
		handler := NewAccessLevelHandler(mockService)

		// Note: strconv.Atoi will parse -1 successfully, but the service should handle validation
		accessLevelID := -1
		mockService.On("GetAccessLevel", mock.Anything, accessLevelID).Return(nil, errors.New("invalid access level ID"))

		request := httptest.NewRequest(http.MethodGet, "/access-levels/-1", nil)
		recorder := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/access-levels/{id}", handler.GetAccessLevel)
		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Access Level Not Found", func(t *testing.T) {
		mockService := new(MockAccessLevelService)
		handler := NewAccessLevelHandler(mockService)

		accessLevelID := 999
		mockService.On("GetAccessLevel", mock.Anything, accessLevelID).Return(nil, errors.New("access level not found"))

		request := httptest.NewRequest(http.MethodGet, "/access-levels/999", nil)
		recorder := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/access-levels/{id}", handler.GetAccessLevel)
		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
		var response dto.ErrorResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, "Access level not found", response.Error)
		mockService.AssertExpectations(t)
	})
}

func TestListAccessLevels(t *testing.T) {
	t.Run("Success With Multiple Access Levels", func(t *testing.T) {
		mockService := new(MockAccessLevelService)
		handler := NewAccessLevelHandler(mockService)

		expectedResponse := []dto.AccessLevelResponse{
			{
				ID:          1,
				Name:        "Admin",
				Description: "Administrator access level",
			},
			{
				ID:          2,
				Name:        "User",
				Description: "Standard user access level",
			},
			{
				ID:   3,
				Name: "Guest",
			},
		}

		mockService.On("ListAccessLevels", mock.Anything).Return(expectedResponse, nil)

		request := httptest.NewRequest(http.MethodGet, "/access-levels", nil)
		recorder := httptest.NewRecorder()

		handler.ListAccessLevels(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		var response []dto.AccessLevelResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, len(expectedResponse), len(response))
		assert.Equal(t, expectedResponse[0].ID, response[0].ID)
		assert.Equal(t, expectedResponse[0].Name, response[0].Name)
		assert.Equal(t, expectedResponse[1].ID, response[1].ID)
		assert.Equal(t, expectedResponse[2].Name, response[2].Name)
		mockService.AssertExpectations(t)
	})

	t.Run("Success With Empty List", func(t *testing.T) {
		mockService := new(MockAccessLevelService)
		handler := NewAccessLevelHandler(mockService)

		expectedResponse := []dto.AccessLevelResponse{}

		mockService.On("ListAccessLevels", mock.Anything).Return(expectedResponse, nil)

		request := httptest.NewRequest(http.MethodGet, "/access-levels", nil)
		recorder := httptest.NewRecorder()

		handler.ListAccessLevels(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		var response []dto.AccessLevelResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, 0, len(response))
		mockService.AssertExpectations(t)
	})

	t.Run("Success With Single Access Level", func(t *testing.T) {
		mockService := new(MockAccessLevelService)
		handler := NewAccessLevelHandler(mockService)

		expectedResponse := []dto.AccessLevelResponse{
			{
				ID:          1,
				Name:        "Admin",
				Description: "Administrator access level",
			},
		}

		mockService.On("ListAccessLevels", mock.Anything).Return(expectedResponse, nil)

		request := httptest.NewRequest(http.MethodGet, "/access-levels", nil)
		recorder := httptest.NewRecorder()

		handler.ListAccessLevels(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		var response []dto.AccessLevelResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, 1, len(response))
		assert.Equal(t, expectedResponse[0].ID, response[0].ID)
		mockService.AssertExpectations(t)
	})

	t.Run("Service Error - Database Error", func(t *testing.T) {
		mockService := new(MockAccessLevelService)
		handler := NewAccessLevelHandler(mockService)

		mockService.On("ListAccessLevels", mock.Anything).Return(nil, errors.New("database connection error"))

		request := httptest.NewRequest(http.MethodGet, "/access-levels", nil)
		recorder := httptest.NewRecorder()

		handler.ListAccessLevels(recorder, request)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		var response dto.ErrorResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, "Failed to list access levels", response.Error)
		assert.Contains(t, response.Message, "database")
		mockService.AssertExpectations(t)
	})

	t.Run("Service Error - Generic Error", func(t *testing.T) {
		mockService := new(MockAccessLevelService)
		handler := NewAccessLevelHandler(mockService)

		mockService.On("ListAccessLevels", mock.Anything).Return(nil, errors.New("unexpected error"))

		request := httptest.NewRequest(http.MethodGet, "/access-levels", nil)
		recorder := httptest.NewRecorder()

		handler.ListAccessLevels(recorder, request)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		var response dto.ErrorResponse
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, "Failed to list access levels", response.Error)
		mockService.AssertExpectations(t)
	})
}

func TestAccessLevelHandlerIntegration(t *testing.T) {
	t.Run("Create Then Get Access Level", func(t *testing.T) {
		mockService := new(MockAccessLevelService)
		handler := NewAccessLevelHandler(mockService)

		// Create
		createReq := &dto.CreateAccessLevelRequest{
			Name:        "Manager",
			Description: "Manager access level",
		}

		createResponse := &dto.AccessLevelResponse{
			ID:          5,
			Name:        "Manager",
			Description: "Manager access level",
		}

		mockService.On("CreateAccessLevel", mock.Anything, createReq).Return(createResponse, nil)

		body, _ := json.Marshal(createReq)
		request := httptest.NewRequest(http.MethodPost, "/access-levels", bytes.NewReader(body))
		recorder := httptest.NewRecorder()

		handler.CreateAccessLevel(recorder, request)
		assert.Equal(t, http.StatusCreated, recorder.Code)

		// Get
		mockService.On("GetAccessLevel", mock.Anything, 5).Return(createResponse, nil)

		request2 := httptest.NewRequest(http.MethodGet, "/access-levels/5", nil)
		recorder2 := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/access-levels/{id}", handler.GetAccessLevel)
		router.ServeHTTP(recorder2, request2)

		assert.Equal(t, http.StatusOK, recorder2.Code)
		var response dto.AccessLevelResponse
		json.Unmarshal(recorder2.Body.Bytes(), &response)
		assert.Equal(t, createResponse.ID, response.ID)
		assert.Equal(t, createResponse.Name, response.Name)
		mockService.AssertExpectations(t)
	})
}

func TestNewAccessLevelHandler(t *testing.T) {
	t.Run("Creates Handler Successfully", func(t *testing.T) {
		mockService := new(MockAccessLevelService)
		handler := NewAccessLevelHandler(mockService)

		assert.NotNil(t, handler)
		assert.Equal(t, mockService, handler.service)
	})
}
