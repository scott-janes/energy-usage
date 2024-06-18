package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rs/zerolog/log"
	"github.com/scott-janes/energy-usage/cli-tool/dao"
	"github.com/scott-janes/energy-usage/shared/config"
	"github.com/spf13/cobra"
)

// topicsCmd represents the topics command
var topicsCmd = &cobra.Command{
	Use:   "topics",
	Short: "Creates required topics for energy usage app",
	Long:  `Checks if topics exist in kafka cluster. If not, it creates them.`,
	Run:   topics,
}

func topics(cmd *cobra.Command, args []string) {
	appConfig := dao.Config{}
	if err := config.LoadConfig(".", "config", "yaml", &appConfig); err != nil {
		log.Error().Stack().Err(err).Msg("Error loading config")
		os.Exit(0)
	}

	broker := fmt.Sprintf("%s:%d", appConfig.Kafka.Host, appConfig.Kafka.Port)
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
    log.Fatal().Err(err).Msg("Failed to create admin client")
	}

	defer adminClient.Close()

	for _, topic := range appConfig.Kafka.Topics {
		exists, err := topicExists(adminClient, topic)
		if err != nil {
      log.Fatal().Err(err).Msg("Failed to check if topic exists")
		}

		if !exists {
      log.Info().Msgf("Creating topic %s", topic)
			err = createTopic(adminClient, topic)
			if err != nil {
        log.Fatal().Err(err).Msg("Failed to create topic")
			}
		} else {
      log.Info().Msgf("Topic %s already exists", topic)
    }
	}
}

func topicExists(adminClient *kafka.AdminClient, topic string) (bool, error) {
	metadata, err := adminClient.GetMetadata(&topic, false, 10000)

	if err != nil {
		return false, err
	}

	for _, t := range metadata.Topics {
		if t.Topic == topic {
			return true, nil
		}
	}
	return false, nil
}

func createTopic(adminClient *kafka.AdminClient, topic string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	retentionPeriod := "604800000" // 1 week in milliseconds
	config := kafka.TopicSpecification{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
		Config:            map[string]string{"retention.ms": retentionPeriod},
	}
	_, err := adminClient.CreateTopics(ctx, []kafka.TopicSpecification{config}, nil)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(topicsCmd)
}

