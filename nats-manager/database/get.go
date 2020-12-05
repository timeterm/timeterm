package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

type Pagination struct {
	Total int
	Limit int
}

type paginatedJWTs struct {
	Pagination
	Subjects []string
}

func (w *bareWrapper) ListJWTs(ctx context.Context, offset int) (p paginatedJWTs, err error) {
	const limit = 100

	p.Limit = limit
	if err = sqlx.GetContext(ctx, w.db, &p.Total, `SELECT COUNT(*) FROM jwt`); err != nil {
		return p, err
	}

	return p, sqlx.SelectContext(ctx, w.db, &p.Subjects, `
		SELECT subject FROM jwt
		ORDER BY subject
		LIMIT $1
		OFFSET $2
	`, limit, offset)
}

type Operator struct {
	Subject string
	Name    string
}

type PaginatedOperators struct {
	Pagination
	Operators []Operator
}

func (w *bareWrapper) ListOperatorsRe(ctx context.Context, offset int, regex string) (p PaginatedOperators, err error) {
	const limit = 100

	p.Limit = limit
	if err = sqlx.GetContext(ctx, w.db, &p.Total, `
		SELECT COUNT(*) FROM operator WHERE name ~ $1
	`, regex); err != nil {
		return p, err
	}

	return p, sqlx.SelectContext(ctx, w.db, &p.Operators, `
		SELECT * FROM operator
		WHERE name ~ $1
		ORDER BY subject
		LIMIT $2
		OFFSET $3
	`, regex, limit, offset)
}

type Account struct {
	Subject         string
	Name            string
	OperatorSubject string
	OperatorName    string
}

type PaginatedAccounts struct {
	Pagination
	Accounts []Account
}

func (w *bareWrapper) ListAccountsRe(
	ctx context.Context,
	offset int,
	nameRegex,
	operatorNameRegex string,
) (p PaginatedAccounts, err error) {
	const limit = 100

	p.Limit = limit
	if err = sqlx.GetContext(ctx, w.db, &p.Total, `
		SELECT COUNT(*) FROM account AS a
		INNER JOIN operator o on a.operator_subject = o.subject	
		WHERE a.name ~ $1
		AND o.name ~ $2
	`, nameRegex, operatorNameRegex); err != nil {
		return p, err
	}

	rows, err := w.db.QueryxContext(ctx, `
		SELECT a.subject, a.name, a.operator_subject, o.name FROM account AS a
		INNER JOIN operator o on a.operator_subject = o.subject	
		WHERE a.name ~ $1
		AND o.name ~ $2
		ORDER BY a.subject
		LIMIT $3
		OFFSET $4
	`, nameRegex, operatorNameRegex, limit, offset)
	if err != nil {
		return p, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var a Account
		if err = rows.Scan(&a.Subject, &a.Name, &a.OperatorSubject, &a.OperatorName); err != nil {
			return p, err
		}
		p.Accounts = append(p.Accounts, a)
	}

	return
}

type User struct {
	Subject         string
	Name            string
	AccountSubject  string
	AccountName     string
	OperatorSubject string
	OperatorName    string
}

type PaginatedUsers struct {
	Pagination
	Users []User
}

func (w *bareWrapper) ListUsersRe(
	ctx context.Context,
	offset int,
	nameRegex,
	accountNameRegex,
	operatorNameRegex string,
) (p PaginatedUsers, err error) {
	const limit = 100

	p.Limit = limit
	if err = sqlx.GetContext(ctx, w.db, &p.Total, `
		SELECT COUNT(*) FROM "user" AS u
		INNER JOIN account a on u.account_subject = a.subject
		INNER JOIN operator o on a.operator_subject = o.subject	
		WHERE u.name ~ $1
		AND a.name ~ $2
		AND o.name ~ $3
	`, nameRegex, accountNameRegex, operatorNameRegex); err != nil {
		return p, err
	}

	rows, err := w.db.QueryxContext(ctx, `
		SELECT u.subject, u.name, u.account_subject, a.name, a.operator_subject, o.name FROM "user" AS u
		INNER JOIN account a on u.account_subject = a.subject
		INNER JOIN operator o on a.operator_subject = o.subject	
		WHERE u.name ~ $1
		AND a.name ~ $2
		AND o.name ~ $3
		ORDER BY u.subject
		LIMIT $4
		OFFSET $5
	`, nameRegex, accountNameRegex, operatorNameRegex, limit, offset)
	if err != nil {
		return p, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var u User
		err = rows.Scan(&u.Subject, &u.Name, &u.AccountSubject, &u.AccountName, &u.OperatorSubject, &u.OperatorName)
		if err != nil {
			return p, err
		}
		p.Users = append(p.Users, u)
	}

	return
}

func (w *bareWrapper) WalkJWTs(ctx context.Context, f func(subj string) bool) error {
	offset := 0
	for {
		jwts, err := w.ListJWTs(ctx, offset)
		if err != nil {
			return fmt.Errorf("could not retrieve JWTs with offset %d: %w", offset, err)
		}
		if len(jwts.Subjects) == 0 {
			break
		}
		offset += jwts.Limit

		for _, sub := range jwts.Subjects {
			if !f(sub) {
				return nil
			}
		}
	}

	return nil
}

func (w *bareWrapper) WalkOperatorSubjectsRe(ctx context.Context, regex string, f func(o Operator) bool) error {
	offset := 0
	for {
		jwts, err := w.ListOperatorsRe(ctx, offset, regex)
		if err != nil {
			return fmt.Errorf("could not retrieve operators with offset %d: %w", offset, err)
		}
		if len(jwts.Operators) == 0 {
			break
		}
		offset += jwts.Limit

		for _, o := range jwts.Operators {
			if !f(o) {
				return nil
			}
		}
	}

	return nil
}

func (w *bareWrapper) WalkAccountSubjectsRe(
	ctx context.Context,
	nameRegex,
	operatorNameRegex string,
	f func(a Account) bool,
) error {
	offset := 0
	for {
		jwts, err := w.ListAccountsRe(ctx, offset, nameRegex, operatorNameRegex)
		if err != nil {
			return fmt.Errorf("could not retrieve accounts with offset %d: %w", offset, err)
		}
		if len(jwts.Accounts) == 0 {
			break
		}
		offset += jwts.Limit

		for _, a := range jwts.Accounts {
			if !f(a) {
				return nil
			}
		}
	}

	return nil
}

func (w *bareWrapper) WalkUserSubjectsRe(
	ctx context.Context,
	nameRegex,
	accountNameRegex,
	operatorNameRegex string,
	f func(u User) bool,
) error {
	offset := 0
	for {
		jwts, err := w.ListUsersRe(ctx, offset, nameRegex, accountNameRegex, operatorNameRegex)
		if err != nil {
			return fmt.Errorf("could not retrieve users with offset %d: %w", offset, err)
		}
		if len(jwts.Users) == 0 {
			break
		}
		offset += jwts.Limit

		for _, u := range jwts.Users {
			if !f(u) {
				return nil
			}
		}
	}

	return nil
}

func (w *bareWrapper) GetJWTMigrationVersion(ctx context.Context) (int, error) {
	var jwtMigrationVersion int

	err := sqlx.GetContext(ctx, w.db, &jwtMigrationVersion, `SELECT version FROM jwt_migration LIMIT 1`)
	if errors.Is(err, sql.ErrNoRows) {
		err = nil
		jwtMigrationVersion = 0
	}
	return jwtMigrationVersion, err
}
