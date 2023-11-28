package v1

import (
	"github.com/happywbfriends/metrics/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

func NewKafkaMetrics() KafkaMetrics {
	m := kafkaMetrics{
		nbMsgProduced:   metrics.NewCounterVec(metrics.MetricsNamespace, MetricsSubsystemKafka, "nb_msg_produced", nil, []string{MetricsLabelTopic}),
		nbErrorProduced: metrics.NewCounterVec(metrics.MetricsNamespace, MetricsSubsystemKafka, "nb_error_produced", nil, []string{MetricsLabelTopic}),
		nbMsgConsumed:   metrics.NewCounterVec(metrics.MetricsNamespace, MetricsSubsystemKafka, "nb_msg_consumed", nil, []string{MetricsLabelTopic}),
		nbErrorConsumed: metrics.NewCounterVec(metrics.MetricsNamespace, MetricsSubsystemKafka, "nb_error_consumed", nil, []string{MetricsLabelTopic}),
	}

	return &m
}

type KafkaMetrics interface {
	IncNbMsgProduced(topic string)
	IncNbErrorProduced(topic string)
	IncNbMsgConsumed(topic string)
	IncNbErrorConsumed(topic string)
}

type kafkaMetrics struct {
	nbMsgProduced   *prometheus.CounterVec
	nbErrorProduced *prometheus.CounterVec
	nbMsgConsumed   *prometheus.CounterVec
	nbErrorConsumed *prometheus.CounterVec
}

func (m *kafkaMetrics) IncNbMsgProduced(topic string) {
	m.nbMsgProduced.WithLabelValues(topic).Inc()
}

func (m *kafkaMetrics) IncNbErrorProduced(topic string) {
	m.nbErrorProduced.WithLabelValues(topic).Inc()
}

func (m *kafkaMetrics) IncNbMsgConsumed(topic string) {
	m.nbMsgConsumed.WithLabelValues(topic).Inc()
}

func (m *kafkaMetrics) IncNbErrorConsumed(topic string) {
	m.nbErrorConsumed.WithLabelValues(topic).Inc()
}
