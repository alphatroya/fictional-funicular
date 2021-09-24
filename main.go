package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/alphatroya/redmine-helper-bot/redmine"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	csvFlag = "csv"

	credEnv  = "REDMINE_API_KEY"
	hostENV  = "REDMINE_HOST"
	debugENV = "DEBUG"
)

func init() {
	for _, key := range []string{
		credEnv,
		hostENV,
		debugENV,
	} {
		if err := viper.BindEnv(key); err != nil {
			errLogger.Fatalf("Binding env error: %v", err)
		}
	}
}

var infoLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
var errLogger = log.New(os.Stderr, "", 0)

func main() {
	main := cobra.Command{
		Short:   "Command for processing csv file for time entries and pass it to project management tool",
		Version: "0.1",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !viper.IsSet(debugENV) {
				infoLogger.SetOutput(io.Discard)
			}
			host := viper.GetString(hostENV)
			token := viper.GetString(credEnv)

			if len(host) == 0 || len(token) == 0 {
				return errors.New(credEnv + " or " + hostENV + " are not set")
			}

			path, err := cmd.Flags().GetString(csvFlag)
			if err != nil {
				return err
			}

			return run(host, token, path)
		},
	}
	main.Flags().String(csvFlag, "", "path to csv file")
	if err := main.MarkFlagRequired(csvFlag); err != nil {
		errLogger.Fatalf("Failed to make csv flag required: %v", err)
	}

	if err := main.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(host string, token string, csv string) error {
	file, err := os.Open(csv)
	if err != nil {
		return err
	}
	defer file.Close()

	entries, err := parseCSV(file)
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		infoLogger.Println("nothing to fill")
		return nil
	}

	infoLogger.Printf("before processing: %+v\n", entries)
	client := http.Client{}
	storage := &storage{
		host:  host,
		token: token,
	}
	r := redmine.NewClientManager(&client, storage, 0)

	resp, err := r.TodayTimeEntries()
	if err != nil {
		return err
	}

	today := 0.0
	for _, item := range resp {
		today += float64(item.Hours)
	}
	infoLogger.Printf("counted hours today: %f\n", today)

	entries, err = refillMissedHours(entries, today)
	if err != nil {
		return err
	}

	infoLogger.Printf("after processing: %+v\n", entries)

	for i, item := range entries {
		res, err := r.FillHoursRequest(item.task, item.hours, item.comment, "")
		if err != nil {
			return err
		}
		infoLogger.Printf("item #%d filled by result: %+v", i, res)
	}

	return nil
}

func refillMissedHours(entries []fillItem, today float64) ([]fillItem, error) {
	const todayGoal float64 = 8

	notFilledCount := 0
	var alreadyFilled float64
	for i, item := range entries {
		hours := item.hours
		if len(hours) == 0 {
			notFilledCount++
			continue
		}
		//nolint 64 is obvious magic number here
		f, err := strconv.ParseFloat(hours, 64)
		if err != nil {
			return entries, fmt.Errorf("Parsing float error: %w, item at line %d (%s) is not a float", err, i, hours)
		}
		alreadyFilled += f
	}

	if alreadyFilled >= todayGoal {
		return entries, errors.New("Cancel task, you already reach the goal")
	}

	remain := fmt.Sprintf("%f", (todayGoal-alreadyFilled)/float64(notFilledCount))
	result := make([]fillItem, 0, len(entries))
	for _, item := range entries {
		if len(item.hours) == 0 {
			item.hours = remain
		}
		result = append(result, item)
	}
	return result, nil
}
