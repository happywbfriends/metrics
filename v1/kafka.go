package v1

import (
	"github.com/happywbfriends/metrics/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

func NewKafkaProducerMetrics() KafkaProducerMetrics {
	m := kafkaMetrics{
		nbDone:  metrics.NewCounterVec(metrics.MetricsNamespace, MetricsSubsystemKafka, "nb_msg_produced", nil, []string{MetricsLabelTopic}),
		nbError: metrics.NewCounterVec(metrics.MetricsNamespace, MetricsSubsystemKafka, "nb_error_produced", nil, []string{MetricsLabelTopic}),
	}

	return &m
}

func NewKafkaConsumerMetrics() KafkaConsumerMetrics {
	m := kafkaMetrics{
		nbDone:  metrics.NewCounterVec(metrics.MetricsNamespace, MetricsSubsystemKafka, "nb_msg_consumed", nil, []string{MetricsLabelTopic}),
		nbError: metrics.NewCounterVec(metrics.MetricsNamespace, MetricsSubsystemKafka, "nb_error_consumed", nil, []string{MetricsLabelTopic}),
	}

	return &m
}

type KafkaProducerMetrics interface {
	IncNbDone(topic string)
	IncNbError(topic string)
}

type KafkaConsumerMetrics interface {
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
