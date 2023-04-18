package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/xerrors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// Store contains all queryable database functions.
// It extends the generated interface to add transaction support.
type Store interface {
	sqlcQuerier

	Ping(ctx context.Context) (time.Duration, error)
	InTx(func(Store) error, *sql.TxOptions) error
}

type sqlQuerier struct {
	sdb *sqlx.DB
	db  DBTX
}

// New creates a new database store using a SQL database connection.
func New(sdb *sql.DB) Store {
	dbx := sqlx.NewDb(sdb, "postgres")
	return &sqlQuerier{
		db:  dbx,
		sdb: dbx,
	}
}

// Ping returns the time it takes to ping the database.
func (q *sqlQuerier) Ping(ctx context.Context) (time.Duration, error) {
	start := time.Now()
	err := q.sdb.PingContext(ctx)
	return time.Since(start), err
}

func (q *sqlQuerier) InTx(function func(Store) error, txOpts *sql.TxOptions) error {
	_, inTx := q.db.(*sqlx.Tx)
	isolation := sql.LevelDefault
	if txOpts != nil {
		isolation = txOpts.Isolation
	}

	// If we are not already in a transaction, and we are running in serializable
	// mode, we need to run the transaction in a retry loop. The caller should be
	// prepared to allow retries if using serializable mode.
	// If we are in a transaction already, the parent InTx call will handle the retry.
	// We do not want to duplicate those retries.
	if !inTx && isolation == sql.LevelSerializable {
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
func (q *sqlQuerier) runTx(function func(Store) error, txOpts *sql.TxOptions) error {
	if _, ok := q.db.(*sqlx.Tx); ok {
		// If the current inner "db" is already a transaction, we just reuse it.
		// We do not need to handle commit/rollback as the outer tx will handle
		// that.
		err := function(q)
		if err != nil {
			return xerrors.Errorf("execute transaction: %w", err)
		}
		return nil
	}

	transaction, err := q.sdb.BeginTxx(context.Background(), txOpts)
	if err != nil {
		return xerrors.Errorf("begin transaction: %w", err)
	}
	defer func() {
		rerr := transaction.Rollback()
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
	err = transaction.Commit()
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
