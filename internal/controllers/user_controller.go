package controllers

import (
	"encoding/json"
	"errors"
	"github.com/AlcantuMurilo77/tic/internal/services"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"log/slog"
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
		slog.Error("failed to create user", "name", request.Name, "error", err)
		http.Error(w, "failed to create user", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(user)
}

func (c *UserController) FindAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := c.userService.FindAll(r.Context())
	if err != nil {
		slog.Error("failed to retrieve users", "error", err)
		http.Error(
			w,
			"failed to retrieve users",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(users); err != nil {
		return
	}
}

func (c *UserController) FindOne(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("user_id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	game, err := c.userService.FindByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "game not found", http.StatusNotFound)
			return
		}
		slog.Error("failed to find user", "user_id", id, "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)
}
