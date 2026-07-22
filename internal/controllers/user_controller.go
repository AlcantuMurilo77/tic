package controllers

import (
	"encoding/json"
	"github.com/AlcantuMurilo77/tic/internal/services"
	"net/http"
)

type UserController struct {
	userService *services.UserService
}

type CreateUserRequest struct {
	Name    string `json:"name"`
	Country string `json:"country"`
	Xman    bool   `json:"xman"`
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{
		userService: service,
	}
}

func (c *UserController) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := c.userService.Create(
		r.Context(),
		request.Name,
		request.Country,
		request.Xman,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(user)
}
