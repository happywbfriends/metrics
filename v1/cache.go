package v1

import (
	"github.com/happywbfriends/metrics/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

func MewCacheMetrics() CacheMetrics {
	buckets := metrics.DefaultDurationMsBuckets // todo
	m := cacheMetrics{
		nbRead:       metrics.NewCounterVec(metrics.MetricsNamespace, MetricsSubsystemCache, "nb_read", nil, []string{MetricsLabelName, MetricsLabelShard, MetricsLabelHit}),
		nbWrite:      metrics.NewCounterVec(metrics.MetricsNamespace, MetricsSubsystemCache, "nb_write", nil, []string{MetricsLabelName, MetricsLabelShard}),
		readDuration: metrics.NewHistogramVec(metrics.MetricsNamespace, MetricsSubsystemCache, "read_duration_ms", nil, buckets, []string{MetricsLabelName, MetricsLabelShard, MetricsLabelHit}),
		size:         nil,
		maxSize:      nil,
	}

	return &m
}

type CacheMetrics interface {
	IncNbReadHit(name string, shard int)
	IncNbReadMiss(name string, shard int)
	ObserveReadHitDuration(name string, shard int, t time.Duration)
	ObserveReadMissDuration(name string, shard int, t time.Duration)
	IncNbWrite(name string, shard int)
}

type cacheMetrics struct {
	nbRead       *prometheus.CounterVec
	nbWrite      *prometheus.CounterVec
	readDuration *prometheus.HistogramVec
	// todo - в каком методе обновлять?
	size    *prometheus.Gauge
	maxSize *prometheus.Gauge
}

func (m *cacheMetrics) IncNbReadHit(name string, shard int) {
	m.incNbRead(name, shard, "1")
}

func (m *cacheMetrics) IncNbReadMiss(name string, shard int) {
	m.incNbRead(name, shard, "0")
}

func (m *cacheMetrics) incNbRead(name string, shard int, hit string) {
	m.nbRead.WithLabelValues(name, strconv.Itoa(shard), hit).Inc()
}

func (m *cacheMetrics) IncNbWrite(name string, shard int) {
	m.nbWrite.WithLabelValues(name, strconv.Itoa(shard)).Inc()
}

func (m *cacheMetrics) ObserveReadHitDuration(name string, shard int, t time.Duration) {
	m.observeReadDuration(name, shard, "1", t)
}

func (m *cacheMetrics) ObserveReadMissDuration(name string, shard int, t time.Duration) {
	m.observeReadDuration(name, shard, "0", t)
}

func (m *cacheMetrics) observeReadDuration(name string, shard int, hit string, t time.Duration) {
	m.readDuration.WithLabelValues(name, strconv.Itoa(shard), hit).Observe(float64(t.Milliseconds()))
}
