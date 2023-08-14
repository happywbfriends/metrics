package v1

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"

	"github.com/happywbfriends/metrics/metrics"
)

func NewHttpClientMetrics() HttpClientMetrics {
	return NewHttpClientRequestMetricsWithBuckets(metrics.DefaultDurationMsBuckets)
}

func NewHttpClientRequestMetricsWithBuckets(requestTimeMsBuckets []float64) HttpClientMetrics {
	m := &httpClientMetrics{
		nbDone:        metrics.NewCounterVec(metrics.MetricsNamespace, metrics.MetricsSubsystemHttpClt, "nb_req_done", nil, []string{metrics.MetricsLabelClient, metrics.MetricsLabelMethod, metrics.MetricsLabelStatusCode}),
		nbError:       metrics.NewCounterVec(metrics.MetricsNamespace, metrics.MetricsSubsystemHttpClt, "nb_req_error", nil, []string{metrics.MetricsLabelClient, metrics.MetricsLabelMethod}),
		requestTimeMs: metrics.NewHistogramVec(metrics.MetricsNamespace, metrics.MetricsSubsystemHttpClt, "req_duration_ms", nil, requestTimeMsBuckets, []string{metrics.MetricsLabelClient, metrics.MetricsLabelMethod}),
	}

	return m
}

type HttpClientMetrics interface {
	IncNbDone(client string, method string, statusCode int)
	IncNbError(client string, method string)
	RequestDuration(client string, method string, t time.Duration)
}

type httpClientMetrics struct {
	nbDone        *prometheus.CounterVec
	nbError       *prometheus.CounterVec
	requestTimeMs *prometheus.HistogramVec
}

func (m *httpClientMetrics) IncNbDone(client string, method string, statusCode int) {
	m.nbDone.WithLabelValues(client, method, strconv.Itoa(statusCode)).Inc()
}

func (m *httpClientMetrics) IncNbError(client string, method string) {
	m.nbError.WithLabelValues(client, method).Inc()
}

func (m *httpClientMetrics) RequestDuration(client string, method string, t time.Duration) {
	m.requestTimeMs.WithLabelValues(client, method).Observe(float64(t.Milliseconds()))
}
