package main

import (
	"errors"
	metricsv1 "github.com/happywbfriends/metrics/v1"
	"net/http"
	"time"
)

func v1HTTPServerExample() {
	supplierOldId := 999

	httpServerMetrics := metricsv1.NewHTTPServerMetrics()

	// Обработчик /bar
	http.HandleFunc("GET/bar", func(w http.ResponseWriter, r *http.Request) {
		httpServerMetrics.IncNbConnections()
		defer httpServerMetrics.DecNbConnections()

		timeBegin := time.Now()
		method := "GET/bar"
		status := http.StatusOK
		defer func() {
			httpServerMetrics.ObserveRequestDuration(method, status, supplierOldId, time.Since(timeBegin))
			httpServerMetrics.IncNbRequest(method, status, supplierOldId)
		}()

		w.WriteHeader(status)
	})
}

func v1HTTPClientExample() {
	httpServerMetrics := metricsv1.NewHttpClientMetrics()

	timeBegin := time.Now()
	client := "foo_api"
	method := "POST/auth"
	status := http.StatusForbidden

	httpServerMetrics.IncNbDone(client, method, status)
	httpServerMetrics.ObserveRequestDuration(client, method, time.Since(timeBegin))
	err := errors.New("test error")
	if err != nil {
		httpServerMetrics.IncNbError(client, method)
	}
}
