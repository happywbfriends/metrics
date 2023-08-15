package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

const (
	MetricsSubsystemHttpClt = "http_clt"
)

type IHttpClientRequestMetrics interface {
	IncDone(statusCode int)
	IncError(e error)
	RequestDuration(duration time.Duration)
}

type NoHttpClientMetrics struct{}

func (m *NoHttpClientMetrics) IncDone(int)                   {}
func (m *NoHttpClientMetrics) IncError(error)                {}
func (m *NoHttpClientMetrics) RequestDuration(time.Duration) {}

func NewHttpClientRequestMetrics(clientName, methodName string) IHttpClientRequestMetrics {
	return NewHttpClientRequestMetricsWithBuckets(clientName, methodName, DefaultDurationMsBuckets)
}

func NewHttpClientRequestMetricsWithBuckets(clientName, methodName string, requestTimeMsBuckets []float64) IHttpClientRequestMetrics {
	labels := map[string]string{
		MetricsLabelClient: clientName,
		MetricsLabelMethod: methodName,
	}

	m := &httpClientMetrics{
		nbDone:        NewCounterVec(MetricsNamespace, MetricsSubsystemHttpClt, "nb_req_done", labels, []string{MetricsLabelStatusCode}),
		nbError:       NewCounter(MetricsNamespace, MetricsSubsystemHttpClt, "nb_req_error", labels),
		requestTimeMs: NewHistogram(MetricsNamespace, MetricsSubsystemHttpClt, "req_duration_ms", labels, requestTimeMsBuckets),
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

func (m *httpClientMetrics) RequestDuration(t time.Duration) {
	m.requestTimeMs.Observe(float64(t.Milliseconds()))
}
