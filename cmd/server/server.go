package main

import (
	"flag"
	"log"
	"log/slog"
	"net"
	"os"

	"github.com/alienxp03/spectral/api"
	"github.com/alienxp03/spectral/api/service"
	"github.com/alienxp03/spectral/sqlite"
	"google.golang.org/grpc"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	_, dbNotFound := os.Stat("storage/data.db")

	db, err := sqlite.Connect("data")
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	if dbNotFound != nil {
		logger.Info("Seeding database for first time run")
		err = sqlite.Seed(db)
		if err != nil {
			log.Fatalf("failed to seed database: %v", err)
		}
	}

	addr := flag.String("addr", "localhost:9090", "HTTP network address")
	listner, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	service.RegisterEnergyServiceServer(s, api.NewAPI(db))

	logger.Info("Ready to accept traffic", "address", *addr)
	if err := s.Serve(listner); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
