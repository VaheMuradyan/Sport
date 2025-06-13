package main

import (
	"github.com/VaheMuradyan/Sport/db"
	"github.com/VaheMuradyan/Sport/proto"
	"github.com/VaheMuradyan/Sport/services"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	database := db.ConnectDB()
	coefficientService := services.NewCoefficientService(database)
	centrifugoService := services.NewCentrifugoService("http://localhost:8000", "0957bfe1-5aa9-40c0-991f-d15150f91594")

	go startGRPCServer(coefficientService, centrifugoService)

	startHTTPServer()

}

func startGRPCServer(coefficientService *services.CoefficientService, centrifugoService *services.CentrifugoService) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterCoefficientServiceServer(grpcServer, services.NewGRPCCoefficientServer(coefficientService, centrifugoService))

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func startHTTPServer() {
	r := gin.Default()

	log.Println("HTTP server listening on :8080")
	r.Run(":8080")
}
