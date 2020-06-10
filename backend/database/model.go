package database

import (
	"context"

	"github.com/google/uuid"
)

type User struct {
	ID   uuid.UUID
	Name string
	Age  int
}

func (w *Wrapper) ReadUser(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User

	err := w.db.GetContext(ctx, &user, `SELECT * FROM "user" WHERE "id" = $1`, id)

	return &user, err
}