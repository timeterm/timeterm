package database

import (
	"context"

	"github.com/jmoiron/sqlx"
)

func (w *bareWrapper) GetOperatorSubject(ctx context.Context, name string) (subj string, err error) {
	err = sqlx.GetContext(ctx, w.db, &subj, `
		SELECT subject FROM operator
		WHERE name = $1
	`, name)
	return
}

func (w *bareWrapper) GetAccountSubject(ctx context.Context, name, operatorSubject string) (subj string, err error) {
	err = sqlx.GetContext(ctx, w.db, &subj, `
		SELECT subject FROM account
		WHERE name = $1 AND operator_subject = $2
	`, name, operatorSubject)
	return
}

func (w *bareWrapper) GetUserSubject(ctx context.Context, name, accountSubject string) (subj string, err error) {
	err = sqlx.GetContext(ctx, w.db, &subj, `
		SELECT subject FROM "user"
		WHERE name = $1 AND account_subject = $2
	`, name, accountSubject)
	return
}
