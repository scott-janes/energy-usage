package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/scott-janes/energy-usage/cli-tool/dao"
	"github.com/scott-janes/energy-usage/shared/config"
	"github.com/scott-janes/energy-usage/shared/storage"
	"github.com/scott-janes/energy-usage/shared/types"
	"github.com/spf13/cobra"
)

var date string
var backoffDays int

// processCmd represents the process command
var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Process data for a given date, or the latest date",
	Run:   process,
}

func process(cmd *cobra.Command, args []string) {
	appConfig := dao.Config{}
	if err := config.LoadConfig(".", "config", "yaml", &appConfig); err != nil {
		log.Error().Stack().Err(err).Msg("Error loading config")
		os.Exit(0)
	}

	if date == "" {
    dateFromDb := getNewestDateFromDb(&appConfig)
    if dateFromDb != nil {
      date = *dateFromDb
    }
	}

	if !isValidDateFormat(date) {
		log.Fatal().Msgf("Invalid date format: %s", date)
	}

	if backoffDays <= 0 {
		log.Fatal().Msg("Backoff days must be greater than 0")
	}

	processDates(date, &appConfig)
	fmt.Println("process called")
}

func isValidDateFormat(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

func processDates(date string, appConfig *dao.Config) {
	broker := fmt.Sprintf("%s:%d", appConfig.Kafka.Host, appConfig.Kafka.Port)
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create producer")
	}

	defer producer.Close()

	startTime, err := time.Parse("2006-01-02", date)

	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to parse date: %s", date)
	}

	currentTime := time.Now().AddDate(0, 0, -backoffDays)

	daysProcessed := 0
	for date := startTime; !date.After(currentTime); date = date.AddDate(0, 0, 1) {
		dateToRun := date.Format("2006-01-02")
		log.Printf("Starting on %s", dateToRun)
		err = func() error {
			return sendEvent(producer, dateToRun, appConfig.Kafka.Topics[0])
		}()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to send event")
		}
		daysProcessed++

		if daysProcessed%7 == 0 {
			log.Printf("Processed 7 days, waiting for 5 seconds....")
			time.Sleep(5 * time.Second)
		}
	}

}

func sendEvent(producer *kafka.Producer, dateToRun string, s string) error {
	deliveryChan := make(chan kafka.Event)

	defer close(deliveryChan)

	event := types.EnergyServiceEvent{
		ID:             uuid.New().String(),
		ChildProcesses: []string{"c02Service", "mixService", "octopusUsageService", "octopusPricingService"},
		Context: types.EnergyServiceContext{
			RequestID: uuid.New().String(),
			Date:      dateToRun,
		},
	}

	message, err := json.Marshal(event)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to marshal event")
	}

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &s, Partition: kafka.PartitionAny},
		Value:          message,
	}, deliveryChan)

	if err != nil {
		log.Error().Err(err).Msg("Failed to produce message")
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}
	return nil
}

func getNewestDateFromDb(appConfig *dao.Config) *string {
	store, err := storage.NewPostgresStore(appConfig.Database.Host, appConfig.Database.Port, appConfig.Database.User, appConfig.Database.Password, appConfig.Database.Name)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating PostgresStore")
		return nil
	}
	log.Info().Msg("PostgresStore created")

	var latestDate time.Time

	err = store.QueryRowData(context.Background(), "select MAX(date) from energy_usage.daily_summary").Scan(&latestDate)

	if err != nil {
		log.Error().Err(err).Msg("Error getting latest date from DB")
		return nil
	}
  var formattedDate string
  formattedDate = latestDate.Format("2006-01-02")
	return &formattedDate
}

func init() {
	processCmd.Flags().StringVarP(&date, "date", "d", "", "Date to process in YYYY-MM-DD format")
	processCmd.Flags().IntVarP(&backoffDays, "backoff-days", "b", 2, "Number of days to subtract from today's date to set the end date for processing")
	rootCmd.AddCommand(processCmd)
}
