package database

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	gomigrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const version = 1

// Wrapper wraps the PostgreSQL database.
type Wrapper struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// New opens the database and creates a new database wrapper.
func New(url string, logger *zap.Logger) (*Wrapper, error) {
	db, err := connect(url)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	db.MapperFunc(nameMapper)

	err = migrate(db)
	if err != nil {
		return nil, fmt.Errorf("could not migrate database: %w", err)
	}

	return &Wrapper{db: db, logger: logger}, nil
}

func (w *Wrapper) Close() error {
	return w.db.Close()
}

func connect(url string) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", url)
}

func migrate(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{
		MigrationsTable: "migrations",
		DatabaseName:    "timeterm",
	})
	if err != nil {
		return err
	}

	migrate, err := gomigrate.NewWithDatabaseInstance("file://database/migrations", "timeterm", driver)
	if err != nil {
		return err
	}

	return doMigrate(migrate)
}

func doMigrate(migrate *gomigrate.Migrate) error {
	currentVersion, isDirty, err := migrate.Version()
	if err != nil && (isDirty || currentVersion != version) {
		err = migrate.Migrate(version)
	}
	return err
}

// nameMapper maps Golang struct field names to table columns.
// For example, ID would be mapped to id and UserID would be mapped to user_id.
func nameMapper(s string) string {
	prevUpper := true // we start with true, so we don't get an underscore before every name.

	var b strings.Builder
	for _, r := range s {
		if unicode.IsUpper(r) {
			if !prevUpper {
				b.WriteByte('_')
				prevUpper = true
			}
		} else {
			prevUpper = false
		}
		b.WriteRune(unicode.ToLower(r))
	}

	return b.String()
}
