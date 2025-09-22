package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

// HTTPMetrics groups HTTP-related metrics
type HTTPMetrics struct {
	RequestsTotal    *prometheus.CounterVec   // 총 요청 수
	RequestDuration  *prometheus.HistogramVec // 요청 처리 시간
	InFlightRequests prometheus.Gauge         // 현재 처리 중 요청 수
}

// NewHTTPMetrics initializes HTTP metrics with labels
func NewHTTPMetrics() *HTTPMetrics {
	return &HTTPMetrics{
		RequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP request handling in seconds",
				Buckets: []float64{0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2, 5},
			},
			[]string{"method", "path", "status"},
		),
		InFlightRequests: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_inflight_requests",
				Help: "Current number of in-flight HTTP requests",
			},
		),
	}
}

// MustRegister registers all HTTP metrics into a registry
func (m *HTTPMetrics) MustRegister(reg prometheus.Registerer) {
	reg.MustRegister(
		m.RequestsTotal,
		m.RequestDuration,
		m.InFlightRequests,
	)
}
