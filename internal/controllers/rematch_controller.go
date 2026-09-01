package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AlcantuMurilo77/tic/internal/models"
	"github.com/AlcantuMurilo77/tic/internal/services"
	"github.com/google/uuid"
)

type RematchPlayerRequest struct {
	PlayerID uuid.UUID `json:"player_id"`
}

type RematchController struct {
	rematchService *services.RematchService
	hub            *services.WebsocketHub
}

func NewRematchController(
	service *services.RematchService,
	hub *services.WebsocketHub,
) *RematchController {
	return &RematchController{
		rematchService: service,
		hub:            hub,
	}
}

func (c *RematchController) Request(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	gameID, err := uuid.Parse(
		r.URL.Query().Get("game_id"),
	)

	if err != nil {
		http.Error(
			w,
			"invalid game id",
			http.StatusBadRequest,
		)
		return
	}

	var body RematchPlayerRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	request, err := c.rematchService.Request(
		r.Context(),
		gameID,
		body.PlayerID,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	c.broadcast(gameID, "rematch_requested", request)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(request)

}

func (c *RematchController) Accept(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	rematchID, err := uuid.Parse(
		r.URL.Query().Get("rematch_id"),
	)

	if err != nil {
		http.Error(
			w,
			"invalid rematch id",
			http.StatusBadRequest,
		)
		return
	}

	var body RematchPlayerRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	newGame, err := c.rematchService.Accept(
		r.Context(),
		rematchID,
		body.PlayerID,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	if newGame.RematchOfGameID != nil {
		c.broadcast(*newGame.RematchOfGameID, "rematch_accepted", newGame)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newGame)

}

func (c *RematchController) broadcast(gameID uuid.UUID, messageType string, payload any) {
	errs := c.hub.Broadcast(gameID, models.WebSocketMessage{
		Type:    messageType,
		Payload: payload,
	})
	for _, err := range errs {
		log.Printf("failed to broadcast %s: %v", messageType, err)
	}
}
