package httpmw

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func PrometheusMW(register prometheus.Registerer) func(http.Handler) http.Handler {
	factory := promauto.With(register)
	requestsProcessed := factory.NewCounterVec(prometheus.CounterOpts{
		Namespace: "strava",
		Subsystem: "api",
		Name:      "requests_processed_total",
		Help:      "The total number of processed API requests",
	}, []string{"method", "path"})
	requestsConcurrent := factory.NewGauge(prometheus.GaugeOpts{
		Namespace: "strava",
		Subsystem: "api",
		Name:      "concurrent_requests",
		Help:      "The number of concurrent API requests.",
	})
	requestsDist := factory.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "strava",
		Subsystem: "api",
		Name:      "request_latencies_seconds",
		Help:      "Latency distribution of requests in seconds.",
		Buckets:   []float64{0.001, 0.005, 0.010, 0.025, 0.050, 0.100, 0.500, 1, 5, 10, 30},
	}, []string{"method", "path"})

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				start  = time.Now()
				method = r.Method
				rctx   = chi.RouteContext(r.Context())
			)

			requestsConcurrent.Inc()
			defer requestsConcurrent.Dec()

			next.ServeHTTP(w, r)

			path := rctx.RoutePattern()

			requestsProcessed.WithLabelValues(method, path).Inc()
			requestsDist.WithLabelValues(method, path).Observe(time.Since(start).Seconds())
		})
	}
}
