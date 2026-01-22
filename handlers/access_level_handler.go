package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/wabtcdi/user_service/dto"
	"github.com/wabtcdi/user_service/service"
)

type AccessLevelHandler struct {
	service service.AccessLevelServiceInterface
}

func NewAccessLevelHandler(service service.AccessLevelServiceInterface) *AccessLevelHandler {
	return &AccessLevelHandler{service: service}
}

func (h *AccessLevelHandler) CreateAccessLevel(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAccessLevelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Errorf("Failed to decode request: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	accessLevel, err := h.service.CreateAccessLevel(r.Context(), &req)
	if err != nil {
		logrus.Errorf("Failed to create access level: %v", err)
		respondWithError(w, http.StatusBadRequest, "Failed to create access level", err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, accessLevel)
}

func (h *AccessLevelHandler) GetAccessLevel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid access level ID", err.Error())
		return
	}

	accessLevel, err := h.service.GetAccessLevel(r.Context(), id)
	if err != nil {
		logrus.Errorf("Failed to get access level: %v", err)
		respondWithError(w, http.StatusNotFound, "Access level not found", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, accessLevel)
}

func (h *AccessLevelHandler) ListAccessLevels(w http.ResponseWriter, r *http.Request) {
	accessLevels, err := h.service.ListAccessLevels(r.Context())
	if err != nil {
		logrus.Errorf("Failed to list access levels: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to list access levels", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, accessLevels)
}
