package main

import (
	metricsv1 "github.com/happywbfriends/metrics/v1"
	"net/http"
	"time"
)

func HTTPServerExample() {
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

func HTTPServerWithHookExample() {
	httpServerMetrics := metricsv1.NewHTTPServerMetrics()
	httpServer := &http.Server{
		ConnState: httpServerMetrics.OnStateChange,
	}
	httpServer.Handler = MyHandler{
		metrics: httpServerMetrics,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			panic("listen probe finished")
		}
	}()
}

type MyHandler struct {
	metrics metricsv1.HTTPServerMetrics
}

func (h MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	supplierOldId := 999

	timeBegin := time.Now()
	method := "GET/bar"
	status := http.StatusOK
	defer func() {
		h.metrics.ObserveRequestDuration(method, status, supplierOldId, time.Since(timeBegin))
		h.metrics.IncNbRequest(method, status, supplierOldId)
	}()

	w.WriteHeader(status)
}
