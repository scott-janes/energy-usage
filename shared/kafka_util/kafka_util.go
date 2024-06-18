package kafka_util

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Config struct {
	BootstrapServers string
	GroupID          string
	Topics           []string
}

func SetupKafkaConsumer(config *Config) (*kafka.Consumer, error) {
	consumerConfig := kafka.ConfigMap{
		"bootstrap.servers": config.BootstrapServers,
		"group.id":          config.GroupID,
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(&consumerConfig)

	if err != nil {
		return nil, err
	}

	err = consumer.SubscribeTopics(config.Topics, nil)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

func SetupKafkaProducer(bootstrapServers string) (*kafka.Producer, error) {
  producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapServers})
  if err != nil {
    return nil, err
  }
  return producer, nil
}
