package v1

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"

	"github.com/happywbfriends/metrics/metrics"
)

func NewHttpClientMetrics(clientName, methodName string) HttpClientMetrics {
	return NewHttpClientRequestMetricsWithBuckets(clientName, methodName, metrics.DefaultDurationMsBuckets)
}

func NewHttpClientRequestMetricsWithBuckets(clientName, methodName string, requestTimeMsBuckets []float64) HttpClientMetrics {
	labels := map[string]string{
		metrics.MetricsLabelClient: clientName,
		metrics.MetricsLabelMethod: methodName,
	}

	m := &httpClientMetrics{
		nbDone:        metrics.NewCounterVec(metrics.MetricsNamespace, metrics.MetricsSubsystemHttpClt, "nb_req_done", labels, []string{metrics.MetricsLabelStatusCode}),
		nbError:       metrics.NewCounter(metrics.MetricsNamespace, metrics.MetricsSubsystemHttpClt, "nb_req_error", labels),
		requestTimeMs: metrics.NewHistogram(metrics.MetricsNamespace, metrics.MetricsSubsystemHttpClt, "req_duration_ms", labels, requestTimeMsBuckets),
	}

	m.nbDone200 = m.nbDone.WithLabelValues("200") // recommended optimization

	return m
}

type HttpClientMetrics interface {
	IncDone(statusCode int)
	IncError(e error)
	RequestDuration(duration time.Duration)
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
