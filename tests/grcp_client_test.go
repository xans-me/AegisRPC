package tests

import (
	"testing"
	"time"

	"github.com/xans-me/AegisRPC/client"
	"google.golang.org/grpc"
)

// TestGRPCClientConnection tests the gRPC client connection with various configurations.
func TestGRPCClientConnection(t *testing.T) {
	addr := "localhost:50051"

	// Test connection with default options
	conn := client.NewGRPCClient(addr)
	if conn == nil {
		t.Fatal("Failed to connect to gRPC server with default options")
	}
	conn.Close()

	// Test connection without retry mechanism
	conn = client.NewGRPCClient(addr, client.WithMaxRetries(0))
	if conn == nil {
		t.Fatal("Failed to connect to gRPC server without retry mechanism")
	}
	conn.Close()

	// Test connection with custom retry policy
	conn = client.NewGRPCClient(addr, client.WithRetryPolicyJSON(`{
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
	}`))
	if conn == nil {
		t.Fatal("Failed to connect to gRPC server with custom retry policy")
	}
	conn.Close()
}

// TestGRPCClientTimeout tests the connection timeout setting.
func TestGRPCClientTimeout(t *testing.T) {
	addr := "localhost:50051"

	// Test connection with custom timeout
	conn := client.NewGRPCClient(addr, client.WithConnectTimeout(5*time.Second))
	if conn == nil {
		t.Fatal("Failed to connect to gRPC server with custom timeout")
	}
	conn.Close()
}

