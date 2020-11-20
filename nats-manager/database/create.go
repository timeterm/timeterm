package database

import "context"

func (w *Wrapper) CreateOperator(ctx context.Context, name, subject string) error {
	_, err := w.db.ExecContext(ctx, `
		INSERT INTO operator (name, subject)
		VALUES ($1, $2)
	`, name, subject)
	return err
}

func (w *Wrapper) CreateAccount(ctx context.Context, name, subject, operatorSubject string) error {
	_, err := w.db.ExecContext(ctx, `
		INSERT INTO account (name, subject, operator_subject)
		VALUES ($1, $2, $3)
	`, name, subject, operatorSubject)
	return err
}

func (w *Wrapper) CreateUser(ctx context.Context, name, subject, accountSubject string) error {
	_, err := w.db.ExecContext(ctx, `
		INSERT INTO user (name, subject, account_subject)
		VALUES ($1, $2, $3, $4)
	`, name, subject, accountSubject)
	return err
}
