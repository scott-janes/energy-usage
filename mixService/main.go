package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/scott-janes/energy-usage/mixService/dao"
	"github.com/scott-janes/energy-usage/shared/config"
	"github.com/scott-janes/energy-usage/shared/kafka_util"
	"github.com/scott-janes/energy-usage/shared/storage"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Info().Msg("Fetching config")

	appConfig := dao.Config{}
	if err := config.LoadConfig(".", "config", "yaml", &appConfig); err != nil {
		log.Error().Stack().Err(err).Msg("Error loading config")
		os.Exit(0)
	}

	log.Info().Msgf("Starting %s", appConfig.ServiceName)

	store, err := storage.NewPostgresStore(appConfig.Database.Host, appConfig.Database.Port, appConfig.Database.User, appConfig.Database.Password, appConfig.Database.Name)

	if err != nil {
		log.Error().Stack().Err(err).Msg("Error creating PostgresStore")
		return
	}

	log.Info().Msg("PostgresStore created")

	consumerConfig := kafka_util.Config{
		BootstrapServers: fmt.Sprintf("%s:%d", appConfig.Kafka.Consumer.Host, appConfig.Kafka.Consumer.Port),
		GroupID:          appConfig.Kafka.Consumer.GroupID,
		Topics:           appConfig.Kafka.Consumer.Topics,
	}

	consumer, err := kafka_util.SetupKafkaConsumer(&consumerConfig)

	if err != nil {
		log.Error().Stack().Err(err).Msg("Error creating Kafka consumer")
		return
	}

	log.Info().Msg("Kafka consumer connection created")

	defer consumer.Close()

	producer, err := kafka_util.SetupKafkaProducer(fmt.Sprintf("%s:%d", appConfig.Kafka.Producer.Host, appConfig.Kafka.Producer.Port))
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error creating Kafka producer")
		return
	}

	log.Info().Msg("Kafka producer connection created")

	defer producer.Close()

	processor := dao.NewProcessor(store, producer, &appConfig)

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(2)

	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()
			consumeMessages(ctx, consumer, processor)
		}()
	}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	<-sigterm
	log.Info().Msg("Termination signal received, shutting down...")
	cancel()
	consumer.Close() // Ensure consumer is closed properly
	wg.Wait()

	log.Info().Msg("Application stopped")
}

func consumeMessages(ctx context.Context, consumer *kafka.Consumer, processor *dao.Processor) {
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Context cancelled, exiting")
			return
		default:
			message, err := consumer.ReadMessage(500)
			if err != nil {
				if kafkaError, ok := err.(kafka.Error); ok {
					if kafkaError.IsTimeout() {
						continue
					}
					if kafkaError.IsFatal() {
						log.Error().Err(err).Msg("Fatal error, stopping consumer")
						return
					}
				}
				log.Error().Err(err).Msg("Error consuming message")
				continue
			}

      go processor.Process(ctx, message)
		}
	}
}

