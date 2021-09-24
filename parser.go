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
			return nil, fmt.Errorf("Parsing line %d failed. Line should have more 3 items", line)
		}
		result = append(result, fillItem{
			task:    record[0],
			hours:   record[1],
			comment: record[2],
		})
	}
}
