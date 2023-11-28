package v1

import (
	"github.com/happywbfriends/metrics/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

func NewKafkaProducerMetrics() KafkaMetrics {
	m := kafkaMetrics{
		nbDone:  metrics.NewCounterVec(metrics.MetricsNamespace, MetricsSubsystemKafka, "nb_done", nil, []string{MetricsLabelTopic}),
		nbError: metrics.NewCounterVec(metrics.MetricsNamespace, MetricsSubsystemKafka, "nb_error", nil, []string{MetricsLabelTopic}),
	}

	return &m
}

type KafkaMetrics interface {
	IncNbDone(topic string)
	IncNbError(topic string)
}

type kafkaMetrics struct {
	nbDone  *prometheus.CounterVec
	nbError *prometheus.CounterVec
}

func (m *kafkaMetrics) IncNbDone(topic string) {
	m.nbDone.WithLabelValues(topic).Inc()
}

func (m *kafkaMetrics) IncNbError(topic string) {
	m.nbError.WithLabelValues(topic).Inc()
}
