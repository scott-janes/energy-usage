package handlers

import (
	"github.com/scott-janes/energy-usage/api/datastore"
	"gofr.dev/pkg/gofr"
)

type DailySummaryHandler struct {
	store datastore.DailySummaryRepository
}

func NewDailySummaryHandler(store datastore.DailySummaryRepository) *DailySummaryHandler {
	return &DailySummaryHandler{store: store}
}

func (h *DailySummaryHandler) GetByDate(ctx *gofr.Context) (interface{}, error) {
	date := ctx.PathParam("date")
	result, err := h.store.GetByDate(ctx, date)

	if err != nil {
		return nil, err
	}
	return result, nil
}
