package dao

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rs/zerolog"
)

func PublishToKafka(ctx context.Context, message *[]byte, producer *kafka.Producer, config *Config, logger *zerolog.Logger) error {
	logger.Info().Msg("Publishing to kafka")
	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &config.Kafka.Producer.Topic, Partition: kafka.PartitionAny},
		Value:          *message,
	}, nil)
	if err != nil {
		return err
	}
	return nil
}
