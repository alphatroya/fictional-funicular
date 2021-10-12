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
		Short: "Command for processing csv file for time entries and pass it to project management tool",
		Long: `
Command for processing csv file for time entries and pass it to project management tool

Passed .csv file should have the following structure:
- It should contain 3 rows.
- First row is a task number, it is required.
- Second row is a float number with the time entry amount (in hours). It's optional.
	- If the field is empty, it will be considered as a filler field. The program will calculate the spare
	amount of hours for today and divide it by fillers.
- Third row is a time entry description, the string, also required
The sum of the hours in the day should not exceed the 8 hours minus the already logged amount.
For example, if you log 3 hours today, your max hours limit is 5 for today.
		`,
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

	return makeFillRequests(r, entries)
}
