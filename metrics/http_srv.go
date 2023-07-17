package metrics

import "github.com/prometheus/client_golang/prometheus"

const (
	metricsSubsystemHttpServer = "http_srv"
)

type IHttpServerMetrics interface {
	IncNbConnections()
	DecNbConnections()
	IncNotFound(path string)
	IncMethodNotAllowed(method, path string)
}

type NoHttpServerMetrics struct{}

func (m *NoHttpServerMetrics) IncNbConnections() {}
func (m *NoHttpServerMetrics) DecNbConnections() {}

func NewHttpServerMetrics() IHttpServerMetrics {
	m := &httpServerMetrics{
		nbConnections:      newGauge(metricsNamespace, metricsSubsystemHttpServer, "current_conns", nil),
		nbNotFound:         newGauge(metricsNamespace, metricsSubsystemHttpServer, "nb_req_not_found", nil),
		nbMethodNotAllowed: newGauge(metricsNamespace, metricsSubsystemHttpServer, "nb_req_not_allowed", nil),
	}
	return m
}

type httpServerMetrics struct {
	nbConnections      prometheus.Gauge
	nbNotFound         prometheus.Counter
	nbMethodNotAllowed prometheus.Counter
}

func (m *httpServerMetrics) IncNbConnections() {
	m.nbConnections.Inc()
}
func (m *httpServerMetrics) DecNbConnections() {
	m.nbConnections.Dec()
}

func (m *httpServerMetrics) IncNotFound(string) {
	m.nbNotFound.Inc()
}
func (m *httpServerMetrics) IncMethodNotAllowed(_, _ string) {
	m.nbMethodNotAllowed.Inc()
}
