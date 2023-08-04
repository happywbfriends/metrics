package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

const (
	metricsSubsystemHttpServer = "http_srv"
)

type IHttpServerMetrics interface {
	IncNbConnections()
	DecNbConnections()
	IncNotFound(path string, supplierOldId int)
	IncMethodNotAllowed(method, path string, supplierOldId int)
}

type NoHttpServerMetrics struct{}

func (m *NoHttpServerMetrics) IncNbConnections()                       {}
func (m *NoHttpServerMetrics) DecNbConnections()                       {}
func (m *NoHttpServerMetrics) IncNotFound(string, int)                 {}
func (m *NoHttpServerMetrics) IncMethodNotAllowed(string, string, int) {}

func NewHttpServerMetrics() IHttpServerMetrics {
	labels := map[string]string{
		metricsLabelMethod: "",
	}

	m := &httpServerMetrics{
		nbConnections: newGauge(metricsNamespace, metricsSubsystemHttpServer, "current_conns", nil),
		nbRequests:    newCounterVec(metricsNamespace, metricsSubsystemHttpServer, "nb_req", labels, []string{metricsLabelStatusCode, metricsLabelSupplierOldId}),
	}
	return m
}

type httpServerMetrics struct {
	nbConnections prometheus.Gauge
	nbRequests    *prometheus.CounterVec
}

func (m *httpServerMetrics) IncNbConnections() {
	m.nbConnections.Inc()
}

func (m *httpServerMetrics) DecNbConnections() {
	m.nbConnections.Dec()
}

func (m *httpServerMetrics) IncNotFound(_ string, supplierOldId int) {
	m.nbRequests.WithLabelValues("404", strconv.Itoa(supplierOldId)).Inc()
}

func (m *httpServerMetrics) IncMethodNotAllowed(_, _ string, supplierOldId int) {
	m.nbRequests.WithLabelValues("405", strconv.Itoa(supplierOldId)).Inc()
}
