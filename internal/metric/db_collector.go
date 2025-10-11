/*
Copyright 2025 steamedEggMaster

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package metric

import (
	"database/sql"

	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/config"
	"github.com/teletwoboy/teleport-plugin-slack-access-request/internal/database"

	"github.com/prometheus/client_golang/prometheus"
)

// DBStatsCollector implements prometheus.Collector
type DBStatsCollector struct {
	db                 *sql.DB
	maxOpenConnections *prometheus.Desc
	openConnections    *prometheus.Desc
	inUseConnections   *prometheus.Desc
	idleConnections    *prometheus.Desc
	waitCountTotal     *prometheus.Desc
	waitDurationTotal  *prometheus.Desc
	maxIdleClosed      *prometheus.Desc
	maxLifetimeClosed  *prometheus.Desc
}

// NewDBStatsCollector returns a custom collector for sql.DB stats
func NewDBStatsCollector(db *database.DB) *DBStatsCollector {
	labels := prometheus.Labels{
		"db_name": config.Cfg.Database.Database,
	}

	return &DBStatsCollector{
		db: db.Conn,

		maxOpenConnections: prometheus.NewDesc(
			"db_max_open_connections",
			"Maximum number of open connections allowed (SetMaxOpenConns)",
			nil, labels,
		),
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
		maxIdleClosed: prometheus.NewDesc(
			"db_max_idle_closed_total",
			"Total connections closed due to exceeding SetMaxIdleConns",
			nil, labels,
		),
		maxLifetimeClosed: prometheus.NewDesc(
			"db_max_lifetime_closed_total",
			"Total connections closed due to exceeding SetConnMaxLifetime",
			nil, labels,
		),
	}
}

// Describe sends metric descriptions to Prometheus
func (c *DBStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.maxOpenConnections
	ch <- c.openConnections
	ch <- c.inUseConnections
	ch <- c.idleConnections
	ch <- c.waitCountTotal
	ch <- c.waitDurationTotal
	ch <- c.maxIdleClosed
	ch <- c.maxLifetimeClosed
}

// Collect is called on each scrape
func (c *DBStatsCollector) Collect(ch chan<- prometheus.Metric) {
	stats := c.db.Stats()

	ch <- prometheus.MustNewConstMetric(
		c.maxOpenConnections, prometheus.GaugeValue, float64(stats.MaxOpenConnections),
	)
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
	ch <- prometheus.MustNewConstMetric(
		c.maxIdleClosed, prometheus.CounterValue, float64(stats.MaxIdleClosed),
	)
	ch <- prometheus.MustNewConstMetric(
		c.maxLifetimeClosed, prometheus.CounterValue, float64(stats.MaxLifetimeClosed),
	)
}
