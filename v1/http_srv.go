package v1

import (
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/happywbfriends/metrics/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

func NewHTTPServerMetrics() HTTPServerMetrics {
	return NewHTTPServerMetricsWithBuckets(metrics.DefaultDurationMsBuckets)
}

func NewHTTPServerMetricsWithBuckets(requestTimeMsBuckets []float64) HTTPServerMetrics {
	m := &httpServerMetrics{
		nbRequests:    metrics.NewCounterVec(metrics.MetricsNamespace, metrics.MetricsSubsystemHttpServer, "nb_req", nil, []string{metrics.MetricsLabelMethod, metrics.MetricsLabelStatusCode, MetricsLabelClientId}),
		requestTimeMs: metrics.NewHistogramVec(metrics.MetricsNamespace, metrics.MetricsSubsystemHttpServer, "req_duration_ms", nil, requestTimeMsBuckets, []string{metrics.MetricsLabelMethod, metrics.MetricsLabelStatusCode, MetricsLabelClientId}),
		nbConnections: metrics.NewGauge(metrics.MetricsNamespace, metrics.MetricsSubsystemHttpServer, "current_conns", nil),
	}
	return m
}

type HTTPServerMetrics interface {
	IncNbRequest(method string, statusCode int, clientId int)
	ObserveRequestDuration(method string, statusCode int, clientId int, duration time.Duration)
	IncNbConnections()
	DecNbConnections()
	OnStateChange(conn net.Conn, state http.ConnState)
}

type NoHTTPServerMetrics struct{}

func (m *NoHTTPServerMetrics) IncNbRequest(method string, statusCode int, clientId int) {}
func (m *NoHTTPServerMetrics) ObserveRequestDuration(method string, statusCode int, clientId int, duration time.Duration) {
}
func (m *NoHTTPServerMetrics) IncNbConnections()                          {}
func (m *NoHTTPServerMetrics) DecNbConnections()                          {}
func (m *NoHTTPServerMetrics) OnStateChange(_ net.Conn, _ http.ConnState) {}

type httpServerMetrics struct {
	nbRequests    *prometheus.CounterVec
	requestTimeMs *prometheus.HistogramVec
	nbConnections prometheus.Gauge
}

func (m *httpServerMetrics) IncNbRequest(method string, statusCode int, clientId int) {
	m.nbRequests.WithLabelValues(method, strconv.Itoa(statusCode), strconv.Itoa(clientId)).Inc()
}

func (m *httpServerMetrics) ObserveRequestDuration(method string, statusCode int, clientId int, duration time.Duration) {
	if statusCode == http.StatusOK || statusCode == http.StatusInternalServerError {
		m.requestTimeMs.WithLabelValues(method, strconv.Itoa(statusCode), strconv.Itoa(clientId)).Observe(float64(duration.Milliseconds()))
	}
}

// IncNbConnections увеличивает количество активных соединений
// Использовать или в явном виде, или через хук OnStateChange
func (m *httpServerMetrics) IncNbConnections() {
	m.nbConnections.Inc()
}

// DecNbConnections уменьшает количество активных соединений
// Использовать или в явном виде, или через хук OnStateChange
func (m *httpServerMetrics) DecNbConnections() {
	m.nbConnections.Dec()
}

// OnStateChange это хук для обновления количества активных соединений
// При использовании хука отпадает необходимость дергать IncNbConnections и DecNbConnections в явном виде
func (m *httpServerMetrics) OnStateChange(_ net.Conn, state http.ConnState) {
	switch state {
	case http.StateNew:
		m.IncNbConnections()
	case http.StateHijacked, http.StateClosed:
		m.DecNbConnections()
	}
}
