package main

import (
	"context"
	"database/sql"
	"errors"
	metricsv1 "github.com/happywbfriends/metrics/v1"
	"net/http"
	"time"
)

func HTTPClientExample() {
	httpServerMetrics := metricsv1.NewHTTPClientMetrics()

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

func sqlDbExample() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, _ := sql.Open("postgres", "...connection string...")
	defer db.Close()

	// запускаем считывание основных метрик инстанса БД раз в 5 сек
	dbMetrics := metricsv1.NewDbMetrics("MyDatabase")
	go metricsv1.DbMetricsHelper(dbMetrics, db, 5*time.Second, ctx)
}

func sqlDbQueryExample() {
	db, _ := sql.Open("postgres", "...connection string...")
	defer db.Close()

	dbQueryMetrics := metricsv1.NewDbQueryMetrics()
	_ = executeSomeUsefulQuery(db, dbQueryMetrics)
}

func executeSomeUsefulQuery(_ *sql.DB, m metricsv1.DbQueryMetrics) (e error) {
	db_name := "MyDatabase"
	query_name := "UsefulQuery"
	defer func(timeBegin time.Time) {
		m.ObserveRequestDuration(db_name, query_name, time.Since(timeBegin))
		if e != nil {
			m.IncNbError(db_name, query_name)
		} else {
			m.IncNbDone(db_name, query_name)
		}
	}(time.Now())

	// do something in db

	return nil
}
