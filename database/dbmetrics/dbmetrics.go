package dbmetrics

import (
	"strconv"
	"time"

	"github.com/Emyrk/strava/database"
	"github.com/jackc/pgx/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

type metricsStore struct {
	database.Store
	logger zerolog.Logger
	// txDuration is how long transactions take to execute.
	txDuration *prometheus.HistogramVec
}

// NewDBMetrics returns a database.Store that registers metrics for the database
// but does not handle individual queries.
// metricsStore is intended to always be used, because queryMetrics are a bit
// too verbose for many use cases.
func NewDBMetrics(s database.Store, logger zerolog.Logger, reg prometheus.Registerer) database.Store {
	txDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "coderd",
		Subsystem: "db",
		Name:      "tx_duration_seconds",
		Help:      "Duration of transactions in seconds.",
		Buckets:   prometheus.DefBuckets,
	}, []string{
		"success", // Did the InTx function return an error?
	})
	reg.MustRegister(txDuration)
	return &metricsStore{
		Store:      s,
		txDuration: txDuration,
		logger:     logger,
	}
}

func (m metricsStore) InTx(f func(database.Store) error, options *pgx.TxOptions) error {
	start := time.Now()
	err := m.Store.InTx(f, options)
	dur := time.Since(start)
	// The number of unique label combinations is
	// 2 x #IDs x #of buckets
	// So IDs should be used sparingly to prevent too much bloat.
	m.txDuration.With(prometheus.Labels{
		"success": strconv.FormatBool(err == nil),
	}).Observe(dur.Seconds())

	return err
}
