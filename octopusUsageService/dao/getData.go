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

const URL = "https://api.octopus.energy/v1/electricity-meter-points/%s/meters/%s/consumption?period_from=%s&period_to=%s&order_by=period"

type ConsumptionOutput struct {
	Count    int                       `json:"count"`
	Next     string                    `json:"next"`
	Previous string                    `json:"previous"`
	Results  []ConsumptionResultOutput `json:"results"`
}

type ConsumptionResultOutput struct {
	Consumption   float64 `json:"consumption"`
	IntervalStart string  `json:"interval_start"`
	IntervalEnd   string  `json:"interval_end"`
}

type ConsumptionDate struct {
	fromDate string
	toDate   string
}

func getDates(date string, logger *zerolog.Logger) (*ConsumptionDate, error) {
	logger.Debug().Msg("Getting dates for consumption")
	fDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		logger.Error().Stack().Err(err).Msg("Error parsing date")
		return nil, err
	}
	fDate = time.Date(fDate.Year(), fDate.Month(), fDate.Day(), 0, 0, 0, 0, time.UTC)

	fromDate := fDate.Format("2006-01-02T15:04:05Z")
	toDate := fDate.AddDate(0, 0, 1).Format("2006-01-02T15:04:05Z")

	return &ConsumptionDate{
		fromDate: fromDate,
		toDate:   toDate,
	}, nil
}

func GetConsumptionData(ctx context.Context, date string, config *Config, logger *zerolog.Logger) (*ConsumptionOutput, error) {

	consumptionDate, err := getDates(date, logger)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(URL, config.Octopus.MPAN, config.Octopus.SerialNumber, consumptionDate.fromDate, consumptionDate.toDate)
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %w", err)
	}

	req.SetBasicAuth(config.Octopus.APIKey, "")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Unable to get consumption data: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Unable to read consumption data: %w", err)
	}
	var result ConsumptionOutput
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("Can not unmarshal json: %w", err)
	}
	return &result, nil
}
