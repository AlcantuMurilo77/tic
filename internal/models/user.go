package models

import (
	"github.com/google/uuid"
)

type User struct {
	UserUuid uuid.UUID `json:"user_id" bson:"_id"`
	Name     string    `json:"name" bson:"name"`
	Country  string    `json:"country" bson:"country"` //no usa users allowed
	Xman     bool      `json:"x_man" bson:"x_man"`
}
