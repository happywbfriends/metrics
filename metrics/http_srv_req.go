package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

const NoSupplierOldId = 0

type IHttpServerRequestMetrics interface {
	IncNbRequest(statusCode int, supplierOldId int)
	RequestDuration(duration time.Duration, supplierOldId int)
}

type NoHttpServerRequestMetrics struct{}

func (m *NoHttpServerRequestMetrics) IncNbRequest(int, int)              {}
func (m *NoHttpServerRequestMetrics) RequestDuration(time.Duration, int) {}

func NewHttpServerRequestMetrics(methodName string) IHttpServerRequestMetrics {
	return NewHttpServerRequestMetricsWithBuckets(methodName, DefaultDurationMsBuckets)
}

func NewHttpServerRequestMetricsWithBuckets(methodName string, requestTimeMsBuckets []float64) IHttpServerRequestMetrics {
	labels := map[string]string{
		metricsLabelMethod: methodName,
	}
	m := &httpServerRequestMetrics{
		nbRequests:    newCounterVec(metricsNamespace, metricsSubsystemHttpServer, "nb_req", labels, []string{metricsLabelStatusCode, metricsLabelSupplierOldId}),
		requestTimeMs: newHistogram(metricsNamespace, metricsSubsystemHttpServer, "req_duration_ms", labels, requestTimeMsBuckets),
	}
	return m
}

func HttpServerRequestHelper(m IHttpServerRequestMetrics, status int, since time.Time, supplierOldId int) {
	m.IncNbRequest(status, supplierOldId)
	m.RequestDuration(time.Since(since), supplierOldId)
}

type httpServerRequestMetrics struct {
	nbRequests    *prometheus.CounterVec
	requestTimeMs prometheus.Histogram
}

func (m *httpServerRequestMetrics) IncNbRequest(statusCode int, supplierOldId int) {
	statusCodeStr := "200"
	if statusCode != 200 {
		statusCodeStr = strconv.Itoa(statusCode)
	}

	supplierOldIdStr := "0"
	if supplierOldId != NoSupplierOldId {
		supplierOldIdStr = strconv.Itoa(supplierOldId)
	}

	m.nbRequests.WithLabelValues(statusCodeStr, supplierOldIdStr).Inc()
}

func (m *httpServerRequestMetrics) RequestDuration(t time.Duration, _ int) {
	m.requestTimeMs.Observe(float64(t.Milliseconds()))
}
