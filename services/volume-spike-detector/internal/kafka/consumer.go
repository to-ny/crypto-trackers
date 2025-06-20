package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type Consumer struct {
	client       sarama.ConsumerGroup
	topics       []string
	eventHandler func(*PriceEvent) error
}

type ConsumerGroupHandler struct {
	eventHandler func(*PriceEvent) error
}

func NewConsumer(brokers []string, groupID string, topics []string, eventHandler func(*PriceEvent) error) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true

	client, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer group: %w", err)
	}

	return &Consumer{
		client:       client,
		topics:       topics,
		eventHandler: eventHandler,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	handler := &ConsumerGroupHandler{
		eventHandler: c.eventHandler,
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := c.client.Consume(ctx, c.topics, handler); err != nil {
				log.Printf("Error consuming messages: %v", err)
				return err
			}
		}
	}
}

func (c *Consumer) Close() error {
	return c.client.Close()
}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var priceEvent PriceEvent
		if err := json.Unmarshal(message.Value, &priceEvent); err != nil {
			log.Printf("Error deserializing price event: %v", err)
			session.MarkMessage(message, "")
			continue
		}

		if err := h.eventHandler(&priceEvent); err != nil {
			log.Printf("Error handling price event: %v", err)
		}

		session.MarkMessage(message, "")
	}
	return nil
}
