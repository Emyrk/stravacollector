package queue

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vgarvardt/gue/v5"
	"github.com/vgarvardt/gue/v5/adapter/pgxv5"
	zadapter "github.com/vgarvardt/gue/v5/adapter/zerolog"
)

type Options struct {
	DBURL  string
	Logger zerolog.Logger
}

// Manager will handle all queue related operations and jobs
type Manager struct {
	Client *gue.Client

	pool *pgxpool.Pool
}

func New(ctx context.Context, opts Options) (*Manager, error) {
	cfg, err := pgxpool.ParseConfig(opts.DBURL)
	if err != nil {
		return nil, fmt.Errorf("parse postgres db url: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("new pool: %w", err)
	}

	poolAdapter := pgxv5.NewConnPool(pool)
	cli, err := gue.NewClient(poolAdapter, gue.WithClientLogger(zadapter.New(opts.Logger)))
	if err != nil {
		return nil, fmt.Errorf("new client: %w", err)
	}

	return &Manager{
		Client: cli,
		pool:   pool,
	}, nil
}

func (m *Manager) Close() {
	m.pool.Close()
}
