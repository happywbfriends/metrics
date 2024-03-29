package v1

import (
	"strconv"
	"time"

	"github.com/happywbfriends/metrics/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

func NewHTTPClientMetrics() HTTPClientMetrics {
	return NewHTTPClientMetricsWithBuckets(metrics.DefaultDurationMsBuckets)
}

func NewHTTPClientMetricsWithBuckets(requestTimeMsBuckets []float64) HTTPClientMetrics {
	m := &httpClientMetrics{
		nbDone:        metrics.NewCounterVec(metrics.MetricsNamespace, metrics.MetricsSubsystemHttpClt, "nb_req_done", nil, []string{MetricsLabelSubject, metrics.MetricsLabelMethod, metrics.MetricsLabelStatusCode}),
		nbError:       metrics.NewCounterVec(metrics.MetricsNamespace, metrics.MetricsSubsystemHttpClt, "nb_req_error", nil, []string{MetricsLabelSubject, metrics.MetricsLabelMethod}),
		requestTimeMs: metrics.NewHistogramVec(metrics.MetricsNamespace, metrics.MetricsSubsystemHttpClt, "req_duration_ms", nil, requestTimeMsBuckets, []string{MetricsLabelSubject, metrics.MetricsLabelMethod}),
	}

	return m
}

type HTTPClientMetrics interface {
	IncNbDone(subject string, method string, statusCode int)
	IncNbError(subject string, method string)
	ObserveRequestDuration(subject string, method string, t time.Duration)
}

type NoHTTPClientMetrics struct{}

func (m *NoHTTPClientMetrics) IncNbDone(subject string, method string, statusCode int) {}
func (m *NoHTTPClientMetrics) IncNbError(subject string, method string)                {}
func (m *NoHTTPClientMetrics) ObserveRequestDuration(subject string, method string, t time.Duration) {
}

type httpClientMetrics struct {
	nbDone        *prometheus.CounterVec
	nbError       *prometheus.CounterVec
	requestTimeMs *prometheus.HistogramVec
}

func (m *httpClientMetrics) IncNbDone(subject string, method string, statusCode int) {
	m.nbDone.WithLabelValues(subject, method, strconv.Itoa(statusCode)).Inc()
}

func (m *httpClientMetrics) IncNbError(subject string, method string) {
	m.nbError.WithLabelValues(subject, method).Inc()
}

func (m *httpClientMetrics) ObserveRequestDuration(subject string, method string, t time.Duration) {
	m.requestTimeMs.WithLabelValues(subject, method).Observe(float64(t.Milliseconds()))
}
