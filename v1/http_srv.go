package v1

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"

	"github.com/happywbfriends/metrics/metrics"
)

func NewHTTPServerMetrics() HTTPServerMetrics {
	return NewHttpServerMetricsWithBuckets(metrics.DefaultDurationMsBuckets)
}

type HTTPServerMetrics interface {
	IncNbRequest(method string, statusCode int, supplierOldId int)
	ObserveOkRequestDuration(method string, duration time.Duration)
}

type httpServerMetrics struct {
	nbRequests    *prometheus.CounterVec
	requestTimeMs *prometheus.HistogramVec
}

func (m *httpServerMetrics) IncNbRequest(method string, statusCode int, supplierOldId int) {
	m.nbRequests.WithLabelValues(method, strconv.Itoa(statusCode), strconv.Itoa(supplierOldId)).Inc()
}

func (m *httpServerMetrics) ObserveOkRequestDuration(method string, duration time.Duration) {
	m.requestTimeMs.WithLabelValues(method).Observe(float64(duration.Milliseconds()))
}

func NewHttpServerMetricsWithBuckets(requestTimeMsBuckets []float64) HTTPServerMetrics {
	m := &httpServerMetrics{
		nbRequests:    metrics.NewCounterVec(metrics.MetricsNamespace, metrics.MetricsSubsystemHttpServer, "nb_req", nil, []string{metrics.MetricsLabelMethod, metrics.MetricsLabelStatusCode, metrics.MetricsLabelSupplierOldId}),
		requestTimeMs: metrics.NewHistogramVec(metrics.MetricsNamespace, metrics.MetricsSubsystemHttpServer, "req_duration_ms", nil, requestTimeMsBuckets, []string{metrics.MetricsLabelMethod}),
	}
	return m
}
