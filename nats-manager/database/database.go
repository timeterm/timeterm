package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/go-logr/logr"
	"github.com/jmoiron/sqlx"

	gomigrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// version describes the current schema version (migration number).
const version = 1

// queries is a generic interface for a sqlx database connection (pool) or transaction.
type querier interface {
	sqlx.ExtContext
	sqlx.PreparerContext
}

// txBeginner is a generic interface for a database connection (pool) allowing for beginning a transaction.
type txBeginner interface {
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
}

// txBeginnerCloser is a generic interface for a database connection (pool) allowing for beginning a transaction
// and closing itself.
type txBeginnerCloser interface {
	txBeginner
	io.Closer
}

// tx is a generic interface for a database transaction allowing committing and rolling back.
type tx interface {
	Commit() error
	Rollback() error
}

// bareWrapper wraps the PostgreSQL database.
type bareWrapper struct {
	db     querier
	logger logr.Logger
}

// Wrapper wraps the PostgreSQL database (connection pool). It can start transactions.
type Wrapper struct {
	bareWrapper
	txbc txBeginnerCloser
}

// BeginTxx creates a new TxWrapper. opts can be specified as nil to use default options.
func (w *Wrapper) BeginTxx(ctx context.Context, opts *sql.TxOptions) (*TxWrapper, error) {
	tx, err := w.txbc.BeginTxx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &TxWrapper{
		bareWrapper: bareWrapper{
			logger: w.logger,
			db:     tx,
		},
		tx: tx,
	}, nil
}

// TxWrapper is a database wrapper, with each call running in the same transaction.
type TxWrapper struct {
	bareWrapper
	tx tx
}

// Commit commits the transaction if it has not already been rolled back.
func (w *TxWrapper) Commit() error {
	return w.tx.Commit()
}

// Rollback rolls back the transaction if it has not already been committed.
func (w *TxWrapper) Rollback() error {
	return w.tx.Rollback()
}

// wrapperOpts contains options for the Wrapper.
type wrapperOpts struct {
	// The path at which the migrations can be found.
	migrationsURL  string
	migrationHooks MigrationHooks
}

// MigrationHook is a function which runs before/after each migration.
type MigrationHook func(from, to uint) error

// MigrationHooks contains hooks to run before/after each migration.
type MigrationHooks struct {
	// Hooks to run before each migration.
	Before []MigrationHook
	// Hooks to run after each migration.
	After []MigrationHook
}

// newWrapperOpts creates a new wrapperOpts struct with default settings.
func newWrapperOpts() wrapperOpts {
	return wrapperOpts{
		migrationsURL: "file://database/migrations",
	}
}

// createWrappersOpts creates a new wrapperOpts struct from a slice of WrapperOpt s.
func createWrapperOpts(opts []WrapperOpt) wrapperOpts {
	o := newWrapperOpts()
	for _, opt := range opts {
		o = opt(o)
	}
	return o
}

// WrapperOpt is a function which configures the options for the Wrapper.
type WrapperOpt func(w wrapperOpts) wrapperOpts

// WithMigrationURL sets the URL at which the migrations can be found.
func WithMigrationsURL(url string) WrapperOpt {
	return func(w wrapperOpts) wrapperOpts {
		w.migrationsURL = url
		return w
	}
}

// WithMigrationHooks adds migration hooks.
func WithMigrationHooks(h MigrationHooks) WrapperOpt {
	return func(w wrapperOpts) wrapperOpts {
		w.migrationHooks.Before = append(w.migrationHooks.Before, h.Before...)
		w.migrationHooks.After = append(w.migrationHooks.After, h.After...)
		return w
	}
}

// New opens the database and creates a new database wrapper.
func New(url string, log logr.Logger, opts ...WrapperOpt) (*Wrapper, error) {
	options := createWrapperOpts(opts)

	db, err := connect(url)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	db.MapperFunc(nameMapper)

	err = migrate(db, options.migrationsURL, options.migrationHooks)
	if err != nil {
		return nil, fmt.Errorf("could not migrate database: %w", err)
	}

	wrapper := &Wrapper{
		bareWrapper: bareWrapper{db: db, logger: log},
		txbc:        db,
	}

	return wrapper, nil
}

// Close closes the Wrapper.
func (w *Wrapper) Close() error {
	return w.txbc.Close()
}

// connect opens a new connection to a PostgreSQL database.
func connect(url string) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", url)
}

// migrate migrates the database db, with sourceURL pointing to the source URL of the migrations.
// hooks can be used for running functions before/after each migration.
func migrate(db *sqlx.DB, sourceURL string, hooks MigrationHooks) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{
		MigrationsTable: "migrations",
		DatabaseName:    "nats-manager",
	})
	if err != nil {
		return err
	}

	migrate, err := gomigrate.NewWithDatabaseInstance(sourceURL, "nats-manager", driver)
	if err != nil {
		return err
	}

	return doMigrate(migrate, hooks)
}

// doMigrate runs a Migrate, running hooks where applicable.
func doMigrate(migrate *gomigrate.Migrate, hooks MigrationHooks) error {
	currentVersion, isDirty, err := migrate.Version()
	if (err == nil && (isDirty || currentVersion != version)) || errors.Is(err, gomigrate.ErrNilVersion) {
		return walkDelta(currentVersion, version, func(cur, ver uint) error {
			// Migrate from cur -> ver, running before/after hooks.
			for _, hook := range hooks.Before {
				if err = hook(cur, ver); err != nil {
					return fmt.Errorf("error running pre-migration hook: %w", err)
				}
			}
			if err = migrate.Migrate(ver); err != nil {
				return err
			}
			for _, hook := range hooks.After {
				if err = hook(cur, ver); err != nil {
					return fmt.Errorf("error running post-migration hook: %w", err)
				}
			}
			return nil
		})
	}
	return err
}

// walkDelta runs f for each number in the range [from, to].
func walkDelta(from, to uint, f func(from, to uint) error) error {
	if from < to {
		for x := from; x < to; x++ {
			if err := f(x, x+1); err != nil {
				return err
			}
		}
	} else {
		for x := from; x > to; x-- {
			if err := f(x, x-1); err != nil {
				return err
			}
		}
	}
	return nil
}

// nameMapper maps Golang struct field names to table columns.
// For example, ID would be mapped to id and UserID would be mapped to user_id, and IDToken to id_token.
func nameMapper(fieldName string) string {
	// prevUpper saves whether the previous characters was an uppercase. If the current character is and the previous
	// character wasn't, we can put an underscore. It starts with true, so we don't get an underscore before
	// the start of the name.
	prevUpper := true

	// We'll put the column name in this character-for-character, and return the built string at the end.
	var b strings.Builder

	for _, r := range fieldName {
		// If the current character is uppercase and the previous wasn't, we can put an underscore.
		if unicode.IsUpper(r) {
			if !prevUpper {
				b.WriteByte('_')
				prevUpper = true
			}
		} else {
			prevUpper = false
		}

		// Write the character in lowercase, we don't want any uppercase letters in the column name.
		b.WriteRune(unicode.ToLower(r))
	}

	return b.String()
}
