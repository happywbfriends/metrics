package v1

import (
	"context"
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"sync/atomic"
	"time"

	"github.com/happywbfriends/metrics/metrics"
)

type DbMetrics interface {
	Update(stats sql.DBStats)
}

func NewDbMetrics(dbName string) DbMetrics {
	labels := map[string]string{
		MetricsLabelSubject: dbName,
	}

	return &dbMetrics{
		nbMaxConns:        metrics.NewGauge(metrics.MetricsNamespace, metrics.MetricsSubsystemDb, "max_conns", labels),
		nbOpenConns:       metrics.NewGauge(metrics.MetricsNamespace, metrics.MetricsSubsystemDb, "open_conns", labels),
		nbUsedConns:       metrics.NewGauge(metrics.MetricsNamespace, metrics.MetricsSubsystemDb, "used_conns", labels),
		waitCount:         metrics.NewGauge(metrics.MetricsNamespace, metrics.MetricsSubsystemDb, "wait_count", labels),
		waitDurationMs:    metrics.NewSummary(metrics.MetricsNamespace, metrics.MetricsSubsystemDb, "wait_duration_count", labels),
		maxIdleClosed:     metrics.NewCounter(metrics.MetricsNamespace, metrics.MetricsSubsystemDb, "max_idle_closed", labels),
		maxIdleTimeClosed: metrics.NewCounter(metrics.MetricsNamespace, metrics.MetricsSubsystemDb, "max_idle_time_closed", labels),
		maxLifetimeClosed: metrics.NewCounter(metrics.MetricsNamespace, metrics.MetricsSubsystemDb, "max_lifetime_closed", labels),
	}
}

type dbMetrics struct {
	nbMaxConns     prometheus.Gauge
	nbOpenConns    prometheus.Gauge
	nbUsedConns    prometheus.Gauge
	waitCount      prometheus.Gauge
	waitDurationMs prometheus.Summary

	maxIdleClosed      prometheus.Counter
	_maxIdleClosed     int64
	maxIdleTimeClosed  prometheus.Counter
	_maxIdleTimeClosed int64
	maxLifetimeClosed  prometheus.Counter
	_maxLifetimeClosed int64
}

func (d *dbMetrics) Update(stats sql.DBStats) {
	d.nbMaxConns.Set(float64(stats.MaxOpenConnections))
	d.nbOpenConns.Set(float64(stats.OpenConnections))
	d.nbUsedConns.Set(float64(stats.InUse))
	d.waitCount.Set(float64(stats.WaitCount))
	d.waitDurationMs.Observe(float64(stats.WaitDuration.Milliseconds()))

	if oldValue := atomic.SwapInt64(&d._maxIdleClosed, stats.MaxIdleClosed); oldValue < stats.MaxIdleClosed {
		d.maxIdleClosed.Add(float64(stats.MaxIdleClosed - oldValue))
	}

	if oldValue := atomic.SwapInt64(&d._maxIdleTimeClosed, stats.MaxIdleTimeClosed); oldValue < stats.MaxIdleTimeClosed {
		d.maxIdleTimeClosed.Add(float64(stats.MaxIdleTimeClosed - oldValue))
	}

	if oldValue := atomic.SwapInt64(&d._maxLifetimeClosed, stats.MaxLifetimeClosed); oldValue < stats.MaxLifetimeClosed {
		d.maxLifetimeClosed.Add(float64(stats.MaxLifetimeClosed - oldValue))
	}
}

func DbMetricsHelper(m DbMetrics, db *sql.DB, updateFreq time.Duration, ctx context.Context) {
	ticker := time.NewTicker(updateFreq)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
		// Безопасно для закрытой БД
		m.Update(db.Stats())
	}
}
