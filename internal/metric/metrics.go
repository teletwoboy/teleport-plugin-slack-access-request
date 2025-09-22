package metric

import (
	"teleport-plugin-slack-access-request/internal/database"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var Registry *prometheus.Registry

// HTTP Metrics
var (
	HTTPRequestsTotal    *prometheus.CounterVec
	HTTPRequestDuration  *prometheus.HistogramVec
	HTTPInFlightRequests prometheus.Gauge
)

// Init initializes all collectors
func Init(db *database.DB) {
	Registry = prometheus.NewRegistry()

	// Default Collectors
	Registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	// HTTP
	http := NewHTTPMetrics()
	http.MustRegister(Registry)

	// DB
	Registry.MustRegister(NewDBStatsCollector(db))
}
