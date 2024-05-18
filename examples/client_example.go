package main

import (
	"fmt"
	"time"

	"github.com/xans-me/AegisRPC/client"
)

func main() {
	addr := "localhost:50051"
	conn := client.NewGRPCClient(addr, 
		client.WithLogLevel(client.LogLevelSimple),
		client.WithMaxRetries(0)) // Disable retry
	defer conn.Close()

	// Example usage of the gRPC client connection
	fmt.Println("Connected to gRPC server:", conn.Target())

	// Simulate some client-side work
	time.Sleep(2 * time.Second)
}

