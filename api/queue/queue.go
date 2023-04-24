package queue

import (
	"context"
	"fmt"
	"github.com/Emyrk/strava/strava/stravalimit"
	"time"

	"golang.org/x/oauth2"

	"github.com/Emyrk/strava/database"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	exp "github.com/vgarvardt/backoff"
	"github.com/vgarvardt/gue/v5"
	"github.com/vgarvardt/gue/v5/adapter/pgxv5"
	zadapter "github.com/vgarvardt/gue/v5/adapter/zerolog"
)

const (
	fetchActivityJob = "fetch_activity"

	stravaFetchQueue = "queue_strava_fetch"
)

type Options struct {
	DBURL    string
	Logger   zerolog.Logger
	DB       database.Store
	OAuthCfg *oauth2.Config
}

// Manager will handle all queue related operations and jobs
type Manager struct {
	Client *gue.Client

	// Pool is used for the queuing library
	pool *pgxpool.Pool

	// DB is used by jobs
	DB database.Store

	Logger   zerolog.Logger
	OAuthCfg *oauth2.Config

	cancel context.CancelFunc
}

func New(ctx context.Context, opts Options) (*Manager, error) {
	cfg, err := pgxpool.ParseConfig(opts.DBURL)
	if err != nil {
		return nil, fmt.Errorf("parse postgres db url: %w", err)
	}
	// Small number of conns
	cfg.MaxConns = 2

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("new pool: %w", err)
	}

	poolAdapter := pgxv5.NewConnPool(pool)
	cli, err := gue.NewClient(poolAdapter,
		gue.WithClientLogger(zadapter.New(opts.Logger)),
		gue.WithClientBackoff(gue.NewExponentialBackoff(exp.Config{
			BaseDelay:  time.Second * 5,
			Multiplier: 1.6,
			Jitter:     0.2,
			MaxDelay:   time.Minute * 15,
		})),
	)
	if err != nil {
		return nil, fmt.Errorf("new client: %w", err)
	}

	return &Manager{
		Client:   cli,
		pool:     pool,
		DB:       opts.DB,
		OAuthCfg: opts.OAuthCfg,
		Logger:   opts.Logger,
	}, nil
}

func (m *Manager) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	m.cancel = cancel

	worker, err := gue.NewWorker(m.Client, m.workMap(),
		gue.WithWorkerHooksJobDone(func(ctx context.Context, j *gue.Job, err error) {
			// TODO: If this is a strava too many requests, we need to sleep.
			if err != nil {
				m.Logger.Error().
					Err(err).
					Str("job_id", j.ID.String()).
					Str("job", j.Type).
					Str("queue", j.Queue).
					Int32("err_count", j.ErrorCount).
					Str("last_error", j.LastError.String).
					Msg("job failed")
			}
		}))
	if err != nil {
		return fmt.Errorf("new worker: %w", err)
	}

	go func() {
		err := worker.Run(ctx)
		if err != nil {
			m.Logger.Error().Err(err).Msg("worker error")
		}
		cancel()
	}()

	return nil
}

func (m *Manager) workMap() gue.WorkMap {
	return gue.WorkMap{
		"online": func(ctx context.Context, j *gue.Job) error {
			m.Logger.Info().Msg("worker online")
			return nil
		},
		fetchActivityJob: func(ctx context.Context, j *gue.Job) error {
			return m.fetchActivity(ctx, j)
		},
	}
}

func (m *Manager) stravaCheck(j *gue.Job, calls int64) error{
	logger := jobLogFields(m.Logger, j)
	ok, limitLogger := stravalimit.CanLogger(1, 100, logger)
	if !ok {
		limitLogger.Error().
			Msg("hitting strava rate limit, job going to fail and try again later")
		return fmt.Errorf("hitting strava rate limit, failing job to try later")
	}
	return nil
}

func jobLogFields(logger zerolog.Logger, j *gue.Job) zerolog.Logger {
	return logger.With().
		Str("job_id", j.ID.String()).
		Str("job", j.Type).
		Str("queue", j.Queue).
		Int32("err_count", j.ErrorCount).
		Str("last_error", j.LastError.String).
		Logger()
}

func (m *Manager) Close() {
	m.cancel()
	m.pool.Close()
}
