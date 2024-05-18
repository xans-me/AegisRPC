package server

import (
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// NewGRPCServer creates a new gRPC server with the specified options.
func NewGRPCServer(opts ...Option) *grpc.Server {
	options := defaultServerOptions()
	for _, opt := range opts {
		opt(options)
	}

	return grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: options.maxConnectionIdle,
		Time:              options.keepaliveTime,
		Timeout:           options.keepaliveTimeout,
	}))
}

// ListenGRPCServer starts listening on the given address for gRPC connections.
func ListenGRPCServer(addr string) net.Listener {
	log.Infof("Attempting to setup listener gRPC Transport on %s", addr)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Infof("Successfully opened listener gRPC Transport on %s", addr)
	return lis
}

