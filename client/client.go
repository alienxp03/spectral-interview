package client

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/alienxp03/spectral/api/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func GetUsages(startTime, endTime string) (*service.GetUsageResponse, error) {
	conn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := service.NewEnergyServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var errMsg error
	usages, err := c.GetUsages(ctx, &service.GetUsageRequest{StartTime: startTime, EndTime: endTime})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return nil, err
		}
		errMsg = errors.New(st.Message())
	}

	return usages, errMsg
}
