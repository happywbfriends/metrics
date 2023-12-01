package v1

import (
	"github.com/happywbfriends/metrics/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

var cacheBuckets = []float64{1.0, 5.0, 10.0, 25.0, 50.0, 100.0, 250.0, 500.0}

func NewCacheMetrics() CacheMetrics {
	m := cacheMetrics{
		nbRead:       metrics.NewCounterVec(metrics.MetricsNamespace, MetricsSubsystemCache, "nb_read", nil, []string{MetricsLabelName, MetricsLabelShard, MetricsLabelHit}),
		nbWrite:      metrics.NewCounterVec(metrics.MetricsNamespace, MetricsSubsystemCache, "nb_write", nil, []string{MetricsLabelName, MetricsLabelShard}),
		readDuration: metrics.NewHistogramVec(metrics.MetricsNamespace, MetricsSubsystemCache, "read_duration_ns", nil, cacheBuckets, []string{MetricsLabelName, MetricsLabelShard, MetricsLabelHit}),
		size:         metrics.NewGaugeVec(metrics.MetricsNamespace, MetricsSubsystemCache, "size", nil, []string{MetricsLabelName, MetricsLabelShard}),
		maxSize:      metrics.NewGaugeVec(metrics.MetricsNamespace, MetricsSubsystemCache, "max_size", nil, []string{MetricsLabelName, MetricsLabelShard}),
	}

	return &m
}

type CacheMetrics interface {
	IncNbReadHit(name string, shard int)
	IncNbReadMiss(name string, shard int)
	ObserveReadHitDuration(name string, shard int, t time.Duration)
	ObserveReadMissDuration(name string, shard int, t time.Duration)
	IncNbWrite(name string, shard int)
	SetSize(name string, shard int, size int)
	IncSize(name string, shard int)
	DecSize(name string, shard int)
	SetMaxSize(name string, shard int, size int)
}

type cacheMetrics struct {
	nbRead       *prometheus.CounterVec
	nbWrite      *prometheus.CounterVec
	readDuration *prometheus.HistogramVec
	size         *prometheus.GaugeVec
	maxSize      *prometheus.GaugeVec
}

func (m *cacheMetrics) IncNbReadHit(name string, shard int) {
	m.incNbRead(name, shard, "1")
}

func (m *cacheMetrics) IncNbReadMiss(name string, shard int) {
	m.incNbRead(name, shard, "0")
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

func (m *cacheMetrics) SetSize(name string, shard int, size int) {
	m.size.WithLabelValues(name, strconv.Itoa(shard)).Set(float64(size))
}

func (m *cacheMetrics) IncSize(name string, shard int) {
	m.size.WithLabelValues(name, strconv.Itoa(shard)).Inc()
}

func (m *cacheMetrics) DecSize(name string, shard int) {
	m.size.WithLabelValues(name, strconv.Itoa(shard)).Dec()
}

func (m *cacheMetrics) SetMaxSize(name string, shard int, size int) {
	m.maxSize.WithLabelValues(name, strconv.Itoa(shard)).Set(float64(size))
}

func (m *cacheMetrics) incNbRead(name string, shard int, hit string) {
	m.nbRead.WithLabelValues(name, strconv.Itoa(shard), hit).Inc()
}

func (m *cacheMetrics) observeReadDuration(name string, shard int, hit string, t time.Duration) {
	m.readDuration.WithLabelValues(name, strconv.Itoa(shard), hit).Observe(float64(t.Nanoseconds()))
}
