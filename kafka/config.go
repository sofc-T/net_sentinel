package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

type KafkaConfig struct {
	Brokers       []string
	ProducerTopic string
	ConsumerGroup string
}

func LoadKafkaConfig() KafkaConfig {
	return KafkaConfig{
		Brokers:       []string{"localhost:9092"}, // Change to your Kafka broker
		ProducerTopic: "network-metrics",
		ConsumerGroup: "network-monitor-group",
	}
}

func NewSyncProducer(brokers []string) sarama.SyncProducer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}

	return producer
}

func NewConsumerGroup(brokers []string, groupID string) sarama.ConsumerGroup {
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0 // Ensure Kafka version compatibility

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer group: %v", err)
	}

	return consumerGroup
}
