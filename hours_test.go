package main

import (
	"reflect"
	"testing"
)

func Test_refillMissedHours(t *testing.T) {
	type args struct {
		today   float64
		entries []fillItem
	}
	tests := []struct {
		name    string
		args    args
		want    []fillItem
		wantErr bool
	}{
		{
			name: "already logged zero and two empty tasks",
			args: args{
				entries: []fillItem{
					{},
					{},
				},
			},
			want: []fillItem{
				{
					hours: "4.00",
				},
				{
					hours: "4.00",
				},
			},
		},
		{
			name: "exact 8 hours",
			args: args{
				entries: []fillItem{
					{
						hours: "1",
					},
					{
						hours: "7",
					},
				},
			},
			want: []fillItem{
				{
					hours: "1",
				},
				{
					hours: "7",
				},
			},
		},
		{
			name: "already logged 5 and two empty tasks",
			args: args{
				entries: []fillItem{
					{},
					{},
				},
				today: 5,
			},
			want: []fillItem{
				{
					hours: "1.50",
				},
				{
					hours: "1.50",
				},
			},
		},
		{
			name: "already logged 5 and one filled and one empty tasks",
			args: args{
				entries: []fillItem{
					{
						hours: "2",
					},
					{},
				},
				today: 5,
			},
			want: []fillItem{
				{
					hours: "2",
				},
				{
					hours: "1.00",
				},
			},
		},
		{
			name: "already logged more 8 hours ",
			args: args{
				entries: []fillItem{
					{},
					{},
				},
				today: 9,
			},
			want: []fillItem{
				{},
				{},
			},
			wantErr: true,
		},
		{
			name: "wrong hours variable",
			args: args{
				entries: []fillItem{
					{
						hours: "a",
					},
				},
			},
			want: []fillItem{
				{
					hours: "a",
				},
			},
			wantErr: true,
		},
		{
			name: "filled hours exceed daily goal",
			args: args{
				entries: []fillItem{
					{
						hours: "5",
					},
					{
						hours: "5",
					},
				},
			},
			want: []fillItem{
				{
					hours: "5",
				},
				{
					hours: "5",
				},
			},
			wantErr: true,
		},
		{
			name: "filled hours with logged hours exceed daily goal",
			args: args{
				entries: []fillItem{
					{
						hours: "1",
					},
					{
						hours: "1",
					},
				},
				today: 8,
			},
			want: []fillItem{
				{
					hours: "1",
				},
				{
					hours: "1",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := refillMissedHours(tt.args.entries, tt.args.today)
			if (err != nil) != tt.wantErr {
				t.Errorf("refillMissedHours() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("refillMissedHours() = %v, want %v", got, tt.want)
			}
		})
	}
}
