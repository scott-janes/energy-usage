package dao

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"github.com/scott-janes/energy-usage/shared/storage"
)

func StoreConsumptionData(ctx context.Context, data *ConsumptionOutput, date string, store *storage.PostgresStore, logger *zerolog.Logger) error {
	logger.Info().Msgf("Storing consumption data for date %s", date)
	tx, err := store.BeginTransaction()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, record := range data.Results {
		if !strings.HasPrefix(record.IntervalStart, date) {
			continue
		}

		_, err = tx.Exec("INSERT INTO energy_usage.octopus_usage (consumption, from_timestamp, to_timestamp) VALUES ($1, $2, $3)", record.Consumption, record.IntervalStart, record.IntervalEnd)
		if err != nil {
			return fmt.Errorf("failed to insert consumption record: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
