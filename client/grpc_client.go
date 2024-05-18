package client

import (
	"context"
	"fmt"
	"math"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

// NewGRPCClient creates and returns a new gRPC client connection using the specified options.
func NewGRPCClient(addr string, opts ...Option) *grpc.ClientConn {
	options := defaultClientOptions()
	for _, opt := range opts {
		opt(options)
	}

	if options.logLevel == LogLevelNone {
		log.SetLevel(log.WarnLevel)
	} else if options.logLevel == LogLevelSimple {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}

	log.Infof("Attempting to connect to gRPC server at %s", addr)

	grpcOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(loggingUnaryClientInterceptor(options.logLevel)),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoff.DefaultConfig,
			MinConnectTimeout: options.connectTimeout,
		}),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                options.keepaliveTime,
			Timeout:             options.keepaliveTimeout,
			PermitWithoutStream: true,
		}),
	}

	if options.enableRetry {
		grpcOpts = append(grpcOpts, grpc.WithDefaultServiceConfig(options.retryPolicyJSON))
	}

	retryCount := 0
	for {
		conn, err := grpc.Dial(addr, grpcOpts...)
		if err != nil {
			if !options.enableRetry {
				log.Fatalf("Failed to connect: %v. Exiting...", err)
			}
			if retryCount >= options.maxRetries {
				log.Fatalf("Could not connect to gRPC server at %s after %d attempts. Exiting...", addr, options.maxRetries)
			}
			log.Warnf("Failed to connect: %v. Retrying...", err)
			time.Sleep(getExponentialBackoffDuration(retryCount, options))
			retryCount++
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), options.connectTimeout)
		defer cancel()

		if waitForConnectionReady(conn, ctx) {
			log.Infof("Successfully connected to gRPC server at %s", addr)
			return conn
		}

		if !options.enableRetry {
			log.Fatalf("Failed to connect: %v. Exiting...", err)
		}
		if retryCount >= options.maxRetries {
			log.Fatalf("Could not connect to gRPC server at %s after %d attempts. Exiting...", addr, options.maxRetries)
		}
		backoffDuration := getExponentialBackoffDuration(retryCount, options)
		log.Warnf("Connection to gRPC server at %s was not ready. Retrying in %v...", addr, backoffDuration)
		conn.Close()
		time.Sleep(backoffDuration)
		retryCount++
	}
}

// loggingUnaryClientInterceptor logs each unary RPC call with details about the retry attempts.
func loggingUnaryClientInterceptor(level LogLevel) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if level != LogLevelNone {
			log.Infof("Calling method: %s", method)
		}
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				log.Warnf("RPC failed with status: %s, code: %s", st.Message(), st.Code())
			} else {
				log.Warnf("RPC failed with error: %v", err)
			}
		}
		return err
	}
}

// waitForConnectionReady waits for the gRPC connection to be ready.
func waitForConnectionReady(conn *grpc.ClientConn, ctx context.Context) bool {
	for {
		state := conn.GetState()
		if state == connectivity.Ready {
			return true
		}
		if !conn.WaitForStateChange(ctx, state) {
			return false
		}
	}
}

// getExponentialBackoffDuration calculates the backoff duration for retries.
func getExponentialBackoffDuration(retryCount int, options *ClientOptions) time.Duration {
	backoffDuration := time.Duration(float64(options.initialBackoff) * math.Pow(options.backoffMultiplier, float64(retryCount)))
	if backoffDuration > options.maxBackoff {
		backoffDuration = options.maxBackoff
	}
	return backoffDuration
}

