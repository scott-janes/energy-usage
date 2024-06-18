package dao

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/scott-janes/energy-usage/shared/storage"
)

func StoreData(ctx context.Context, data *MixResponse, date string, store *storage.PostgresStore, logger *zerolog.Logger) error {
	logger.Info().Msgf("Storing mix data for date %s", date)
	tx, err := store.BeginTransaction()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	for _, record := range data.Data {
		if !strings.HasPrefix(record.From, date) {
			continue
		}

		var foo uuid.NullUUID

		err := tx.QueryRow("INSERT INTO energy_usage.mix (from_timestamp, to_timestamp) VALUES ($1, $2) RETURNING id", record.From, record.To).Scan(&foo)
		if err != nil {
			return fmt.Errorf("failed to insert mix record: %w", err)
		}

		for _, gen := range record.Generation {
			_, err = tx.Exec("INSERT INTO energy_usage.mix_generation (mix_id, fuel, perc) VALUES ($1, $2, $3)", foo, gen.Fuel, gen.Perc)
			if err != nil {
				return fmt.Errorf("failed to insert mix generation record: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
