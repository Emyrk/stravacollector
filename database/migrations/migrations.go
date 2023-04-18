package migrations

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"golang.org/x/xerrors"
)

//go:embed *.sql
var migrations embed.FS

func setup(db *sql.DB) (source.Driver, *migrate.Migrate, error) {
	ctx := context.Background()
	sourceDriver, err := iofs.New(migrations, ".")
	if err != nil {
		return nil, nil, fmt.Errorf("create iofs: %w", err)
	}

	// there is a postgres.WithInstance() method that takes the DB instance,
	// but, when you close the resulting Migrate, it closes the DB, which
	// we don't want.  Instead, create just a connection that will get closed
	// when migration is done.
	conn, err := db.Conn(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("postgres connection: %w", err)
	}
	dbDriver, err := postgres.WithConnection(ctx, conn, &postgres.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("wrap postgres connection: %w", err)
	}

	m, err := migrate.NewWithInstance("", sourceDriver, "", dbDriver)
	if err != nil {
		return nil, nil, fmt.Errorf("new migrate instance: %w", err)
	}

	return sourceDriver, m, nil
}

// Up runs SQL migrations to ensure the database schema is up-to-date.
func Up(db *sql.DB) (retErr error) {
	_, m, err := setup(db)
	if err != nil {
		return fmt.Errorf("migrate setup: %w", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if retErr != nil {
			return
		}
		if dbErr != nil {
			retErr = dbErr
			return
		}
		retErr = srcErr
	}()

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			// It's OK if no changes happened!
			return nil
		}

		return fmt.Errorf("up: %w", err)
	}

	return nil
}

// Down runs all down SQL migrations.
func Down(db *sql.DB) error {
	_, m, err := setup(db)
	if err != nil {
		return xerrors.Errorf("migrate setup: %w", err)
	}

	err = m.Down()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			// It's OK if no changes happened!
			return nil
		}

		return xerrors.Errorf("down: %w", err)
	}

	return nil
}
