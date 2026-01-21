package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/wabtcdi/user_service/dto"
	"github.com/wabtcdi/user_service/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Errorf("Failed to decode request: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	user, err := h.userService.CreateUser(r.Context(), &req)
	if err != nil {
		logrus.Errorf("Failed to create user: %v", err)
		respondWithError(w, http.StatusBadRequest, "Failed to create user", err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	user, err := h.userService.GetUser(r.Context(), id)
	if err != nil {
		logrus.Errorf("Failed to get user: %v", err)
		respondWithError(w, http.StatusNotFound, "User not found", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Errorf("Failed to decode request: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	user, err := h.userService.UpdateUser(r.Context(), id, &req)
	if err != nil {
		logrus.Errorf("Failed to update user: %v", err)
		respondWithError(w, http.StatusBadRequest, "Failed to update user", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	err = h.userService.DeleteUser(r.Context(), id)
	if err != nil {
		logrus.Errorf("Failed to delete user: %v", err)
		respondWithError(w, http.StatusNotFound, "User not found", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	response, err := h.userService.ListUsers(r.Context(), page, pageSize)
	if err != nil {
		logrus.Errorf("Failed to list users: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to list users", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Errorf("Failed to decode request: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	response, err := h.userService.AuthenticateUser(r.Context(), &req)
	if err != nil {
		logrus.Errorf("Authentication failed: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Authentication failed", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *UserHandler) AssignAccessLevels(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	var req dto.AssignAccessLevelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Errorf("Failed to decode request: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	err = h.userService.AssignAccessLevels(r.Context(), id, &req)
	if err != nil {
		logrus.Errorf("Failed to assign access levels: %v", err)
		respondWithError(w, http.StatusBadRequest, "Failed to assign access levels", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Access levels assigned successfully"})
}

func (h *UserHandler) GetUserAccessLevels(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	accessLevels, err := h.userService.GetUserAccessLevels(r.Context(), id)
	if err != nil {
		logrus.Errorf("Failed to get user access levels: %v", err)
		respondWithError(w, http.StatusNotFound, "Failed to get access levels", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, accessLevels)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		logrus.Errorf("Failed to marshal response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, error string, message string) {
	respondWithJSON(w, code, dto.ErrorResponse{
		Error:   error,
		Message: message,
	})
}
