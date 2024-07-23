package main

import (
	"net/http"

	"github.com/happywbfriends/metrics/metered"
	v1 "github.com/happywbfriends/metrics/v1"
)

func main() {
	clientName := "example_api"

	httpextrametrics := v1.NewHTTPClientMetrics()
	tp := metered.NewTracerProvider(
		httpextrametrics,
	)
	_ /*httpTransport*/ = metered.NewTrasport(
		clientName,
		http.DefaultTransport,
		httpextrametrics,
		func(ht *metered.HTTPClientTracker) {
			ht.SetTracerProvider(tp)
		},
	)

	// client := http.Client{
	// 	Transport: httpTransport,
	// }

	// client.Do(req)
}
