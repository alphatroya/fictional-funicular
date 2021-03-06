package main

import (
	"errors"
	"fmt"
	"strconv"
)

func refillMissedHours(entries []fillItem, today float64) ([]fillItem, error) {
	const todayGoal float64 = 8
	remainToday := todayGoal - today
	if remainToday <= 0 {
		return entries, errors.New("already logged 8 hours")
	}

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
			return entries, fmt.Errorf("refillMissedHours: parsing float item at line %d (%s) is not a float: %w", i, hours, err)
		}
		alreadyFilled += f
	}

	if almostEqualOfGreat(alreadyFilled, remainToday) {
		return entries, errors.New("refillMissedHours: you already reach the goal today")
	}

	remain := fmt.Sprintf("%.2f", (remainToday-alreadyFilled)/float64(notFilledCount))
	result := make([]fillItem, 0, len(entries))
	for _, item := range entries {
		if len(item.hours) == 0 {
			item.hours = remain
		}
		result = append(result, item)
	}
	return result, nil
}

const float64EqualityThreshold = 1e-9

func almostEqualOfGreat(a, b float64) bool {
	return a-b > float64EqualityThreshold
}
