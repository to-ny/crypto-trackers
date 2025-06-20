package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type SignalProducer interface {
	PublishSignal(ctx context.Context, topic string, signal *TradingSignal) error
	Close() error
}

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer(brokers []string) (SignalProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	return &Producer{producer: producer}, nil
}

func (p *Producer) PublishSignal(ctx context.Context, topic string, signal *TradingSignal) error {
	data, err := json.Marshal(signal)
	if err != nil {
		return fmt.Errorf("failed to marshal signal: %w", err)
	}

	message := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(signal.Symbol),
		Value: sarama.ByteEncoder(data),
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	partition, offset, err := p.producer.SendMessage(message)
	if err != nil {
		return fmt.Errorf("failed to send message to kafka: %w", err)
	}

	log.Printf("Published signal to topic %s, partition %d, offset %d", topic, partition, offset)
	return nil
}

func (p *Producer) Close() error {
	if p.producer != nil {
		return p.producer.Close()
	}
	return nil
}
