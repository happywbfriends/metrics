package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type HTTPServerRequestMetricsV1 interface {
	IncNbRequest(method string, statusCode int, supplierOldId int)
	ObserveRequestDuration(method string, duration time.Duration)
}

type httpServerRequestMetricsV1 struct {
	nbRequests    *prometheus.CounterVec
	requestTimeMs *prometheus.HistogramVec
}

func (m *httpServerRequestMetricsV1) IncNbRequest(method string, statusCode int, supplierOldId int) {
	m.nbRequests.WithLabelValues(method, strconv.Itoa(statusCode), strconv.Itoa(supplierOldId)).Inc()
}

func (m *httpServerRequestMetricsV1) ObserveRequestDuration(method string, duration time.Duration) {
	m.requestTimeMs.WithLabelValues(method).Observe(float64(duration.Milliseconds()))
}

func NewHTTPServerRequestMetricsV1() HTTPServerRequestMetricsV1 {
	return NewHttpServerRequestMetricsWithBucketsV1(DefaultDurationMsBuckets)
}

func NewHttpServerRequestMetricsWithBucketsV1(requestTimeMsBuckets []float64) HTTPServerRequestMetricsV1 {
	m := &httpServerRequestMetricsV1{
		nbRequests:    newCounterVec(metricsNamespace, metricsSubsystemHttpServer, "nb_req", nil, []string{metricsLabelMethod, metricsLabelStatusCode, metricsLabelSupplierOldId}),
		requestTimeMs: newHistogramVec(metricsNamespace, metricsSubsystemHttpServer, "req_duration_ms", nil, requestTimeMsBuckets, []string{metricsLabelMethod}),
	}
	return m
}
