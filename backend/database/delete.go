package database

import (
	"context"

	"github.com/google/uuid"
)

func (w *Wrapper) DeleteOrganization(ctx context.Context, id uuid.UUID) error {
	_, err := w.db.ExecContext(ctx, `DELETE FROM "organization" WHERE "id" = $1`, id)
	return err
}

func (w *Wrapper) DeleteStudent(ctx context.Context, id uuid.UUID) error {
	_, err := w.db.ExecContext(ctx, `DELETE FROM "student" WHERE "id" = $1`, id)
	return err
}

func (w *Wrapper) DeleteDevice(ctx context.Context, id uuid.UUID) error {
	_, err := w.db.ExecContext(ctx, `DELETE FROM "device" WHERE "id" = $1`, id)
	return err
}