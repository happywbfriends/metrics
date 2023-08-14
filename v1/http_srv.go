package v1

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"

	"github.com/happywbfriends/metrics/metrics"
)

func NewHTTPServerRequestMetrics() HTTPServerRequestMetrics {
	return NewHttpServerRequestMetricsWithBuckets(metrics.DefaultDurationMsBuckets)
}

type HTTPServerRequestMetrics interface {
	IncNbRequest(method string, statusCode int, supplierOldId int)
	ObserveRequestDuration(method string, duration time.Duration)
}

type httpServerRequestMetrics struct {
	nbRequests    *prometheus.CounterVec
	requestTimeMs *prometheus.HistogramVec
}

func (m *httpServerRequestMetrics) IncNbRequest(method string, statusCode int, supplierOldId int) {
	m.nbRequests.WithLabelValues(method, strconv.Itoa(statusCode), strconv.Itoa(supplierOldId)).Inc()
}

func (m *httpServerRequestMetrics) ObserveRequestDuration(method string, duration time.Duration) {
	m.requestTimeMs.WithLabelValues(method).Observe(float64(duration.Milliseconds()))
}

func NewHttpServerRequestMetricsWithBuckets(requestTimeMsBuckets []float64) HTTPServerRequestMetrics {
	m := &httpServerRequestMetrics{
		nbRequests:    metrics.NewCounterVec(metrics.MetricsNamespace, metrics.MetricsSubsystemHttpServer, "nb_req", nil, []string{metrics.MetricsLabelMethod, metrics.MetricsLabelStatusCode, metrics.MetricsLabelSupplierOldId}),
		requestTimeMs: metrics.NewHistogramVec(metrics.MetricsNamespace, metrics.MetricsSubsystemHttpServer, "req_duration_ms", nil, requestTimeMsBuckets, []string{metrics.MetricsLabelMethod}),
	}
	return m
}
