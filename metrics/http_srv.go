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

func (m *NoHttpServerMetrics) IncNbConnections()                  {}
func (m *NoHttpServerMetrics) DecNbConnections()                  {}
func (m *NoHttpServerMetrics) IncNotFound(string)                 {}
func (m *NoHttpServerMetrics) IncMethodNotAllowed(string, string) {}

func NewHttpServerMetrics() IHttpServerMetrics {
	labels := map[string]string{
		metricsLabelMethod: "",
	}

	m := &httpServerMetrics{
		nbConnections: newGauge(metricsNamespace, metricsSubsystemHttpServer, "current_conns", nil),
		nbRequests:    newCounterVec(metricsNamespace, metricsSubsystemHttpServer, "nb_req", labels, []string{metricsLabelStatusCode}),
	}
	m.nbRequests404 = m.nbRequests.WithLabelValues("404")
	m.nbRequests405 = m.nbRequests.WithLabelValues("405")
	return m
}

type httpServerMetrics struct {
	nbConnections prometheus.Gauge
	nbRequests    *prometheus.CounterVec
	nbRequests404 prometheus.Counter
	nbRequests405 prometheus.Counter
}

func (m *httpServerMetrics) IncNbConnections() {
	m.nbConnections.Inc()
}
func (m *httpServerMetrics) DecNbConnections() {
	m.nbConnections.Dec()
}

func (m *httpServerMetrics) IncNotFound(string) {
	m.nbRequests404.Inc()
}
func (m *httpServerMetrics) IncMethodNotAllowed(_, _ string) {
	m.nbRequests405.Inc()
}
