package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

const (
	metricsSubsystemHttpClt = "http_clt"
)

type IHttpClientMetrics interface {
	IncDone(statusCode int)
	IncError(e error)
	Measure(duration time.Duration)
}

type NoHttpClientMetrics struct{}

func (m *NoHttpClientMetrics) IncDone(int)           {}
func (m *NoHttpClientMetrics) IncError(error)        {}
func (m *NoHttpClientMetrics) Measure(time.Duration) {}

func NewHttpClientMetrics(clientName, methodName string, requestTimeMsBuckets []float64) IHttpClientMetrics {
	labels := map[string]string{
		metricsLabelClient: clientName,
		metricsLabelMethod: methodName,
	}

	m := &httpClientMetrics{
		nbDone:        newCounterVec(metricsNamespace, metricsSubsystemHttpClt, "nb_success", labels, []string{metricsLabelStatusCode}),
		nbError:       newCounter(metricsNamespace, metricsSubsystemHttpClt, "nb_error", labels),
		requestTimeMs: newHistogram(metricsNamespace, metricsSubsystemHttpClt, "req_duration_ms", labels, requestTimeMsBuckets),
	}

	m.nbDone200 = m.nbDone.WithLabelValues("200") // recommended optimization

	return m
}

type httpClientMetrics struct {
	nbDone        *prometheus.CounterVec
	nbDone200     prometheus.Counter // cached for optimization
	nbError       prometheus.Counter
	requestTimeMs prometheus.Histogram
}

func (m *httpClientMetrics) IncDone(statusCode int) {
	if statusCode == 200 {
		m.nbDone200.Inc()
	} else {
		m.nbDone.WithLabelValues(strconv.Itoa(statusCode)).Inc()
	}
}

func (m *httpClientMetrics) IncError(error) {
	m.nbError.Inc()
}

func (m *httpClientMetrics) Measure(t time.Duration) {
	m.requestTimeMs.Observe(float64(t.Milliseconds()))
}
