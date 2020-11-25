package database

import "context"

func (w *Wrapper) CreateOperator(ctx context.Context, name, subject string) error {
	tx, err := w.txbc.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO jwt (subject)
		VALUES ($1)
	`, subject)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO operator (name, subject)
		VALUES ($1, $2)
	`, name, subject)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (w *Wrapper) CreateAccount(ctx context.Context, name, subject, operatorSubject string) error {
	tx, err := w.txbc.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO jwt (subject)
		VALUES ($1)
	`, subject)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO account (name, subject, operator_subject)
		VALUES ($1, $2, $3)
	`, name, subject, operatorSubject)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (w *Wrapper) CreateUser(ctx context.Context, name, subject, accountSubject string) error {
	tx, err := w.txbc.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO jwt (subject)
		VALUES ($1)
	`, subject)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO "user" (name, subject, account_subject)
		VALUES ($1, $2, $3)
	`, name, subject, accountSubject)
	if err != nil {
		return err
	}

	return tx.Commit()
}
