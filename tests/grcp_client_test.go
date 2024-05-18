package tests

import (
	"google.golang.org/grpc"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xans-me/AegisRPC/client"
)

// TestGRPCClientConnection tests the gRPC client connection with various configurations.
func TestGRPCClientConnection(t *testing.T) {
	addr := "localhost:50051"

	// Test connection with default options
	t.Run("DefaultOptions", func(t *testing.T) {
		conn, err := connectWithRetries(addr, client.WithLogLevel(client.LogLevelDetailed))
		if err != nil {
			t.Fatalf("Failed to connect to gRPC server with default options: %v", err)
		}
		conn.Close()
	})

	// Test connection without retry mechanism
	t.Run("NoRetries", func(t *testing.T) {
		conn, err := connectWithRetries(addr, client.WithMaxRetries(0), client.WithLogLevel(client.LogLevelDetailed))
		if err != nil {
			t.Fatalf("Failed to connect to gRPC server without retry mechanism: %v", err)
		}
		conn.Close()
	})

	// Test connection with custom retry policy
	t.Run("CustomRetryPolicy", func(t *testing.T) {
		conn, err := connectWithRetries(addr, client.WithRetryPolicyJSON(`{
			"methodConfig": [{
				"name": [{"service": "custom.Service"}],
				"waitForReady": true,
				"retryPolicy": {
					"MaxAttempts": 10,
					"InitialBackoff": "1s",
					"MaxBackoff": "15s",
					"BackoffMultiplier": 2.0,
					"RetryableStatusCodes": ["UNKNOWN", "UNAVAILABLE"]
				}
			}]
		}`), client.WithLogLevel(client.LogLevelDetailed))
		if err != nil {
			t.Fatalf("Failed to connect to gRPC server with custom retry policy: %v", err)
		}
		conn.Close()
	})
}

// TestGRPCClientTimeout tests the connection timeout setting.
func TestGRPCClientTimeout(t *testing.T) {
	addr := "localhost:50051"

	// Test connection with custom timeout
	t.Run("CustomTimeout", func(t *testing.T) {
		conn, err := connectWithRetries(addr, client.WithConnectTimeout(5*time.Second), client.WithLogLevel(client.LogLevelDetailed))
		if err != nil {
			t.Fatalf("Failed to connect to gRPC server with custom timeout: %v", err)
		}
		conn.Close()
	})
}

// connectWithRetries tries to connect with retries and returns the connection or an error.
func connectWithRetries(addr string, opts ...client.Option) (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	var err error
	for i := 0; i < 5; i++ {
		conn = client.NewGRPCClient(addr, opts...)
		if conn != nil {
			return conn, nil
		}
		logrus.Warnf("Failed to connect, retrying... (%d/5)", i+1)
		time.Sleep(1 * time.Second)
	}
	return nil, err
}