package river

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type managerMetrics struct {
	rideActivitySummaries prometheus.Gauge
	rideActivityDetails   prometheus.Gauge
}

func (m *Manager) initMetrics(registry *prometheus.Registry) {
	factory := promauto.With(registry)
	m.rideActivitySummaries = factory.NewGauge(prometheus.GaugeOpts{
		Namespace: "strava",
		Subsystem: "manager",
		Name:      "activity_summary_total",
		Help:      "The total number of ride activities known about",
	})
	m.rideActivityDetails = factory.NewGauge(prometheus.GaugeOpts{
		Namespace: "strava",
		Subsystem: "manager",
		Name:      "activity_detail_total",
		Help:      "The total number of ride activities synced",
	})
}
