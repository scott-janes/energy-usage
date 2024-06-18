package dao

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/scott-janes/energy-usage/shared/storage"
)

func StoreData(ctx context.Context, response *C02IntensityResponse, date string, store *storage.PostgresStore, logger *zerolog.Logger) error {
  logger.Info().Msgf("Storing c02 data for date %s", date)
	tx, err := store.BeginTransaction()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO energy_usage.c02 (from_timestamp, to_timestamp, forecast, actual) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}

	for _, record := range response.Data {
		if _, err := stmt.Exec(record.From, record.To, record.Intensity.Forecast, record.Intensity.Actual); err != nil {
			return fmt.Errorf("failed to insert c02 record: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
