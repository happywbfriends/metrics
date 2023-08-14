package main

import (
	metricsv1 "github.com/happywbfriends/metrics/v1"
	"net/http"
	"time"
)

func v1HTTPServerExample() {
	supplierOldId := 999

	httpServerMetrics := metricsv1.NewHTTPServerMetrics()

	// Обработчик /bar
	http.HandleFunc("GET/bar", func(w http.ResponseWriter, r *http.Request) {
		method := "GET/bar"
		httpServerMetrics.IncNbConnections()
		defer httpServerMetrics.DecNbConnections()

		timeBegin := time.Now()
		status := http.StatusOK
		defer func() {
			httpServerMetrics.ObserveRequestDuration(method, status, supplierOldId, time.Since(timeBegin))
			httpServerMetrics.IncNbRequest(method, status, supplierOldId)
		}()

		w.WriteHeader(status)
	})
}
