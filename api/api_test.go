package api

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/alienxp03/spectral/api/service"
	"github.com/alienxp03/spectral/sqlite"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/status"
)

func TestGetUsages(t *testing.T) {
	tests := []struct {
		name    string
		db      *mockDB
		start   string
		end     string
		want    *service.GetUsageResponse
		wantErr string
	}{
		{
			name:  "OK",
			start: "2022-01-01 00:00:00",
			end:   "2022-01-01 01:00:00",
			db: &mockDB{
				getUsageHistory: func(t *testing.T, startTime, endTime time.Time) ([]sqlite.Usage, error) {
					return []sqlite.Usage{
						{
							Time:  time.Date(2022, 1, 1, 1, 0, 0, 0, time.UTC),
							Usage: 100,
						},
						{
							Time:  time.Date(2022, 1, 1, 2, 0, 0, 0, time.UTC),
							Usage: 200,
						},
					}, nil
				},
			},
			want: &service.GetUsageResponse{
				Data: &service.UsageData{
					Total: float64(300),
					Usages: []*service.Usage{
						{
							Time:  "Sat, 01 Jan 2022 01:00:00 UTC",
							Usage: float32(100),
						},
						{
							Time:  "Sat, 01 Jan 2022 02:00:00 UTC",
							Usage: float32(200),
						},
					},
				},
			},
		},
		{
			name:    "invalid start time",
			start:   "2022-33-31 00:00:00",
			end:     "2022-01-01 01:00:00",
			wantErr: status.Errorf(http.StatusBadRequest, "invalid request: invalid start time").Error(),
		},
		{
			name:    "invalid end time",
			start:   "2022-01-01 01:00:00",
			end:     "2022-33-31 00:00:00",
			wantErr: status.Errorf(http.StatusBadRequest, "invalid request: invalid end time").Error(),
		},
		{
			name:    "period is longer than allowed",
			start:   "2022-01-01 01:00:00",
			end:     "2022-04-30 00:00:00",
			wantErr: status.Errorf(http.StatusBadRequest, "invalid request: only periods of up to 90 days are allowed").Error(),
		},
		{
			name:  "db error",
			start: "2022-01-01 00:00:00",
			end:   "2022-01-01 01:00:00",
			db: &mockDB{
				getUsageHistory: func(t *testing.T, startTime, endTime time.Time) ([]sqlite.Usage, error) {
					return []sqlite.Usage{}, errors.New("db error")
				},
			},
			wantErr: status.Errorf(http.StatusInternalServerError, "failed to get usages: db error").Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewAPI(tt.db)

			resp, err := api.GetUsages(context.Background(), &service.GetUsageRequest{
				StartTime: tt.start,
				EndTime:   tt.end,
			})

			if err != nil {
				assert.Equal(t, tt.wantErr, err.Error())
			}

			if tt.want != nil {
				assert.Equal(t, tt.want.Data.Total, resp.Data.Total)
				assert.Equal(t, tt.want.Data.Usages, resp.Data.Usages)
			}
		})
	}
}

type mockDB struct {
	T               *testing.T
	getUsageHistory func(t *testing.T, startTime, endTime time.Time) ([]sqlite.Usage, error)
}

func (db *mockDB) GetUsageHistory(_ context.Context, startTime, endTime time.Time) ([]sqlite.Usage, error) {
	return db.getUsageHistory(db.T, startTime, endTime)
}
