package api

import (
	"github.com/google/uuid"
	"gitlab.com/timeterm/timeterm/backend/database"
)

type User struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Age  int       `json:"age"`
}

func UserFrom(user *database.User) *User {
	return (*User)(user)
}
