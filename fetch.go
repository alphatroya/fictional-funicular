package main

import "github.com/alphatroya/redmine-helper-bot/redmine"

func fetchTodayTotalSum(r redmine.Client) (float64, error) {
	resp, err := r.TodayTimeEntries()
	if err != nil {
		return 0, err
	}

	today := 0.0
	for _, item := range resp {
		today += float64(item.Hours)
	}
	infoLogger.Printf("total today hours sum: %f\n", today)
	return today, nil
}
