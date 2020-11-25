package database

import (
	"context"
)

func (w *Wrapper) SetJWTMigrationVersion(ctx context.Context, to int) error {
	tx, err := w.txbc.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err = tx.ExecContext(ctx, `TRUNCATE jwt_migration`); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO jwt_migration (version) VALUES ($1)`, to); err != nil {
		return err
	}

	return tx.Commit()
}
