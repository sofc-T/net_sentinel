package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

// SendMessage sends a message to Kafka
func SendMessage(producer sarama.SyncProducer, topic, message string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	_, _, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send Kafka message: %v", err)
		return err
	}

	log.Printf("Sent message to Kafka topic %s: %s", topic, message)
	return nil
}
