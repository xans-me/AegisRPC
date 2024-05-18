package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/xans-me/AegisRPC/server"
	"google.golang.org/grpc"
)

func main() {
	addr := "localhost:50051"
	lis := server.ListenGRPCServer(addr)
	grpcServer := server.NewGRPCServer(server.WithMaxConnectionIdle(10 * time.Minute))

	// Register your gRPC service here
	// pb.RegisterYourServiceServer(grpcServer, &yourService{})

	go func() {
		log.Infof("Starting gRPC server on %s", addr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Handle graceful shutdown
	waitForShutdown(grpcServer)
}

func waitForShutdown(grpcServer *grpc.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Info("Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Info("gRPC server stopped")
}