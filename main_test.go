package main

import (
	"fmt"
	"testing"
	"time"
)

func Test_checkVacation(t *testing.T) {
	testCases := []struct {
		filename   string
		isVacation bool
		isErr      bool
	}{
		{
			filename:   "vacation-empty.txt",
			isVacation: false,
		},
		{
			filename:   "vacation-today.txt",
			isVacation: true,
		},
		{
			filename:   "vacation-not-today.txt",
			isVacation: false,
		},
		{
			filename:   "vacation-bad.txt",
			isVacation: false,
			isErr:      true,
		},
	}

	today, err := time.Parse("02.01.2006", "12.01.2022")
	if err != nil {
		t.Fatalf("wrong today date, err=%s", err)
	}
	for _, tc := range testCases {
		isVacation, err := checkVacation(fmt.Sprintf("testdata/%s", tc.filename), today)
		if isVacation != tc.isVacation {
			t.Errorf("Expected %v but got %v", tc.isVacation, isVacation)
		}

		if (err != nil) != tc.isErr {
			t.Errorf("Expected error result %v but got %v", tc.isErr, err)
		}
	}
}
