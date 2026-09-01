package models

import "github.com/google/uuid"

type MoveMessage struct {
	GameID   uuid.UUID `json:"game_id"`
	PlayerID uuid.UUID `json:"player_id"`
	Row      int       `json:"row"`
	Col      int       `json:"col"`
}

type WebSocketMessage struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}
