package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/alienxp03/spectral/api/service"
	"github.com/alienxp03/spectral/sqlite"
	"google.golang.org/grpc/status"
)

type DB interface {
	GetUsageHistory(ctx context.Context, startTime, endTime time.Time) ([]sqlite.Usage, error)
}

type API struct {
	service.UnimplementedEnergyServiceServer
	Logger *slog.Logger
	DB     DB
}

const (
	maxPeriod = 90
)

func NewAPI(db DB) *API {
	return &API{
		Logger: slog.Default(),
		DB:     db,
	}
}

type usageRequest struct {
	StartTime time.Time
	EndTime   time.Time
}

func (a *API) GetUsages(ctx context.Context, req *service.GetUsageRequest) (*service.GetUsageResponse, error) {
	params, err := parseRequest(req)
	if err != nil {
		return nil, status.Errorf(http.StatusBadRequest, "invalid request: %v", err)
	}

	records, err := a.DB.GetUsageHistory(ctx, params.StartTime, params.EndTime)
	if err != nil {
		return nil, status.Errorf(http.StatusInternalServerError, "failed to get usages: %v", err)
	}

	usages := make([]*service.Usage, 0)
	total := float64(0)

	for _, usage := range records {
		total += usage.Usage
		usages = append(usages, &service.Usage{
			Time:  usage.Time.UTC().Format(time.RFC1123),
			Usage: float32(usage.Usage),
		})
	}

	return &service.GetUsageResponse{
		Data: &service.UsageData{
			Total:  total,
			Usages: usages,
		},
	}, nil
}

func parseRequest(req *service.GetUsageRequest) (*usageRequest, error) {
	startTime, err := time.Parse("2006-01-02 15:04:05", req.StartTime)
	if err != nil {
		return nil, errors.New("invalid start time")
	}

	endTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime)
	if err != nil {
		return nil, errors.New("invalid end time")
	}

	days := endTime.Sub(startTime).Hours() / 24
	if days > maxPeriod {
		return nil, fmt.Errorf("only periods of up to %d days are allowed", maxPeriod)
	}

	return &usageRequest{
		StartTime: startTime,
		EndTime:   endTime,
	}, nil
}
