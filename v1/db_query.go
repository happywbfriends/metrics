package v1

import (
	"github.com/happywbfriends/metrics/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type DbQueryMetrics interface {
	IncNbDone()
	IncNbError(error)
	ObserveRequestDuration(time.Duration)
}

func NewDbQueryMetrics(dbName, queryName string) DbQueryMetrics {
	labels := map[string]string{
		MetricsLabelSubject:               dbName,
		metrics.MetricsLabelDatabaseQuery: queryName,
	}

	return &dbRequestMetrics{
		durationMs: metrics.NewHistogram(metrics.MetricsNamespace, metrics.MetricsSubsystemDbQuery, "duration_ms", labels, metrics.DefaultDurationMsBuckets),
		nbDone:     metrics.NewCounter(metrics.MetricsNamespace, metrics.MetricsSubsystemDbQuery, "nb_done", labels),
		nbError:    metrics.NewCounter(metrics.MetricsNamespace, metrics.MetricsSubsystemDbQuery, "nb_error", labels),
	}
}

type dbRequestMetrics struct {
	durationMs prometheus.Histogram
	nbDone     prometheus.Counter
	nbError    prometheus.Counter
}

func (m *dbRequestMetrics) ObserveRequestDuration(duration time.Duration) {
	m.durationMs.Observe(float64(duration.Milliseconds()))
}
func (m *dbRequestMetrics) IncNbDone() {
	m.nbDone.Inc()
}
func (m *dbRequestMetrics) IncNbError(error) {
	m.nbError.Inc()
}

/*
Helper для быстрого расчета метрик запроса. Плюсом идет то, что метод сам анализирует ошибку и может инкрементить
нужные вспомогательные метрики.

Пример:

```

	func SomeDatabaseMethod() (e error){
		defer func(from time.Time) {
			metrics.DbQueryMetricsHelper(metrics, from, e)
		}(time.Now())

		// Your code goes here
	}

```
*/
func DbQueryMetricsHelper(m DbQueryMetrics, startTm time.Time, err error) {
	m.ObserveRequestDuration(time.Since(startTm))
	if err != nil {
		m.IncNbError(err)
	} else {
		m.IncNbDone()
	}
}
