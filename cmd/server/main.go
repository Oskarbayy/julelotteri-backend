package main

import (
	"log"
	"net"

	julelotteri "github.com/oskarbayy/julelotteri-backend/generated"
	"github.com/oskarbayy/julelotteri-backend/internal/services"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen on port :8080, %v", err)
	}

	grpcServer := grpc.NewServer()
	lotteriService := &services.LotteriService{}

	julelotteri.RegisterLotteriServiceServer(grpcServer, lotteriService)

	print("gRPC server is running on port 8080")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port :8080 %v", err)
	}
}
