package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Emyrk/strava/database/migrations"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
	"github.com/rs/zerolog"
	"golang.org/x/xerrors"
)

// Store contains all queryable database functions.
// It extends the generated interface to add transaction support.
type Store interface {
	sqlcQuerier
	manualQuerier

	Ping(ctx context.Context) (time.Duration, error)
	InTx(func(Store) error, *pgx.TxOptions) error
	Close() error
}

// DBTX represents a database connection or transaction.
type DBTX interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

type sqlQuerier struct {
	sdb *pgxpool.Pool
	db  DBTX
}

func NewPostgresDB(ctx context.Context, logger zerolog.Logger, dbURL string) (Store, error) {
	logger = logger.With().Str("db_url", dbURL).Logger()
	logger.Info().Msg("connecting to postgres database")

	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("parse postgres db url: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}

	pingCtx, pingCancel := context.WithTimeout(ctx, 15*time.Second)
	defer pingCancel()
	err = pool.Ping(pingCtx)
	if err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	err = migrations.Up(pool)
	if err != nil {
		return nil, fmt.Errorf("migrate up: %w", err)
	}

	return New(pool), nil
}

// New creates a new database store using a SQL database connection.
func New(sdb *pgxpool.Pool) Store {
	return &sqlQuerier{
		db:  sdb,
		sdb: sdb,
	}
}

func (q *sqlQuerier) Close() error {
	q.sdb.Close()
	return nil
}

// Ping returns the time it takes to ping the database.
func (q *sqlQuerier) Ping(ctx context.Context) (time.Duration, error) {
	start := time.Now()
	err := q.sdb.Ping(ctx)
	return time.Since(start), err
}

func (q *sqlQuerier) InTx(function func(Store) error, txOpts *pgx.TxOptions) error {
	_, inTx := q.db.(*pgxpool.Tx)
	isolation := pgx.ReadCommitted
	if txOpts != nil {
		isolation = txOpts.IsoLevel
	}

	// If we are not already in a transaction, and we are running in serializable
	// mode, we need to run the transaction in a retry loop. The caller should be
	// prepared to allow retries if using serializable mode.
	// If we are in a transaction already, the parent InTx call will handle the retry.
	// We do not want to duplicate those retries.
	if !inTx && isolation == pgx.Serializable {
		// This is an arbitrarily chosen number.
		const retryAmount = 3
		var err error
		attempts := 0
		for attempts = 0; attempts < retryAmount; attempts++ {
			err = q.runTx(function, txOpts)
			if err == nil {
				// Transaction succeeded.
				return nil
			}
			if err != nil && !IsSerializedError(err) {
				// We should only retry if the error is a serialization error.
				return err
			}
		}
		// Transaction kept failing in serializable mode.
		return xerrors.Errorf("transaction failed after %d attempts: %w", attempts, err)
	}
	return q.runTx(function, txOpts)
}

// InTx performs database operations inside a transaction.
func (q *sqlQuerier) runTx(function func(Store) error, txOpts *pgx.TxOptions) error {
	if _, ok := q.db.(*pgxpool.Tx); ok {
		// If the current inner "db" is already a transaction, we just reuse it.
		// We do not need to handle commit/rollback as the outer tx will handle
		// that.
		err := function(q)
		if err != nil {
			return xerrors.Errorf("execute transaction: %w", err)
		}
		return nil
	}

	opts := txOpts
	if opts == nil {
		opts = &pgx.TxOptions{
			IsoLevel: pgx.ReadCommitted,
		}
	}

	transaction, err := q.sdb.BeginTx(context.Background(), *opts)
	if err != nil {
		return xerrors.Errorf("begin transaction: %w", err)
	}
	defer func() {
		rerr := transaction.Rollback(context.Background())
		if rerr == nil || errors.Is(rerr, sql.ErrTxDone) {
			// no need to do anything, tx committed successfully
			return
		}
		// couldn't roll back for some reason, extend returned error
		err = xerrors.Errorf("defer (%s): %w", rerr.Error(), err)
	}()
	err = function(&sqlQuerier{db: transaction})
	if err != nil {
		return xerrors.Errorf("execute transaction: %w", err)
	}
	err = transaction.Commit(context.Background())
	if err != nil {
		return xerrors.Errorf("commit transaction: %w", err)
	}
	return nil
}

func IsSerializedError(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code.Name() == "serialization_failure"
	}
	return false
}
