package database

import (
	"context"
)

func (w *Wrapper) ReplaceOrganization(ctx context.Context, org Organization) error {
	_, err := w.db.ExecContext(ctx, `UPDATE "organization" SET "name" = $1 WHERE "id" = $2`, org.Name, org.ID)

	return err
}

func (w *Wrapper) ReplaceDevice(ctx context.Context, dev Device) error {
	_, err := w.db.ExecContext(ctx, `UPDATE "device" SET "name" = $1 WHERE "id" = $2`, dev.Name, dev.ID)

	return err
}