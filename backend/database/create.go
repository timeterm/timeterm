package database

import (
	"context"

	"github.com/google/uuid"
)

type Organization struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

func (w *Wrapper) CreateOrganization(ctx context.Context, name string) (Organization, error) {
	org := Organization{
		Name: name,
	}

	row := w.db.QueryRowContext(ctx, `INSERT INTO "organization" ("id", "name") VALUES (DEFAULT, $1) RETURNING "id"`, name)

	return org, row.Scan(&org.ID)
}
