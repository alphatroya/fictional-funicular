package main

import (
	"errors"
	"testing"

	"github.com/alphatroya/fictional-funicular/mocks"
	"github.com/alphatroya/redmine-helper-bot/redmine"
	"github.com/golang/mock/gomock"
)

func Test_fetchTodayTotalSum(t *testing.T) {
	type args struct {
		data []*redmine.TimeEntryResponse
		err  error
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "Success empty request",
			args: args{},
		},
		{
			name: "Success data single request",
			args: args{
				data: []*redmine.TimeEntryResponse{
					{
						Hours: 1,
					},
				},
			},
			want: 1,
		},
		{
			name: "Success data single request",
			args: args{
				data: []*redmine.TimeEntryResponse{
					{
						Hours: 1,
					},
					{
						Hours: 1.5,
					},
				},
			},
			want: 2.5,
		},
		{
			name: "Failed request",
			args: args{
				err: errors.New("Error"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mocks.NewMockClient(ctrl)
			mock.
				EXPECT().
				TodayTimeEntries().
				DoAndReturn(func() ([]*redmine.TimeEntryResponse, error) {
					return tt.args.data, tt.args.err
				}).
				AnyTimes()

			got, err := fetchTodayTotalSum(mock)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchTodayTotalSum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("fetchTodayTotalSum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeFillRequests(t *testing.T) {
	type args struct {
		failedTask string
		entries    []fillItem
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should make len(entries) successful requests",
			args: args{
				entries: []fillItem{
					{
						task:    "1",
						hours:   "1a",
						comment: "1b",
					},
					{
						task:    "2",
						hours:   "2a",
						comment: "2b",
					},
					{
						task:    "3",
						hours:   "3a",
						comment: "3b",
					},
				},
			},
		},
		{
			name: "Should return error when one request failed",
			args: args{
				failedTask: "2",
				entries: []fillItem{
					{
						task:    "1",
						hours:   "1a",
						comment: "1b",
					},
					{
						task:    "2",
						hours:   "2a",
						comment: "2b",
					},
					{
						task:    "3",
						hours:   "3a",
						comment: "3b",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Should return error when first request failed",
			args: args{
				failedTask: "1",
				entries: []fillItem{
					{
						task:    "1",
						hours:   "1a",
						comment: "1b",
					},
					{
						task:    "2",
						hours:   "2a",
						comment: "2b",
					},
					{
						task:    "3",
						hours:   "3a",
						comment: "3b",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mocks.NewMockClient(ctrl)
			for _, item := range tt.args.entries {
				mock.
					EXPECT().
					FillHoursRequest(item.task, item.hours, item.comment, "").
					DoAndReturn(func(id string, _, _, _ interface{}) (*redmine.TimeEntryBodyResponse, error) {
						if id == tt.args.failedTask {
							return nil, errors.New("Request error")
						}
						return nil, nil
					}).
					MaxTimes(1)
			}
			if err := makeFillRequests(mock, tt.args.entries); (err != nil) != tt.wantErr {
				t.Errorf("makeFillRequests() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
