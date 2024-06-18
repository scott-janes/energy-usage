package dao

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"github.com/scott-janes/energy-usage/shared/storage"
)

func StorePricingData(ctx context.Context, data *[]PricingResultOut, date string, store *storage.PostgresStore, logger *zerolog.Logger) error {
	logger.Info().Msgf("Storing pricing data for date %s", date)
	tx, err := store.BeginTransaction()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()


	for _, record := range *data {
		if !strings.HasPrefix(record.ValidFrom.Format("2006-01-02"), date) {
			continue
		}

		_, err = tx.Exec("INSERT INTO energy_usage.octopus_pricing (from_timestamp, to_timestamp, value_inc_vat, value_exc_vat) VALUES ($1, $2, $3, $4)", record.ValidFrom, record.ValidTo, record.ValueIncVat, record.ValueExcVat)
		if err != nil {
			return fmt.Errorf("failed to insert pricing record: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
