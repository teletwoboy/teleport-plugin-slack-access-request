package metric

import (
	"database/sql"
	"teleport-plugin-slack-access-request/internal/config"
	"teleport-plugin-slack-access-request/internal/database"

	"github.com/prometheus/client_golang/prometheus"
)

// DBStatsCollector implements prometheus.Collector
type DBStatsCollector struct {
	db                *sql.DB
	openConnections   *prometheus.Desc
	inUseConnections  *prometheus.Desc
	idleConnections   *prometheus.Desc
	waitCountTotal    *prometheus.Desc
	waitDurationTotal *prometheus.Desc
}

// NewDBStatsCollector returns a custom collector for sql.DB stats
func NewDBStatsCollector(db *database.DB) *DBStatsCollector {
	labels := prometheus.Labels{
		"db_name": config.Cfg.Database.Database,
	}

	return &DBStatsCollector{
		db: db.Conn,
		openConnections: prometheus.NewDesc(
			"db_open_connections",
			"Total number of established connections both in use and idle",
			nil, labels,
		),
		inUseConnections: prometheus.NewDesc(
			"db_in_use_connections",
			"Number of connections currently in use",
			nil, labels,
		),
		idleConnections: prometheus.NewDesc(
			"db_idle_connections",
			"Number of idle connections",
			nil, labels,
		),
		waitCountTotal: prometheus.NewDesc(
			"db_wait_count_total",
			"Total number of connections waited for",
			nil, labels,
		),
		waitDurationTotal: prometheus.NewDesc(
			"db_wait_duration_seconds_total",
			"Total time blocked waiting for a new connection (in seconds)",
			nil, labels,
		),
	}
}

// Describe sends metric descriptions to Prometheus
func (c *DBStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.openConnections
	ch <- c.inUseConnections
	ch <- c.idleConnections
	ch <- c.waitCountTotal
	ch <- c.waitDurationTotal
}

// Collect is called on each scrape
func (c *DBStatsCollector) Collect(ch chan<- prometheus.Metric) {
	stats := c.db.Stats()

	ch <- prometheus.MustNewConstMetric(
		c.openConnections, prometheus.GaugeValue, float64(stats.OpenConnections),
	)
	ch <- prometheus.MustNewConstMetric(
		c.inUseConnections, prometheus.GaugeValue, float64(stats.InUse),
	)
	ch <- prometheus.MustNewConstMetric(
		c.idleConnections, prometheus.GaugeValue, float64(stats.Idle),
	)
	ch <- prometheus.MustNewConstMetric(
		c.waitCountTotal, prometheus.CounterValue, float64(stats.WaitCount),
	)
	ch <- prometheus.MustNewConstMetric(
		c.waitDurationTotal, prometheus.CounterValue, stats.WaitDuration.Seconds(),
	)
}
