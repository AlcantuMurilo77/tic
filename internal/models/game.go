package models

import (
	"github.com/google/uuid"
	"time"
)

type Game struct {
	GameUuid  uuid.UUID `bson:"game_uuid"`
	UserX     uuid.UUID `bson:"user_x"`
	UserO     uuid.UUID `bson:"user_o"`
	WinnerID  uuid.UUID `bson:"winner_id"`
	StartedAt time.Time `bson:"started_at"`
	EndedAt   time.Time `bson:"ended_at"`
	Board     [][]int   `bson:"board"`
}
