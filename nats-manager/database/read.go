package database

import "context"

func (w *Wrapper) GetOperatorByName(ctx context.Context, name, subject string) error {
	_, err := w.db.ExecContext(ctx, `
		INSERT INTO operator (name, subject)
		VALUES ($1, $2)
	`, name, subject)
	return err
}
