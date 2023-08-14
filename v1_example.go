package main

import (
	v1 "github.com/happywbfriends/metrics/v1"
	"net/http"
	"time"
)

func v1HTTPServerExample() {
	supplierOldId := 999

	timeBegin := time.Now()
	status := http.StatusNotFound
	httpServerMetrics := v1.NewHTTPServerMetrics()
	httpServerMetrics.IncNbRequest("POST/foo", status, supplierOldId)
	if status == http.StatusOK {
		httpServerMetrics.ObserveOkRequestDuration("POST/foo", time.Since(timeBegin))
	}
}
