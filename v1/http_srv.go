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

func NewHttpServerMetricsWithBuckets(requestTimeMsBuckets []float64) HTTPServerMetrics {
	m := &httpServerMetrics{
		nbRequests:    metrics.NewCounterVec(metrics.MetricsNamespace, metrics.MetricsSubsystemHttpServer, "nb_req", nil, []string{metrics.MetricsLabelMethod, metrics.MetricsLabelStatusCode, metrics.MetricsLabelSupplierOldId}),
		requestTimeMs: metrics.NewHistogramVec(metrics.MetricsNamespace, metrics.MetricsSubsystemHttpServer, "req_duration_ms", nil, requestTimeMsBuckets, []string{metrics.MetricsLabelMethod}),
		nbConnections: metrics.NewGauge(metrics.MetricsNamespace, metrics.MetricsSubsystemHttpServer, "current_conns", nil),
	}
	return m
}

type HTTPServerMetrics interface {
	IncNbRequest(method string, statusCode int, supplierOldId int)
	ObserveOkRequestDuration(method string, duration time.Duration)
	IncNbConnections()
	DecNbConnections()
}

type httpServerMetrics struct {
	nbRequests    *prometheus.CounterVec
	requestTimeMs *prometheus.HistogramVec
	nbConnections prometheus.Gauge
}

func (m *httpServerMetrics) IncNbRequest(method string, statusCode int, supplierOldId int) {
	m.nbRequests.WithLabelValues(method, strconv.Itoa(statusCode), strconv.Itoa(supplierOldId)).Inc()
}

func (m *httpServerMetrics) ObserveOkRequestDuration(method string, duration time.Duration) {
	m.requestTimeMs.WithLabelValues(method).Observe(float64(duration.Milliseconds()))
}

func (m *httpServerMetrics) IncNbConnections() {
	m.nbConnections.Inc()
}

func (m *httpServerMetrics) DecNbConnections() {
	m.nbConnections.Dec()
}
