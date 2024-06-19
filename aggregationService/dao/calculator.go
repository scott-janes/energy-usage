package dao

import (
	"context"
	"math"

	"github.com/rs/zerolog"
)

type MixPercentage struct {
	Fuel       string  `json:"fuel"`
	Percentage float64 `json:"perc"`
}

type DailySummary struct {
	Date                   string
	CarbonIntensity        float64
	AverageCarbonIntensity float64
	TotalEnergyUsed        float64
	TotalEnergyCostIncVat  float64
	TotalEnergyCostExcVat  float64
	MixPercentage          []MixPercentage
}

func CalculateDailySummary(ctx context.Context, data *[]EnergyData, date string, logger *zerolog.Logger) (*DailySummary, error) {
	logger.Info().Msgf("Calculating daily summary for date %s", date)

	fuelTypeTotals := make(map[string]float64)
	var totalEnergyUsed, totalEnergyCostIncVat, totalEnergyCostExcVat, totalCarbonIntensity float64

	for _, item := range *data {
		totalEnergyUsed += item.Usage
		totalEnergyCostIncVat += item.Usage * item.CostIncVat
		totalEnergyCostExcVat += item.Usage * item.CostExcVat
		totalCarbonIntensity += item.Usage * item.Intensity

		for _, gen := range item.Generations {
			fuelTypeTotals[gen.Fuel] += item.Usage * gen.Percentage
		}

	}

	var mixPercentage []MixPercentage

	for fuel, usage := range fuelTypeTotals {
		mixPercentage = append(mixPercentage, MixPercentage{Fuel: fuel, Percentage: formatFloat(usage / totalEnergyUsed)})
	}
	averageCarbonIntensity := totalCarbonIntensity / totalEnergyUsed

	return &DailySummary{
		Date:                   date,
		CarbonIntensity:        formatFloat(totalCarbonIntensity),
		AverageCarbonIntensity: formatFloat(averageCarbonIntensity),
		TotalEnergyUsed:        formatFloat(totalEnergyUsed),
		TotalEnergyCostIncVat:  formatFloat(totalEnergyCostIncVat),
		TotalEnergyCostExcVat:  formatFloat(totalEnergyCostExcVat),
		MixPercentage:          mixPercentage,
	}, nil
}

func formatFloat(f float64) float64 {
	return math.Round(f*100) / 100
}
