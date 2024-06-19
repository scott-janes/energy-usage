package model

import "github.com/google/uuid"


type MixPercentage struct {
  Fuel string `json:"fuel"`
  Perc float64 `json:"perc"`
}

type DailySummary struct {
  Id uuid.UUID `json:"id"`
  Date string `json:"date"`
  Carbonintensity float64 `json:"carbon_intensity"`
  Averagecarbonintensity float64 `json:"average_carbon_intensity"`
  Totalenergyused float64 `json:"total_energy_used"`
  Totalenergycostincvat float64 `json:"total_energy_cost_inc_vat"`
  Totalenergycostexvat float64 `json:"total_energy_cost_ex_vat"`
  Mixpercentage []MixPercentage `json:"mix_percentage"`
}
