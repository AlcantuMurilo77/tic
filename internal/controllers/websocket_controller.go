package controllers

import (
	"github.com/AlcantuMurilo77/tic/internal/services"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type WebSocketController struct {
	webSocketService *services.WebSocketService
	gameService      *services.GameService
	hub              *services.WebsocketHub
}

func NewWebSocketController(
	webSocketService *services.WebSocketService,
	gameService *services.GameService,
	hub *services.WebsocketHub,
) *WebSocketController {
	return &WebSocketController{
		webSocketService: webSocketService,
		gameService:      gameService,
		hub:              hub,
	}
}

func (c *WebSocketController) Connect(
	w http.ResponseWriter,
	r *http.Request,
) {

	gameIDStr := r.URL.Query().Get("game_id")
	playerIDStr := r.URL.Query().Get("player_id")

	gameID, err := uuid.Parse(gameIDStr)
	if err != nil {
		http.Error(w, "invalid game id", http.StatusBadRequest)
		return
	}

	playerID, err := uuid.Parse(playerIDStr)
	if err != nil {
		http.Error(w, "invalid player id", http.StatusBadRequest)
		return
	}

	gameModel, err := c.gameService.FindByID(
		r.Context(),
		gameID,
	)
	if err != nil {
		http.Error(w, "game not found", http.StatusNotFound)
		return
	}

	if playerID != gameModel.UserX && playerID != gameModel.UserO {
		http.Error(w, "player dfoes not belong to this game", http.StatusForbidden)
		return
	}

	conn, err := c.webSocketService.Connect(w, r)
	if err != nil {
		log.Printf("failed to upgrade websocket: %v", err)
		return
	}
	defer conn.Close()

	c.hub.AddPlayer(
		gameID,
		playerID,
		gameModel.UserX,
		gameModel.UserO,
		conn,
	)
	defer c.hub.RemovePlayer(
		gameID,
		playerID,
		gameModel.UserX,
		gameModel.UserO,
	)

	for {
		move, err := c.webSocketService.ReadMove(conn)
		if err != nil {
			log.Printf("failed to read websocket message: %v", err)
			return
		}

		game, err := c.gameService.Move(
			r.Context(),
			move.GameID,
			move.PlayerID,
			move.Row,
			move.Col,
		)

		if err != nil {
			log.Printf("failed to process move: %v", err)

			c.webSocketService.Send(conn, map[string]string{
				"error": err.Error(),
			})

			continue
		}

		room := c.hub.GetRoom(move.GameID)

		if room == nil {
			log.Printf("room not found for game %s", move.GameID)
			continue
		}

		if room.PlayerX != nil {
			if err := c.webSocketService.Send(room.PlayerX, game); err != nil {
				log.Printf("failed to send game state to player X: %v", err)
			}
		}

		if room.PlayerO != nil {
			if err := c.webSocketService.Send(room.PlayerO, game); err != nil {
				log.Printf("failed to send game state to player O: %v", err)
			}
		}

		if err := c.webSocketService.Send(conn, game); err != nil {
			log.Printf("failed to send game state: %v", err)
			return
		}
	}
}
