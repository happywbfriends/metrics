package v1

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
	"time"

	"github.com/happywbfriends/metrics/metrics"
)

func NewHTTPServerMetrics() HTTPServerMetrics {
	return NewHTTPServerMetricsWithBuckets(metrics.DefaultDurationMsBuckets)
}

func NewHTTPServerMetricsWithBuckets(requestTimeMsBuckets []float64) HTTPServerMetrics {
	m := &httpServerMetrics{
		nbRequests:    metrics.NewCounterVec(metrics.MetricsNamespace, metrics.MetricsSubsystemHttpServer, "nb_req", nil, []string{metrics.MetricsLabelMethod, metrics.MetricsLabelStatusCode, metrics.MetricsLabelSupplierOldId}),
		requestTimeMs: metrics.NewHistogramVec(metrics.MetricsNamespace, metrics.MetricsSubsystemHttpServer, "req_duration_ms", nil, requestTimeMsBuckets, []string{metrics.MetricsLabelMethod, metrics.MetricsLabelStatusCode, metrics.MetricsLabelSupplierOldId}),
		nbConnections: metrics.NewGauge(metrics.MetricsNamespace, metrics.MetricsSubsystemHttpServer, "current_conns", nil),
	}
	return m
}

type HTTPServerMetrics interface {
	IncNbRequest(method string, statusCode int, supplierOldId int)
	ObserveRequestDuration(method string, statusCode int, supplierOldId int, duration time.Duration)
	IncNbConnections()
	DecNbConnections()
}

type NoHTTPServerMetrics struct{}

func (m *NoHTTPServerMetrics) IncNbRequest(method string, statusCode int, supplierOldId int) {}
func (m *NoHTTPServerMetrics) ObserveRequestDuration(method string, statusCode int, supplierOldId int, duration time.Duration) {
}
func (m *NoHTTPServerMetrics) IncNbConnections() {}
func (m *NoHTTPServerMetrics) DecNbConnections() {}

type httpServerMetrics struct {
	nbRequests    *prometheus.CounterVec
	requestTimeMs *prometheus.HistogramVec
	nbConnections prometheus.Gauge
}

func (m *httpServerMetrics) IncNbRequest(method string, statusCode int, supplierOldId int) {
	m.nbRequests.WithLabelValues(method, strconv.Itoa(statusCode), strconv.Itoa(supplierOldId)).Inc()
}

func (m *httpServerMetrics) ObserveRequestDuration(method string, statusCode int, supplierOldId int, duration time.Duration) {
	if statusCode == http.StatusOK || statusCode == http.StatusInternalServerError {
		m.requestTimeMs.WithLabelValues(method, strconv.Itoa(statusCode), strconv.Itoa(supplierOldId)).Observe(float64(duration.Milliseconds()))
	}
}

func (m *httpServerMetrics) IncNbConnections() {
	m.nbConnections.Inc()
}

func (m *httpServerMetrics) DecNbConnections() {
	m.nbConnections.Dec()
}
