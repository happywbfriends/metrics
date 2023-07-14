package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type IHttpServerRequestMetrics interface {
	IncNbRequest(statusCode int)
	RequestDuration(duration time.Duration)
}

type NoHttpServerRequestMetrics struct{}

func (m *NoHttpServerRequestMetrics) IncNbRequest(int)              {}
func (m *NoHttpServerRequestMetrics) RequestDuration(time.Duration) {}

func NewHttpServerRequestMetrics(methodName string) IHttpServerRequestMetrics {
	return NewHttpServerRequestMetricsWithBuckets(methodName, DefaultDurationMsBuckets)
}

func NewHttpServerRequestMetricsWithBuckets(methodName string, requestTimeMsBuckets []float64) IHttpServerRequestMetrics {
	labels := map[string]string{
		metricsLabelMethod: methodName,
	}
	m := &httpServerRequestMetrics{
		nbRequests:    newCounterVec(metricsNamespace, metricsSubsystemHttpServer, "nb_req", labels, []string{metricsLabelStatusCode}),
		requestTimeMs: newHistogram(metricsNamespace, metricsSubsystemHttpServer, "req_duration_ms", labels, requestTimeMsBuckets),
	}
	m.nbRequests200 = m.nbRequests.WithLabelValues("200") // recommended optimization
	return m
}

func HttpServerRequestHelper(m IHttpServerRequestMetrics, status int, since time.Time) {
	m.IncNbRequest(status)
	m.RequestDuration(time.Since(since))
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

func (m *httpServerRequestMetrics) RequestDuration(t time.Duration) {
	m.requestTimeMs.Observe(float64(t.Milliseconds()))
}
