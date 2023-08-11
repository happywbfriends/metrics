package main

import (
	"context"
	"database/sql"
	"github.com/happywbfriends/metrics/metrics"
	"log"
	"net/http"
	"time"
)

func httpServerExample() {
	// Метрики инстанса сервера
	httpServerMetrics := metrics.NewHttpServerMetrics()

	// Метрики конкретного запроса
	httpServerRequestMetricsBar := metrics.NewHttpServerRequestMetrics("/bar")

	supplierOldId := 999

	// Обработчик для 404
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		httpServerMetrics.IncNbConnections()
		defer httpServerMetrics.DecNbConnections()

		httpServerMetrics.IncNotFound(r.URL.Path, supplierOldId)
	})

	// Обработчик /bar
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		httpServerMetrics.IncNbConnections()
		defer httpServerMetrics.DecNbConnections()

		since := time.Now()
		status := http.StatusOK
		defer func() {
			httpServerRequestMetricsBar.RequestDuration(time.Since(since), supplierOldId)
			httpServerRequestMetricsBar.IncNbRequest(status, supplierOldId)
		}()

		w.WriteHeader(status)
	})

	log.Fatal(http.ListenAndServe(":80", nil))
}

func sqlDbExample() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, _ := sql.Open("postgres", "...connection string...")
	defer db.Close()

	// запускаем считываение основных метрик инстанса БД раз в 5 сек
	dbMetrics := metrics.NewDbMetrics("MyDatabase")
	go metrics.DbMetricsHelper(dbMetrics, db, 5*time.Second, ctx)

	dbQueryMetrics := metrics.NewDbQueryMetrics("MyDatabase", "UsefulQuery")
	_ = executeSomeUsefulQuery(db, dbQueryMetrics)
}

func executeSomeUsefulQuery(_ *sql.DB, queryMetrics metrics.IDbQueryMetrics) (e error) {
	defer func(since time.Time) {
		metrics.DbQueryMetricsHelper(queryMetrics, since, e)
	}(time.Now())

	// do something in db
	return nil
}
