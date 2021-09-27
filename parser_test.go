package main

import (
	"reflect"
	"strings"
	"testing"
)

func Test_parseCSV(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []fillItem
		wantErr bool
	}{
		{
			name:  "correct csv",
			input: "64421,,Test\n65345,1,Test 2",
			want: []fillItem{
				{
					task:    "64421",
					hours:   "",
					comment: "Test",
				},
				{
					task:    "65345",
					hours:   "1",
					comment: "Test 2",
				},
			},
		},
		{
			name:    "wrong number of rows",
			input:   "64421,",
			wantErr: true,
		},
		{
			name:    "wrong csv file",
			input:   "1,,Test\n2,",
			wantErr: true,
		},
		{
			name:    "valid csv file but empty issue id",
			input:   ",,Test",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			got, err := parseCSV(reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCSV() = %v, want %v", got, tt.want)
			}
		})
	}
}
