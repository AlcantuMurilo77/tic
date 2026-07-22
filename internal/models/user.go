package models

import (
	"github.com/google/uuid"
)

type User struct {
	UserUuid uuid.UUID
	Name     string
	Country  string //no usa users allowed
	Xman     bool
}
