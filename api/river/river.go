package river

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/Emyrk/strava/database"
	"github.com/Emyrk/strava/internal/debounce"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/riverlog"
	"github.com/riverqueue/river/rivertype"
	"github.com/riverqueue/rivercontrib/otelriver"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog"
	promotel "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
	"riverqueue.com/riverui"
)

const (
	PriorityHighest = 1
	PriorityHigh    = 2
	PriorityDefault = 3
	PriorityLow     = 4
)

const (
	riverBackloadQueue = "backload_queue"
	riverStravaQueue   = "strava_queue"
	riverControlQueue  = "control_queue"
	riverDatabaseQueue = "database_operations_queue"
)

type Options struct {
	DBURL    string
	Logger   zerolog.Logger
	DB       database.Store
	OAuthCfg *oauth2.Config
	Registry *prometheus.Registry
}

type Manager struct {
	logger   zerolog.Logger
	db       database.Store
	pool     *pgxpool.Pool
	cli      *river.Client[pgx.Tx]
	oauthCfg *oauth2.Config

	rateLimitLogger *debounce.Debouncer
	appCtx          context.Context

	managerMetrics
}

func New(ctx context.Context, opts Options) (*Manager, error) {
	cfg, err := database.PoolConfig(opts.DBURL)
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
	hourly, err := cron.ParseStandard("0 * * * *")
	if err != nil {
		return nil, fmt.Errorf("parse cron schedule: %w", err)
	}

	sixly, err := cron.ParseStandard("0 0/6 * * *")
	if err != nil {
		return nil, fmt.Errorf("parse cron schedule: %w", err)
	}

	periodicJobs := []*river.PeriodicJob{
		river.NewPeriodicJob(
			// Always resume after some amount of time to prevent the queue from sleeping
			// forever.
			hourly,
			func() (river.JobArgs, *river.InsertOpts) {
				return ResumeArgs{
					Queues: []string{riverStravaQueue, riverBackloadQueue},
				}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: true, ID: "strava_resume"},
		),
		// TODO: When hugel goes live, this needs to be updated more frequently
		river.NewPeriodicJob(
			sixly,
			func() (river.JobArgs, *river.InsertOpts) {
				return RefreshViewsArgs{}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: false, ID: "refresh_views"},
		),
		river.NewPeriodicJob(
			hourly,
			func() (river.JobArgs, *river.InsertOpts) {
				return ReloadSegmentsArgs{}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: true, ID: "reload_segments"},
		),
	}

	exporter, err := promotel.New(promotel.WithRegisterer(opts.Registry))
	if err != nil {
		return nil, fmt.Errorf("new prometheus exporter: %w", err)
	}
	otelMeterProvider := metric.NewMeterProvider(metric.WithReader(exporter))

	riverClient, err := river.NewClient(riverpgxv5.New(pool), (&river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 1},
			riverStravaQueue:   {MaxWorkers: 1},
			riverControlQueue:  {MaxWorkers: 1},
			riverDatabaseQueue: {MaxWorkers: 1},
			riverBackloadQueue: {MaxWorkers: 1},
		},
		Workers: workers,
		Middleware: []rivertype.Middleware{
			otelriver.NewMiddleware(&otelriver.MiddlewareConfig{
				MeterProvider: otelMeterProvider,
			}),
			riverlog.NewMiddleware(func(w io.Writer) slog.Handler {
				return slog.NewJSONHandler(w, nil)
			}, nil),
		},

		CancelledJobRetentionPeriod: time.Hour * 24 * 7,
		CompletedJobRetentionPeriod: time.Hour * 24,
		DiscardedJobRetentionPeriod: time.Hour * 24 * 30,
		Logger:                      slog.New(slogzerolog.Option{Level: slog.LevelInfo, Logger: &opts.Logger}.NewZerologHandler()),
		PeriodicJobs:                periodicJobs,
	}).WithDefaults())
	if err != nil {
		return nil, fmt.Errorf("new river: %w", err)
	}

	m := &Manager{
		logger:          opts.Logger,
		db:              opts.DB,
		pool:            pool,
		cli:             riverClient,
		rateLimitLogger: debounce.New(time.Minute * 7),
		oauthCfg:        opts.OAuthCfg,
		appCtx:          ctx,
	}

	m.initWorkers(workers)
	m.initMetrics(opts.Registry)
	m.background(ctx)

	if err := riverClient.Start(ctx); err != nil {
		return nil, fmt.Errorf("start river client: %w", err)
	}

	return m, nil
}

func (m *Manager) Close(ctx context.Context) error {
	grp := &errgroup.Group{}
	grp.Go(func() error {
		err := m.cli.Stop(ctx)
		if err != nil {
			return fmt.Errorf("stop river client: %w", err)
		}
		select {
		case <-ctx.Done():
		case <-m.cli.Stopped():
		}
		return nil
	})

	grpErr := grp.Wait()
	m.pool.Close()
	m.logger.Info().Msgf("River client stopped")
	return grpErr
}

func (m *Manager) Attach(ctx context.Context, r chi.Router) error {
	opts := &riverui.ServerOpts{
		Client:                   m.cli,
		DB:                       m.pool,
		DevMode:                  false,
		JobListHideArgsByDefault: false,
		LiveFS:                   false,
		Logger:                   slog.New(slogzerolog.Option{Level: slog.LevelInfo, Logger: &m.logger}.NewZerologHandler()),
		Prefix:                   "/river",
	}

	srv, err := riverui.NewServer(opts)
	if err != nil {
		return fmt.Errorf("new riverui server: %w", err)
	}

	err = srv.Start(ctx)
	if err != nil {
		return fmt.Errorf("start riverui server: %w", err)
	}

	r.Mount("/river", srv)
	m.logger.Info().
		Str("path", "/river").
		Msg("River UI server started")
	return nil
}

func (m *Manager) initWorkers(workers *river.Workers) {
	river.AddWorker[FetchActivityArgs](workers, &FetchActivityWorker{
		mgr: m,
	})
	river.AddWorker[ResumeArgs](workers, &ResumeWorker{
		mgr: m,
	})
	river.AddWorker[UpdateActivityArgs](workers, &UpdateActivityWorker{
		mgr: m,
	})
	river.AddWorker[RefreshViewsArgs](workers, &RefreshViewsWorker{
		mgr: m,
	})
	river.AddWorker[ReloadSegmentsArgs](workers, &ReloadSegmentsWorker{
		mgr: m,
	})
	river.AddWorker[ForwardLoadArgs](workers, &ForwardLoadWorker{
		mgr: m,
	})
}

func (m *Manager) StravaSnooze(ctx context.Context) error {
	// TODO: Pause the queue until the next interval, not just 15minutes
	_ = river.RecordOutput(ctx, "hitting strava rate limit, job going to pause for 15 minutes")
	_ = m.Pause(time.Now().Add(time.Minute*15), riverStravaQueue)
	_ = m.Pause(time.Now().Add(time.Minute*15), riverBackloadQueue)
	return river.JobSnooze(time.Minute * 15)
}

func (m *Manager) StravaMaintaince(ctx context.Context, reason string) error {
	// TODO: Pause the queue until the next interval, not just 15minutes
	_ = river.RecordOutput(ctx, fmt.Sprintf("strava is offline, or in maintaince: %s", reason))
	_ = m.Pause(time.Now().Add(time.Minute*15), riverStravaQueue)
	_ = m.Pause(time.Now().Add(time.Minute*15), riverBackloadQueue)
	return river.JobSnooze(time.Minute * 15)
}
