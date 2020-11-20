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

const version = 1

type querier interface {
	sqlx.ExtContext
	sqlx.PreparerContext
}

type txBeginner interface {
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
}

type txBeginnerCloser interface {
	txBeginner
	io.Closer
}

type tx interface {
	Commit() error
	Rollback() error
}

// bareWrapper wraps the PostgreSQL database.
type bareWrapper struct {
	db     querier
	logger logr.Logger
}

type Wrapper struct {
	bareWrapper
	txbc txBeginnerCloser
}

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

type TxWrapper struct {
	bareWrapper
	tx tx
}

func (w *TxWrapper) Commit() error {
	return w.tx.Commit()
}

func (w *TxWrapper) Rollback() error {
	return w.tx.Rollback()
}

type wrapperOpts struct {
	migrationsURL string
}

func newWrapperOpts() wrapperOpts {
	return wrapperOpts{
		migrationsURL: "file://database/migrations",
	}
}

func createWrapperOpts(opts []WrapperOpt) wrapperOpts {
	o := newWrapperOpts()
	for _, opt := range opts {
		o = opt(o)
	}
	return o
}

type WrapperOpt func(w wrapperOpts) wrapperOpts

func WithMigrationsURL(url string) WrapperOpt {
	return func(w wrapperOpts) wrapperOpts {
		w.migrationsURL = url
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

	err = migrate(db, options.migrationsURL)
	if err != nil {
		return nil, fmt.Errorf("could not migrate database: %w", err)
	}

	wrapper := &Wrapper{
		bareWrapper: bareWrapper{db: db, logger: log},
		txbc:        db,
	}

	return wrapper, nil
}

func (w *Wrapper) Close() error {
	return w.txbc.Close()
}

func connect(url string) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", url)
}

func migrate(db *sqlx.DB, sourceURL string) error {
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

	return doMigrate(migrate)
}

func doMigrate(migrate *gomigrate.Migrate) error {
	currentVersion, isDirty, err := migrate.Version()
	if (err == nil && (isDirty || currentVersion != version)) || errors.Is(err, gomigrate.ErrNilVersion) {
		err = migrate.Migrate(version)
	}
	return err
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
