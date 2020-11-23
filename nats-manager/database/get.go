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

func (w *bareWrapper) GetAccountSubject(ctx context.Context, name, operatorName string) (subj string, err error) {
	err = sqlx.GetContext(ctx, w.db, &subj, `
		SELECT a.subject FROM account AS a
		INNER JOIN operator o ON o.subject = a.operator_subject
		WHERE a.name = $1 
		  AND o.name = $2
	`, name, operatorName)
	return
}

func (w *bareWrapper) GetUserSubject(ctx context.Context, name, accountName, operatorName string) (subj string, err error) {
	err = sqlx.GetContext(ctx, w.db, &subj, `
		SELECT u.subject FROM "user" AS u
		INNER JOIN account a ON a.subject = u.account_subject
		INNER JOIN operator o ON o.subject = a.operator_subject
		WHERE u.name = $1 
		  AND a.name = $2 
		  AND o.name = $3
	`, name, accountName, operatorName)
	return
}
