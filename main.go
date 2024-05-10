package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alphatroya/redmine-helper-bot/redmine"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	csvFlag      = "csv"
	vacationFlag = "vacation-file"

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

var (
	infoLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errLogger  = log.New(os.Stderr, "", 0)
)

func main() {
	main := cobra.Command{
		Short: "Command for processing csv file for time entries and pass it to project management tool",
		Long: `
Command for processing a CSV file containing time entries and transferring it to a project management tool.

The CSV file should adhere to the following structure:
- It must consist of 3 rows.
- The first row should denote the task number, which is mandatory.
- The second row should represent the time entry amount in hours as a floating-point number. This field is optional.
  - If left empty, it will be treated as a filler field. The program will then calculate the remaining available hours for the day and distribute them among the filler entries.
- The third row should contain the time entry description, which is also mandatory.
The total hours logged in a day should not exceed 8 hours minus the already logged amount.
For instance, if 3 hours are logged today, the maximum allowable hours for today would be 5..
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

			vacationFile, err := cmd.Flags().GetString(vacationFlag)
			if err == nil {
				today := time.Now()
				isVacation, err := checkVacation(vacationFile, today)
				if err != nil {
					return err
				}
				if isVacation {
					infoLogger.Println("today is vacation, skipping")
					return nil
				}
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

func checkVacation(vacation string, today time.Time) (bool, error) {
	file, err := os.Open(vacation)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dateStr := scanner.Text()
		date, err := time.Parse("02.01.2006", dateStr)
		if err != nil {
			return false, err
		}

		if date.Day() == today.Day() && date.Month() == today.Month() && date.Year() == today.Year() {
			return true, nil
		}
	}

	return false, nil
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
