package metrics

import (
	"context"
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
	"sync/atomic"
	"time"
)

const metricsSubsystemDb = "db"

type IDbMetrics interface {
	Update(stats sql.DBStats)
}

func NewDbMetrics(dbName string) IDbMetrics {
	labels := map[string]string{
		metricsLabelDatabase: dbName,
	}

	return &dbMetrics{
		nbMaxConns:        newGauge(metricsNamespace, metricsSubsystemDb, "max_conns", labels),
		nbOpenConns:       newGauge(metricsNamespace, metricsSubsystemDb, "open_conns", labels),
		nbUsedConns:       newGauge(metricsNamespace, metricsSubsystemDb, "open_conns", labels),
		waitCount:         newGauge(metricsNamespace, metricsSubsystemDb, "wait_count", labels),
		waitDurationMs:    newSummary(metricsNamespace, metricsSubsystemDb, "wait_count", labels),
		maxIdleClosed:     newCounter(metricsNamespace, metricsSubsystemDb, "max_idle_closed", labels),
		maxIdleTimeClosed: newCounter(metricsNamespace, metricsSubsystemDb, "max_idle_time_closed", labels),
		maxLifetimeClosed: newCounter(metricsNamespace, metricsSubsystemDb, "max_lifetime_closed", labels),
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

func DbMetricsHelper(m IDbMetrics, db *sql.DB, updateFreq time.Duration, ctx context.Context) {
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
