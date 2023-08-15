package v1

import (
	"github.com/happywbfriends/metrics/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

func NewDbQueryMetrics() DbQueryMetrics {
	return &dbRequestMetrics{
		durationMs: metrics.NewHistogramVec(metrics.MetricsNamespace, metrics.MetricsSubsystemDbQuery, "duration_ms", nil, metrics.DefaultDurationMsBuckets, []string{MetricsLabelSubject, metrics.MetricsLabelDatabaseQuery}),
		nbDone:     metrics.NewCounterVec(metrics.MetricsNamespace, metrics.MetricsSubsystemDbQuery, "nb_done", nil, []string{MetricsLabelSubject, metrics.MetricsLabelDatabaseQuery}),
		nbError:    metrics.NewCounterVec(metrics.MetricsNamespace, metrics.MetricsSubsystemDbQuery, "nb_error", nil, []string{MetricsLabelSubject, metrics.MetricsLabelDatabaseQuery}),
	}
}

type DbQueryMetrics interface {
	IncNbDone(subject string, query string)
	IncNbError(subject string, query string)
	ObserveRequestDuration(subject string, query string, duration time.Duration)
}

type dbRequestMetrics struct {
	durationMs *prometheus.HistogramVec
	nbDone     *prometheus.CounterVec
	nbError    *prometheus.CounterVec
}

func (m *dbRequestMetrics) IncNbDone(subject string, query string) {
	m.nbDone.WithLabelValues(subject, query).Inc()
}
func (m *dbRequestMetrics) IncNbError(subject string, query string) {
	m.nbError.WithLabelValues(subject, query).Inc()
}
func (m *dbRequestMetrics) ObserveRequestDuration(subject string, query string, duration time.Duration) {
	m.durationMs.WithLabelValues(subject, query).Observe(float64(duration.Milliseconds()))
}
