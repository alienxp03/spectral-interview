package main

import (
	"context"
	"log"
	"time"

	"github.com/alienxp03/spectral/api/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := service.NewEnergyServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetUsages(ctx, &service.GetUsageRequest{
		StartTime: "2019-01-01 00:00:00",
		EndTime:   "2019-01-01 00:20:00",
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %+v", r.Data)
}
