package dao

import (
	"context"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jellydator/ttlcache/v3"
	"github.com/rs/zerolog/log"
	"github.com/scott-janes/energy-usage/shared/storage"
	"github.com/scott-janes/energy-usage/shared/types"
)

type Processor struct {
	Store  *storage.PostgresStore
	Config *Config
	Cache  *ttlcache.Cache[string, CacheValue]
}

func NewProcessor(store *storage.PostgresStore, config *Config, cache *ttlcache.Cache[string, CacheValue]) *Processor {
	return &Processor{
		Store:  store,
		Config: config,
		Cache:  cache,
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

  process := ShouldProcessRequest(ctx, p.Cache, &energyServiceEvent)

  if process {
    data, err := GetData(ctx, energyServiceEvent.Context.Date, p.Store, &logger)
    if err != nil {
      logger.Error().Stack().Err(err).Msg("Error getting consumption data")
      return
    }

    dailySummary, err := CalculateDailySummary(ctx, data, energyServiceEvent.Context.Date, &logger)
    if err != nil {
      logger.Error().Stack().Err(err).Msg("Error calculating daily summary")
      return
    }

    err = StoreDailySummaryData(ctx, dailySummary, energyServiceEvent.Context.Date, p.Store, &logger)
    if err != nil {
      logger.Error().Stack().Err(err).Msg("Error storing daily summary")
      return
    }
  }
}
