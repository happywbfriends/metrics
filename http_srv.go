package metrics

import "github.com/prometheus/client_golang/prometheus"

const (
	metricsSubsystemHttpServer = "http_srv"
)

type IHttpServerMetrics interface {
	IncNbConnections()
	DecNbConnections()
}

type NoHttpServerMetrics struct{}

func (m *NoHttpServerMetrics) IncNbConnections() {}
func (m *NoHttpServerMetrics) DecNbConnections() {}

func NewHttpServerMetrics() IHttpServerMetrics {
	m := &httpServerMetrics{
		nbConnections: newGauge(metricsNamespace, metricsSubsystemHttpServer, "nb_current_conns", nil),
	}
	return m
}

type httpServerMetrics struct {
	nbConnections prometheus.Gauge
}

func (m *httpServerMetrics) IncNbConnections() {
	m.nbConnections.Inc()
}
func (m *httpServerMetrics) DecNbConnections() {
	m.nbConnections.Dec()
}
