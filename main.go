package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"

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

	today, err := fetchTodayTotalSum(r)
	if err != nil {
		return err
	}
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
