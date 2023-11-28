package v1

import (
	"github.com/happywbfriends/metrics/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

func NewKafkaProducerMetrics() KafkaProducerMetrics {
	m := kafkaProducerMetrics{
		nbDone:  metrics.NewCounterVec(metrics.MetricsNamespace, MetricsSubsystemKafka, "nb_msg_produced", nil, []string{MetricsLabelTopic}),
		nbError: metrics.NewCounterVec(metrics.MetricsNamespace, MetricsSubsystemKafka, "nb_error_produced", nil, []string{MetricsLabelTopic}),
	}

	return &m
}

type KafkaProducerMetrics interface {
	IncNbDone(topic string)
	IncNbError(topic string)
}

type kafkaProducerMetrics struct {
	nbDone  *prometheus.CounterVec
	nbError *prometheus.CounterVec
}

func (m *kafkaProducerMetrics) IncNbDone(topic string) {
	m.nbDone.WithLabelValues(topic).Inc()
}

func (m *kafkaProducerMetrics) IncNbError(topic string) {
	m.nbError.WithLabelValues(topic).Inc()
}
