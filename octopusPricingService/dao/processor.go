package dao

import (
	"context"
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gookit/goutil/arrutil"
	"github.com/rs/zerolog/log"
	"github.com/scott-janes/energy-usage/shared/storage"
	"github.com/scott-janes/energy-usage/shared/types"
)

type Processor struct {
	store    *storage.PostgresStore
	producer *kafka.Producer
	config   *Config
}

func NewProcessor(store *storage.PostgresStore, producer *kafka.Producer, config *Config) *Processor {
	return &Processor{
		store:    store,
		producer: producer,
		config:   config,
	}
}

func (p *Processor) Process(ctx context.Context, message *kafka.Message) {
	energyServiceEvent := types.EnergyServiceEvent{}

	if err := json.Unmarshal(message.Value, &energyServiceEvent); err != nil {
		log.Error().Stack().Err(err).Msg("Error unmarshalling message")
		return
	}

	logger := log.With().Str("eventId", energyServiceEvent.ID).Logger()

	if err := types.Validate(energyServiceEvent); err != nil {
		logger.Error().Stack().Err(err).Msg("Error validating event")
		return
	}

	if !arrutil.StringsHas(energyServiceEvent.ChildProcesses, p.config.ServiceName) {
		logger.Info().Msg("Skipping event as it does not have the current service as a child process")
		return
	}

	result, err := GetPricingData(ctx, energyServiceEvent.Context.Date, p.config, &logger)
	if err != nil {
		logger.Error().Stack().Err(err).Msg("Error getting consumption data")
		return
	}

  formattedResult := FormatData(ctx, result, &logger)

	err = StorePricingData(ctx, formattedResult, energyServiceEvent.Context.Date, p.store, &logger)

	energyServiceEvent.Context.Status = "OK"
	energyServiceEvent.Context.Service = p.config.ServiceName
	energyServiceEvent.Completed_At = time.Now().Format(time.RFC3339)

	newMessage, err := json.Marshal(energyServiceEvent)
	if err != nil {
		logger.Error().Stack().Err(err).Msg("Error marshalling message")
		return
	}

	err = PublishToKafka(ctx, &newMessage, p.producer, p.config, &logger)

	if err != nil {
		logger.Error().Stack().Err(err).Msg("Error publishing message to kafka")
		return
	}
}
