package kafka

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer(brokers []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{producer: producer}, nil
}

func (p *Producer) PublishSignal(topic string, signal *TradingSignal) error {
	data, err := json.Marshal(signal)
	if err != nil {
		return err
	}

	message := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(signal.Symbol),
		Value: sarama.ByteEncoder(data),
	}

	partition, offset, err := p.producer.SendMessage(message)
	if err != nil {
		return err
	}

	log.Printf("Published signal to topic %s, partition %d, offset %d", topic, partition, offset)
	return nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}