package models

import (
	"github.com/google/uuid"
	"time"
)

type RematchStatus string

const (
	RematchPending  RematchStatus = "pending"
	RematchAccepted RematchStatus = "accepted"
)

type RematchRequest struct {
	RematchUuid         uuid.UUID     `json:"id" bson:"_id"`
	OriginalGameId      uuid.UUID     `json:"original_game_id" bson:"original_game_id"`
	RequestedByPlayerId uuid.UUID     `json:"requested_by_player_id" bson:"requested_by_player_id"`
	AcceptedByPlayerID  *uuid.UUID    `json:"accepted_by_player_id,omitempty" bson:"accepted_by_player_id,omitempty"`
	Status              RematchStatus `json:"status" bson:"status"`
	NewGameId           *uuid.UUID    `json:"new_game_id,omitempty" bson:"new_game_id,omitempty"`
	CreatedAt           time.Time     `json:"created_at" bson:"created_at"`
	AcceptedAt          *time.Time    `json:"accepted_at,omitempty" bson:"accepted_at,omitempty"`
}
