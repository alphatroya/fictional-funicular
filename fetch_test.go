package main

import (
	"errors"
	"testing"

	"github.com/alphatroya/redmine-helper-bot/redmine"

	m "github.com/alphatroya/fictional-funicular/mocks"
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
			mock := m.NewMockClient(ctrl)
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
