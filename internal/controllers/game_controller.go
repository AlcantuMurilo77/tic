package controllers

import (
	"encoding/json"
	"errors"
	"github.com/AlcantuMurilo77/tic/internal/services"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"log"
	"net/http"
)

type GameController struct {
	gameService *services.GameService
}

type CreateGameRequest struct {
	UserX uuid.UUID `json:"user_x"`
	UserO uuid.UUID `json:"user_o"`
}

func NewGameController(service *services.GameService) *GameController {
	return &GameController{
		gameService: service,
	}
}

func (c *GameController) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request CreateGameRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	game, err := c.gameService.Create(
		r.Context(),
		request.UserX,
		request.UserO,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(game)
}

func (c *GameController) FindAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	games, err := c.gameService.FindAll(r.Context())
	if err != nil {
		http.Error(
			w,
			"failed to retrieve games",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(games); err != nil {
		return
	}
}

func (c *GameController) FindOne(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("game_id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	game, err := c.gameService.FindByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "game not found", http.StatusNotFound)
			return
		}
		log.Printf("error finding game by id %q: %v", id, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)
}
