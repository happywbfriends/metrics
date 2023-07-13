package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type IHttpServerRequestMetrics interface {
	IncNbRequest(statusCode int)
	RequestTime(duration time.Duration)
}

type NoHttpServerRequestMetrics struct{}

func (m *NoHttpServerRequestMetrics) IncNbRequest(int)          {}
func (m *NoHttpServerRequestMetrics) RequestTime(time.Duration) {}

func NewHttpServerRequestMetrics(methodName string, requestTimeMsBuckets []float64) IHttpServerRequestMetrics {
	labels := map[string]string{
		metricsLabelMethod: methodName,
	}
	m := &httpServerRequestMetrics{
		nbRequests:    newCounterVec(metricsNamespace, metricsSubsystemHttpServer, "nb_req", labels, []string{metricsLabelStatusCode}),
		requestTimeMs: newHistogram(metricsNamespace, metricsSubsystemHttpServer, "req_time_ms", labels, requestTimeMsBuckets),
	}
	m.nbRequests200 = m.nbRequests.WithLabelValues("200") // recommended optimization
	return m
}

type httpServerRequestMetrics struct {
	nbRequests    *prometheus.CounterVec
	nbRequests200 prometheus.Counter
	requestTimeMs prometheus.Histogram
}

func (m *httpServerRequestMetrics) IncNbRequest(statusCode int) {
	if statusCode == 200 {
		m.nbRequests200.Inc()
	} else {
		m.nbRequests.WithLabelValues(strconv.Itoa(statusCode)).Inc()
	}
}

func (m *httpServerRequestMetrics) RequestTime(t time.Duration) {
	m.requestTimeMs.Observe(float64(t.Milliseconds()))
}
