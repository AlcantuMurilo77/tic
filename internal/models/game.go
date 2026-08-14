package models

import (
	"github.com/google/uuid"
	"time"
)

type Game struct {
	GameUuid  uuid.UUID `json:"id" bson:"_id"`
	UserX     uuid.UUID `json:"user_x" bson:"user_x"`
	UserO     uuid.UUID `json:"user_o" bson:"user_o"`
	WinnerID  uuid.UUID `json:"winner_id" bson:"winner_id"`
	StartedAt time.Time `json:"started_at" bson:"started_at"`
	EndedAt   time.Time `json:"ended_at" bson:"ended_at"`
	Board     [][]int   `json:"board" bson:"board"`
}
