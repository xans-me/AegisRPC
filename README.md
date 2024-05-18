# AegisRPC

AegisRPC is a flexible and robust gRPC client and server library that provides enhanced features such as configurable retry mechanisms, exponential backoff, customizable logging levels, and easy integration with your gRPC services. The library is designed to make gRPC integration smoother and more reliable for both client and server applications.

## Why AegisRPC?

AegisRPC is built on top of the powerful gRPC-Go library. While gRPC-Go is a fantastic library for creating high-performance RPC (Remote Procedure Call) applications, AegisRPC extends its functionality by offering:

- **Configurable Retry Mechanisms**: Easily set up retries with exponential backoff and custom retry policies.
- **Enhanced Logging**: Choose from various logging levels to get the right amount of information you need.
- **Flexible Options**: Customize client and server options to fit your specific needs.
- **Ease of Use**: Simplified client and server creation with sensible defaults.

## Installation

To install AegisRPC, use the following command:

```sh
go install github.com/xans-me/AegisRPC
```

## Usage

### Client Example
Here's an example of how to use the AegisRPC client:

```go
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
```

### Server Example
Here's an example of how to set up an AegisRPC server:

```go
package main

import (
	"os"
	"os/signal"
	"syscall"

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
```

## Configuration Options
### Client Options
- **WithInitialBackoff**: Sets the initial backoff duration for retries.
- **WithMaxBackoff**: Sets the maximum backoff duration for retries.
- **WithBackoffMultiplier**: Sets the backoff multiplier for retries.
- **WithMaxRetries**: Sets the maximum number of retries.
- **WithConnectTimeout**: Sets the connection timeout.
- **WithKeepaliveTime**: Sets the keepalive time duration.
- **WithKeepaliveTimeout**: Sets the keepalive timeout duration.
- **WithLogLevel**: Sets the logging level.
- **WithRetryPolicyJSON**: Sets the retry policy JSON.

### Server Options
- **WithMaxConnectionIdle**: Sets the maximum connection idle time.
- **WithKeepaliveTime**: Sets the keepalive time duration.
- **WithKeepaliveTimeout**: Sets the keepalive timeout duration.

##Contributing
Contributions are welcome! Please open an issue or submit a pull request with your changes.

## License
This project is licensed under the MIT License. See the LICENSE file for details.
```
This `README.md` provides clear guidance on using the library, its benefits, and the reasons for its creation. It also includes installation instructions, example code for both client and server usage, and detailed descriptions of the configurable options.
```