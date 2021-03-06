package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
)

type fillItem struct {
	task    string
	hours   string
	comment string
}

const itemRowCount = 3

func parseCSV(file io.Reader) ([]fillItem, error) {
	result := make([]fillItem, 0)
	reader := csv.NewReader(file)

	for line := 0; ; line++ {
		record, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return result, nil
			}
			return nil, err
		}
		if len(record) < itemRowCount {
			return nil, fmt.Errorf("parseCSV: parsing line %d failed: line should have more 3 items", line)
		}
		item := fillItem{
			task:    record[0],
			hours:   record[1],
			comment: record[2],
		}

		if item.task == "" {
			return nil, errors.New("parseCSV: task id should not be empty")
		}

		result = append(result, item)
	}
}
