package database

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/go-logr/logr"
	"github.com/jmoiron/sqlx"

	gomigrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const version = 20

// Wrapper wraps the PostgreSQL database.
type Wrapper struct {
	db     *sqlx.DB
	logger logr.Logger

	stopJanitor func()
}

type wrapperOpts struct {
	migrationsURL string
	janitorOn     bool
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

func WithJanitor(on bool) WrapperOpt {
	return func(w wrapperOpts) wrapperOpts {
		w.janitorOn = on
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

	wrapper := &Wrapper{db: db, logger: log}

	if options.janitorOn {
		ctx, cancel := context.WithCancel(context.Background())

		donec := make(chan struct{})
		go func() {
			wrapper.runJanitor(ctx)
			close(donec)
		}()

		wrapper.stopJanitor = func() {
			cancel()
			<-donec
		}
	}

	return wrapper, nil
}

func (w *Wrapper) Close() error {
	if w.stopJanitor != nil {
		w.stopJanitor()
	}
	return w.db.Close()
}

func connect(url string) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", url)
}

func migrate(db *sqlx.DB, sourceURL string) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{
		MigrationsTable: "migrations",
		DatabaseName:    "timeterm",
	})
	if err != nil {
		return err
	}

	migrate, err := gomigrate.NewWithDatabaseInstance(sourceURL, "timeterm", driver)
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
