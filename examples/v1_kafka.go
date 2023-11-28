package main

import (
	"github.com/IBM/sarama"
	metricsv1 "github.com/happywbfriends/metrics/v1"
)

func KafkaProducerExample() {
	metrics := metricsv1.NewKafkaProducerMetrics()
	producer, _ := newKafkaProducer([]string{"kafka:9092"}, "clientID", "user", "password", metrics)
	logger := newKafkaLogger(producer, "topic")
	logger.SendMessage("test")
}

func newKafkaProducer(urls []string, clientID, user, password string, metrics metricsv1.KafkaProducerMetrics) (sarama.AsyncProducer, error) {
	sarConf := sarama.NewConfig()
	sarConf.ClientID = clientID
	sarConf.Producer.RequiredAcks = sarama.NoResponse
	sarConf.Producer.Retry.Max = 0
	sarConf.Producer.Return.Successes = true
	sarConf.Net.SASL.Enable = true
	sarConf.Net.SASL.User = user
	sarConf.Net.SASL.Password = password
	sarConf.Net.SASL.Mechanism = sarama.SASLTypePlaintext

	producer, err := sarama.NewAsyncProducer(urls, sarConf)
	if err != nil {
		return nil, err
	}

	go func() {
		for msg := range producer.Successes() {
			metrics.IncNbDone(msg.Topic)
		}
	}()

	go func() {
		for err := range producer.Errors() {
			metrics.IncNbError(err.Msg.Topic)
		}
	}()

	return producer, nil
}

func newKafkaLogger(
	kafkaProducer sarama.AsyncProducer,
	kafkaTopic string,
) KafkaLogger {
	l := kafkaLogger{
		Producer: kafkaProducer,
		Topic:    kafkaTopic,
	}

	return &l
}

type KafkaLogger interface {
	SendMessage(msg string)
}

type kafkaLogger struct {
	Producer sarama.AsyncProducer
	Topic    string
}

func (l *kafkaLogger) SendMessage(msg string) {
	l.Producer.Input() <- &sarama.ProducerMessage{
		Topic: l.Topic,
		Value: sarama.ByteEncoder(msg),
	}
}
