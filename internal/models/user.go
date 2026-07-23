package models

import (
	"github.com/google/uuid"
)

type User struct {
	UserUuid uuid.UUID `bson:"user_uuid"`
	Name     string    `bson:"name"`
	Country  string    `bson:"country"` //no usa users allowed
	Xman     bool      `bson:"xman"`
}
