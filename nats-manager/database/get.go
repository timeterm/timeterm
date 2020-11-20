package database

import "context"

func (w *Wrapper) GetOperatorSubject(ctx context.Context, name string) (subj string, err error) {
	err = w.db.GetContext(ctx, &subj, `
		SELECT subject FROM operator
		WHERE name = $1
	`, name)
	return
}

func (w *Wrapper) GetAccountSubject(ctx context.Context, name, operatorSubject string) (subj string, err error) {
	err = w.db.GetContext(ctx, &subj, `
		SELECT subject FROM account
		WHERE name = $1 AND operator_subject = $2
	`, name, operatorSubject)
	return
}

func (w *Wrapper) GetUserSubject(ctx context.Context, name, accountSubject string) (subj string, err error) {
	err = w.db.GetContext(ctx, &subj, `
		SELECT subject FROM "user"
		WHERE name = $1 AND account_subject = $2
	`, name, accountSubject)
	return
}
