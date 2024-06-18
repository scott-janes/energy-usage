package dao

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

func FormatData(ctx context.Context, result *PricingOutput, logger *zerolog.Logger) *[]PricingResultOut {
	logger.Info().Msg("Formatting data")
	var formattedData []PricingResultOut
	interval := 30 * time.Minute

	for _, result := range result.Results {
		current := result.ValidFrom

		for current.Before(result.ValidTo) {
			end := current.Add(interval)
			if end.After(result.ValidTo) {
				end = result.ValidTo
			}

			formattedData = append(formattedData, PricingResultOut{
				ValidFrom:          current,
				ValidTo:            end,
				ValueIncVat:        result.ValueIncVat,
				ValueExcVat: result.ValueExcVat,
			})

			current = end
		}
	}
	return &formattedData
}
