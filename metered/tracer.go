package metered

import (
	"net/http/httptrace"
	"time"

	v1 "github.com/happywbfriends/metrics/v1"
)

type TracerProvider struct {
	metrics v1.HTTPClientMetricsExtra
}

type TraceMetricsHolder struct {
	dnsStart     time.Time
	dnsHost      string
	getConnStart time.Time
	getConnHost  string
	metrics      v1.HTTPClientMetricsExtra
}

func NewTracerProvider(metrics v1.HTTPClientMetricsExtra) *TracerProvider {
	return &TracerProvider{
		metrics: metrics,
	}
}

func (t *TracerProvider) Trace() *httptrace.ClientTrace {
	mytrace := &TraceMetricsHolder{
		metrics: t.metrics,
	}

	return &httptrace.ClientTrace{
		DNSStart: mytrace.DNSStart,
		DNSDone:  mytrace.DNSDone,
		GetConn:  mytrace.GetConn,
		GotConn:  mytrace.GotConn,
	}
}

func (t *TraceMetricsHolder) DNSStart(dnsInfo httptrace.DNSStartInfo) {
	t.dnsStart = time.Now()
	t.dnsHost = dnsInfo.Host
}

func (t *TraceMetricsHolder) DNSDone(dnsInfo httptrace.DNSDoneInfo) {
	coalesced := "true"
	if !dnsInfo.Coalesced {
		coalesced = "false"
	}

	t.metrics.ObserveDnsDuration(t.dnsHost, coalesced, time.Since(t.dnsStart))
}

func (t *TraceMetricsHolder) GetConn(hostPort string) {
	t.getConnStart = time.Now()
	t.getConnHost = hostPort
}
func (t *TraceMetricsHolder) GotConn(connInfo httptrace.GotConnInfo) {
	reused := "true"
	if !connInfo.Reused {
		reused = "false"
	}
	t.metrics.ObserveConnectDuration(t.getConnHost, reused, time.Since(t.getConnStart))
}
