package sqlite

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetUsageHistory(t *testing.T) {
	tests := []struct {
		name      string
		startTime time.Time
		endTime   time.Time
		setup     func(*Sqlite)
		want      []Usage
		wantErr   string
	}{
		{
			name:      "OK",
			startTime: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			endTime:   time.Date(2022, 1, 1, 2, 0, 0, 0, time.UTC),
			setup: func(db *Sqlite) {
				usages := []Usage{
					{
						ID:    1,
						Time:  time.Date(2022, 1, 1, 1, 0, 0, 0, time.UTC),
						Usage: 100,
					},
					{
						ID:    2,
						Time:  time.Date(2022, 1, 1, 2, 0, 0, 0, time.UTC),
						Usage: 200,
					},
				}

				for _, usage := range usages {
					_, err := db.DB.Exec("INSERT INTO meter_usages (id, time, usage) VALUES (?, ?, ?)", usage.ID, usage.Time, usage.Usage)
					if err != nil {
						t.Fatalf("failed to insert usage: %v", err)
					}
				}
			},
			want: []Usage{
				{
					ID:    2,
					Time:  time.Date(2022, 1, 1, 2, 0, 0, 0, time.UTC),
					Usage: 200,
				},
				{
					ID:    1,
					Time:  time.Date(2022, 1, 1, 1, 0, 0, 0, time.UTC),
					Usage: 100,
				},
			},
		},
		{
			name:      "Empty",
			startTime: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			endTime:   time.Date(2022, 1, 1, 2, 0, 0, 0, time.UTC),
			setup: func(db *Sqlite) {
			},
			want: []Usage{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := connect(t)
			tt.setup(db)

			got, err := db.GetUsageHistory(context.Background(), tt.startTime, tt.endTime)

			if err != nil {
				if tt.wantErr == "" {
					t.Errorf("unexpected error: %v", err)
				}
				if !(err.Error() == tt.wantErr) {
					t.Errorf("unexpected error: %v, want %v", err, tt.wantErr)
				}
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSeed(t *testing.T) {
	tests := []struct {
		name  string
		check func(*Sqlite)
	}{
		{
			name: "OK",
			check: func(db *Sqlite) {
				row, err := db.DB.Query("SELECT COUNT(*) FROM meter_usages")
				if err != nil {
					t.Fatalf("failed to select meter_usages: %v", err)
				}

				count := 0
				row.Next()
				_ = row.Scan(&count)
				if count == 0 {
					t.Fatalf("failed to seed: meter_usages is empty")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := connect(t)

			err := Seed(db)

			if err != nil {
				t.Fatalf("failed to seed: %v", err)
			}

			tt.check(db)
		})
	}
}

func connect(t *testing.T) *Sqlite {
	t.Helper()
	db, err := Connect("data_test")
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}

	_, err = db.DB.Exec("DELETE FROM meter_usages")
	if err != nil {
		t.Fatalf("failed to delete meter_usages: %v", err)
	}

	return db
}
