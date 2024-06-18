package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/scott-janes/energy-usage/shared/storage"
)

func StoreDailySummaryData(ctx context.Context, data *DailySummary, date string, store *storage.PostgresStore, logger *zerolog.Logger) error {
	logger.Info().Msgf("Storing daily summary data for date %s", date)
	const query = "INSERT INTO energy_usage.daily_summary (date, carbon_intensity, average_carbon_intensity, total_energy_used, total_energy_cost_inc_vat, total_energy_cost_exc_vat, mix_percentage) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	jsonMixPercentage, err := json.Marshal(data.MixPercentage)
	if err != nil {
		return fmt.Errorf("failed to marshal mix percentage: %w", err)
	}
	err = store.ExecData(ctx, query, date, data.CarbonIntensity, data.AverageCarbonIntensity, data.TotalEnergyUsed, data.TotalEnergyCostIncVat, data.TotalEnergyCostExcVat, jsonMixPercentage)
	if err != nil {
		return fmt.Errorf("failed to insert data: %w", err)
	}
	return nil
}
