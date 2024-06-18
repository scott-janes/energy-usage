package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog"
)

const URL = "https://api.carbonintensity.org.uk/intensity/date/%s"

func GetC02Data(ctx context.Context, date string, logger *zerolog.Logger) (*C02IntensityResponse, error) {
	logger.Info().Msgf("Getting C02 data for date: %s", date)
	url := fmt.Sprintf(URL, date)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Unable to get C02 data: %w", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	var result C02IntensityResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("Can not unmarshal json: %w", err)
	}

	return &result, nil
}

type C02IntensityRecord struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Intensity struct {
		Forecast int    `json:"forecast"`
		Actual   int    `json:"actual"`
		Index    string `json:"index"`
	} `json:"intensity"`
}

type C02IntensityResponse struct {
	Data []C02IntensityRecord `json:"data"`
}
