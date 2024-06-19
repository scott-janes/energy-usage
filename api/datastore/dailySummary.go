package datastore

import (
	"encoding/json"

	"github.com/scott-janes/energy-usage/api/model"
	"gofr.dev/pkg/gofr"
)

type DailySummaryRepository interface {
	GetByDate(ctx *gofr.Context, date string) (*model.DailySummary, error)
}

type dailySummaryRepo struct {
}

func NewDailySummaryRepo() *dailySummaryRepo {
	return &dailySummaryRepo{}
}

func (d *dailySummaryRepo) GetByDate(ctx *gofr.Context, date string) (*model.DailySummary, error) {
	var result model.DailySummary
	var mixPercentageJson []byte
	err := ctx.SQL.QueryRowContext(ctx, "select * from energy_usage.daily_summary where date=$1", date).Scan(&result.Id, &result.Date, &result.Carbonintensity, &result.Averagecarbonintensity, &result.Totalenergyused, &result.Totalenergycostincvat, &result.Totalenergycostexvat, &mixPercentageJson)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(mixPercentageJson, &result.Mixpercentage)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
