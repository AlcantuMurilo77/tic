package models

import (
	"github.com/google/uuid"
)

type User struct {
	UserUuid uuid.UUID `json:"id" bson:"_id"`
	Name     string    `json:"name" bson:"name"`
	Country  string    `json:"country" bson:"country"`
	Xman     bool      `json:"xman" bson:"xman"`
}
