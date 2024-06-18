package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

const URL = "https://api.octopus.energy/v1/products/%s/electricity-tariffs/%s/standard-unit-rates/?period_from=%s&period_to=%s"

type PricingOutput struct {
	Count    int                `json:"count"`
	Next     string             `json:"next"`
	Previous string             `json:"previous"`
	Results  []PricingResultOut `json:"results"`
}

type PricingResultOut struct {
	ValueExcVat float64   `json:"value_exc_vat"`
	ValueIncVat float64   `json:"value_inc_vat"`
	ValidFrom   time.Time `json:"valid_from"`
	ValidTo     time.Time `json:"valid_to"`
}

type PricingDate struct {
	fromDate string
	toDate   string
}

func getDates(date string, logger *zerolog.Logger) (*PricingDate, error) {
	logger.Debug().Msg("Getting dates for pricing")
	fDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		logger.Error().Stack().Err(err).Msg("Error parsing date")
		return nil, err
	}
	fDate = time.Date(fDate.Year(), fDate.Month(), fDate.Day(), 0, 0, 0, 0, time.UTC)

	fromDate := fDate.Format("2006-01-02T15:04:05Z")
	toDate := fDate.AddDate(0, 0, 1).Format("2006-01-02T15:04:05Z")

	return &PricingDate{
		fromDate: fromDate,
		toDate:   toDate,
	}, nil
}

func GetPricingData(ctx context.Context, date string, config *Config, logger *zerolog.Logger) (*PricingOutput, error) {
  pricingDate, err := getDates(date, logger)
  if err != nil {
    return nil, err
  }

  url := fmt.Sprintf(URL, config.Octopus.ProductCode, config.Octopus.TarrifCode, pricingDate.fromDate, pricingDate.toDate)
  client := &http.Client{}

  req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
  if err != nil {
    logger.Error().Stack().Err(err).Msg("Error creating request")
    return nil, err
  }

  resp, err := client.Do(req)
  if err != nil {
    return nil, fmt.Errorf("Unable to get pricing data: %w", err)
  }

  defer resp.Body.Close()

  body, err := io.ReadAll(resp.Body)
  if err != nil {
    return nil, fmt.Errorf("Unable to read pricing data: %w", err)
  }

  var result PricingOutput
  if err := json.Unmarshal(body, &result); err != nil {
    return nil, fmt.Errorf("Can not unmarshal json: %w", err)
  }
  return &result, nil
}
