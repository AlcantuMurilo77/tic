package controllers

import (
	"github.com/AlcantuMurilo77/tic/internal/services"
	"log"
	"net/http"
)

type WebSocketController struct {
	webSocketService *services.WebSocketService
	gameService      *services.GameService
}

func NewWebSocketController(
	webSocketService *services.WebSocketService,
	gameService *services.GameService,
) *WebSocketController {
	return &WebSocketController{
		webSocketService: webSocketService,
		gameService:      gameService,
	}
}

func (c *WebSocketController) Connect(
	w http.ResponseWriter,
	r *http.Request,
) {
	conn, err := c.webSocketService.Connect(w, r)
	if err != nil {
		log.Printf("failed to upgrade websocket: %v", err)
		return
	}
	defer conn.Close()

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

		if err := c.webSocketService.Send(conn, game); err != nil {
			log.Printf("failed to send game state: %v", err)
			return
		}
	}
}
