package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/alphatroya/redmine-helper-bot/redmine"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	csvFlag = "csv"

	tokenENV = "REDMINE_API_KEY"
	hostENV  = "REDMINE_HOST"
)

func init() {
	err := viper.BindEnv(tokenENV)
	if err != nil {
		panic(err)
	}
	err = viper.BindEnv(hostENV)
	if err != nil {
		panic(err)
	}
}

func main() {
	main := cobra.Command{
		Short:   "Command for processing csv file for time entries and pass it to project management tool",
		Version: "0.1",
		RunE: func(cmd *cobra.Command, args []string) error {
			host := viper.GetString(hostENV)
			token := viper.GetString(tokenENV)

			if len(host) == 0 || len(token) == 0 {
				return errors.New(tokenENV + " or " + hostENV + " are not set")
			}

			path, err := cmd.Flags().GetString(csvFlag)
			if err != nil {
				return err
			}

			return runExecution(host, token, path)
		},
	}
	main.Flags().String(csvFlag, "", "Path to csv file")
	main.MarkFlagRequired(csvFlag)

	err := main.Execute()
	if err != nil {
		panic(err)
	}
}

func runExecution(host string, token string, csv string) error {
	file, err := os.Open(csv)
	if err != nil {
		return err
	}
	defer file.Close()

	entries, err := parseCSV(file)
	if err != nil {
		return err
	}

	fmt.Printf("Before: %+v\n", entries)
	client := http.Client{}
	storage := &Storage{
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
	fmt.Printf("%f\n", today)

	entries, err = refillMissedHours(entries, today)
	if err != nil {
		return err
	}

	fmt.Printf("After: %+v\n", entries)

	for i, item := range entries {
		res, err := r.FillHoursRequest(item.task, item.hours, item.comment, "")
		if err != nil {
			return err
		}
		fmt.Printf("Item %d filled by result: %+v", i, res)
	}

	return nil
}

func refillMissedHours(entries []entry, today float64) ([]entry, error) {
	const todayGoal float64 = 8

	notFilledCount := 0
	var alreadyFilled float64
	for i, item := range entries {
		hours := item.hours
		if len(hours) == 0 {
			notFilledCount++
			continue
		}
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
	result := make([]entry, 0, len(entries))
	for _, item := range entries {
		if len(item.hours) == 0 {
			item.hours = remain
		}
		result = append(result, item)
	}
	return result, nil
}

type entry struct {
	task    string
	hours   string
	comment string
}

func parseCSV(file *os.File) ([]entry, error) {
	result := make([]entry, 0)
	reader := csv.NewReader(file)

	for line := 0; ; line++ {
		record, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return result, nil
			}
			return nil, err
		}
		if len(record) < 3 {
			return result, fmt.Errorf("Parsing line %d failed. Line should have more 3 items", line)
		}
		result = append(result, entry{
			task:    record[0],
			hours:   record[1],
			comment: record[2],
		})
	}
}
