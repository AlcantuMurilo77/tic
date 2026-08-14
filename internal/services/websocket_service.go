package services

import (
	"encoding/json"
	"net/http"

	"github.com/AlcantuMurilo77/tic/internal/models"
	"github.com/gorilla/websocket"
)

type WebSocketService struct {
	upgrader websocket.Upgrader
}

func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (s *WebSocketService) Connect(
	w http.ResponseWriter,
	r *http.Request,
) (*websocket.Conn, error) {
	return s.upgrader.Upgrade(w, r, nil)
}

func (s *WebSocketService) ReadMove(
	conn *websocket.Conn,
) (*models.MoveMessage, error) {

	var move models.MoveMessage

	if err := conn.ReadJSON(&move); err != nil {
		return nil, err
	}

	return &move, nil
}

func (s *WebSocketService) Send(
	conn *websocket.Conn,
	payload any,
) error {
	return conn.WriteJSON(payload)
}

func (s *WebSocketService) SendRaw(
	conn *websocket.Conn,
	payload []byte,
) error {
	var data any

	if err := json.Unmarshal(payload, &data); err != nil {
		return err
	}
	return conn.WriteJSON(data)
}
