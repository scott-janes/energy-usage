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

const URL = "https://api.carbonintensity.org.uk/generation/%s/%s"

type MixDate struct {
	fromDate string
	toDate   string
}

func formatDate(date string) (*MixDate, error) {
	fDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return nil, err
	}
	fDate = time.Date(fDate.Year(), fDate.Month(), fDate.Day(), 0, 0, 0, 0, time.UTC)

	fromDate := fDate.Format("2006-01-02T15:04:05Z")
	toDate := fDate.AddDate(0, 0, 1).Format("2006-01-02T15:04:05Z")

	return &MixDate{fromDate, toDate}, nil

}

func GetMixData(ctx context.Context, date string, logger *zerolog.Logger) (*MixResponse, error) {
  logger.Info().Msgf("Getting mix data for date: %s", date)
	mixDate, err := formatDate(date)
	if err != nil {
		return nil, fmt.Errorf("Unable to format date: %w", err)
	}

	url := fmt.Sprintf(URL, mixDate.fromDate, mixDate.toDate)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Unable to get mix data: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	var result MixResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("Can not unmarshal json: %w", err)
	}

	return &result, nil
}

type GenerationMix struct {
	Fuel string  `json:"fuel"`
	Perc float64 `json:"perc"`
}

type MixResponseRecord struct {
	From       string          `json:"from"`
	To         string          `json:"to"`
	Generation []GenerationMix `json:"generationmix"`
}

type MixResponse struct {
	Data []MixResponseRecord `json:"data"`
}

type MixRecord struct {
	From    string
	To      string
	Max     int
	Average int
	Min     int
	Index   string
}
