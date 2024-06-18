package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/scott-janes/energy-usage/shared/storage"
)

type Generation struct {
	Fuel       string  `json:"fuel"`
	Percentage float64 `json:"perc"`
}

type EnergyData struct {
	FromTimestamp string       `json:"from_timestamp"`
	ToTimestamp   string       `json:"to_timestamp"`
	Usage         float64      `json:"usage"`
	Intensity     float64      `json:"intensity"`
	CostExcVat    float64      `json:"cost_exc_vat"`
	CostIncVat    float64      `json:"cost_inc_vat"`
	Generations   []Generation `json:"generations"`
}

func GetData(ctx context.Context, date string, store *storage.PostgresStore, logger *zerolog.Logger) (*[]EnergyData, error) {
	var data []EnergyData
	query := `
SELECT ou.from_timestamp,
       ou.to_timestamp,
       ou.consumption,
       co2.actual                                                          as c02_consumption,
       pr.value_exc_vat                                                    as cost_exc_vat,
       pr.value_inc_vat                                                    as cost_inc_vat,
       json_agg(DISTINCT mg.*) as generations
FROM energy_usage.octopus_usage ou
         JOIN energy_usage.mix m ON ou.from_timestamp = m.from_timestamp AND ou.to_timestamp = m.to_timestamp
         JOIN energy_usage.mix_generation mg ON m.id = mg.mix_id
         JOIN energy_usage.c02 co2 ON ou.from_timestamp = co2.from_timestamp AND ou.to_timestamp = co2.to_timestamp
         JOIN energy_usage.octopus_pricing pr
              ON ou.from_timestamp = pr.from_timestamp AND ou.to_timestamp = pr.to_timestamp
WHERE DATE(ou.from_timestamp) = $1
GROUP BY ou.from_timestamp, ou.to_timestamp, ou.consumption, co2.actual, pr.value_exc_vat, pr.value_inc_vat
ORDER BY ou.from_timestamp, ou.to_timestamp;
  `
	rows, err := store.QueryData(ctx, query, date)

	if err != nil {
		return nil, fmt.Errorf("error getting data: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var singleData EnergyData
		var generations string

		if err := rows.Scan(&singleData.FromTimestamp, &singleData.ToTimestamp, &singleData.Usage, &singleData.Intensity, &singleData.CostExcVat, &singleData.CostIncVat, &generations); err != nil {
			return nil, fmt.Errorf("error scanning rows: %w", err)
		}

		if err := json.Unmarshal([]byte(generations), &singleData.Generations); err != nil {
			return nil, fmt.Errorf("error unmarshalling generations: %w", err)
		}

		data = append(data, singleData)
	}

	return &data, nil
}
