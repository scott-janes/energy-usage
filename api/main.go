package main

import (
	"github.com/scott-janes/energy-usage/api/datastore"
	"github.com/scott-janes/energy-usage/api/handlers"
	"gofr.dev/pkg/gofr"
)

func main() {
	// initialise gofr object
	app := gofr.New()

	// register route greet
	app.GET("/greet", func(ctx *gofr.Context) (interface{}, error) {

		return "Hello World!", nil
	})

	registerDailySummaryHandlers(app)
	// Runs the server, it will listen on the default port 8000.
	// it can be over-ridden through configs
	app.Run()
}

func registerDailySummaryHandlers(app *gofr.App) {
  s := datastore.NewDailySummaryRepo()
  h := handlers.NewDailySummaryHandler(s)
  app.GET("/daily-summary/{date}", h.GetByDate)
}
