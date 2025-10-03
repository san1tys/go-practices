package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type UserHandler struct{}

type getUserResponse struct {
	UserID int `json:"user_id"`
}

type createUserRequest struct {
	Name string `json:"name"`
}

type createUserResponse struct {
	Created string `json:"created"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "invalid id"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "invalid id"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(getUserResponse{UserID: id})
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "invalid name"})
		return
	}

	if req.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "invalid name"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createUserResponse{Created: req.Name})
}
