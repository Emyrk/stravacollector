package river

import (
	"context"
	"fmt"

	"github.com/Emyrk/strava/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
)

type Options struct {
	DBURL    string
	Logger   zerolog.Logger
	DB       database.Store
	OAuthCfg *oauth2.Config
	Registry *prometheus.Registry
}

type Manager struct {
	logger zerolog.Logger
	db     database.Store
	pool   *pgxpool.Pool
	cli    *river.Client[pgx.Tx]
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

	workers := river.NewWorkers()
	riverClient, err := river.NewClient(riverpgxv5.New(pool), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 1},
		},
		Workers: workers,
	})
	if err != nil {
		return nil, fmt.Errorf("new river: %w", err)
	}

	if err := riverClient.Start(ctx); err != nil {
		return nil, fmt.Errorf("start river client: %w", err)
	}

	m := &Manager{
		logger: opts.Logger,
		db:     opts.DB,
		pool:   pool,
		cli:    riverClient,
	}

	m.initWorkers(workers)

	return m, nil
}

func (m *Manager) initWorkers(workers *river.Workers) {

}
