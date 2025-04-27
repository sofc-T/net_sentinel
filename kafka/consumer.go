package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type ConsumerHandler struct{}

func (c *ConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (c *ConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		fmt.Printf("Received message: %s\n", string(message.Value))
		session.MarkMessage(message, "")
	}
	return nil
}

func StartConsumer(brokers []string, groupID, topic string) {
	consumerGroup := NewConsumerGroup(brokers, groupID)
	handler := &ConsumerHandler{}

	ctx := context.Background()
	for {
		err := consumerGroup.Consume(ctx, []string{topic}, handler)
		if err != nil {
			log.Fatalf("Error in Kafka consumer: %v", err)
		}
	}
}




