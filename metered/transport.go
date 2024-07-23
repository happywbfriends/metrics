package metered

import (
	"net/http"
	"net/http/httptrace"
	"time"

	metricshwb "github.com/happywbfriends/metrics/v1"
)

type tracerProvider interface {
	Trace() *httptrace.ClientTrace
}

type ClientMetricsOpt func(*HTTPClientTracker)

func NewTrasport(
	clientName string,
	transport http.RoundTripper,
	hwbmetrics metricshwb.HTTPClientMetricsExtra,
	opts ...ClientMetricsOpt,
) http.RoundTripper {
	metricsTransport := newHTTPClientTracker(transport, clientName, hwbmetrics)
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(metricsTransport)
		}
	}

	return metricsTransport
}

func newHTTPClientTracker(
	next http.RoundTripper,
	client string,
	metrics metricshwb.HTTPClientMetrics,
) *HTTPClientTracker {
	return &HTTPClientTracker{
		next:    next,
		client:  client,
		metrics: metrics,
	}
}

func (t *HTTPClientTracker) SetTracerProvider(tracerProvider tracerProvider) {
	t.tracerProvider = tracerProvider
}

type HTTPClientTracker struct {
	next           http.RoundTripper
	client         string
	metrics        metricshwb.HTTPClientMetrics
	tracerProvider tracerProvider
}

func (t HTTPClientTracker) RoundTrip(req *http.Request) (*http.Response, error) {
	method := req.Method + req.URL.Path
	start := time.Now()

	if t.tracerProvider != nil {
		req = req.WithContext(httptrace.WithClientTrace(req.Context(), t.tracerProvider.Trace()))
	}

	resp, err := t.next.RoundTrip(req)
	if err != nil {
		t.metrics.IncNbError(t.client, method)

		return resp, err
	}

	t.metrics.ObserveRequestDuration(t.client, method, time.Since(start))
	t.metrics.IncNbDone(t.client, method, resp.StatusCode)

	return resp, err
}
