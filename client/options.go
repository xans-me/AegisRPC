package client

import "time"

// ClientOptions holds configuration options for the gRPC client.
type ClientOptions struct {
	initialBackoff    time.Duration
	maxBackoff        time.Duration
	backoffMultiplier float64
	maxRetries        int
	connectTimeout    time.Duration
	keepaliveTime     time.Duration
	keepaliveTimeout  time.Duration
	logLevel          LogLevel
	retryPolicyJSON   string
	enableRetry       bool
}

// LogLevel defines the level of logging.
type LogLevel int

const (
	LogLevelNone LogLevel = iota
	LogLevelSimple
	LogLevelDetailed
)

// Option is a function that configures a ClientOptions.
type Option func(*ClientOptions)

// WithInitialBackoff sets the initial backoff duration for retries.
func WithInitialBackoff(d time.Duration) Option {
	return func(o *ClientOptions) {
		o.initialBackoff = d
	}
}

// WithMaxBackoff sets the maximum backoff duration for retries.
func WithMaxBackoff(d time.Duration) Option {
	return func(o *ClientOptions) {
		o.maxBackoff = d
	}
}

// WithBackoffMultiplier sets the backoff multiplier for retries.
func WithBackoffMultiplier(m float64) Option {
	return func(o *ClientOptions) {
		o.backoffMultiplier = m
	}
}

// WithMaxRetries sets the maximum number of retries.
func WithMaxRetries(r int) Option {
	return func(o *ClientOptions) {
		o.maxRetries = r
		o.enableRetry = r > 0
	}
}

// WithConnectTimeout sets the connection timeout.
func WithConnectTimeout(d time.Duration) Option {
	return func(o *ClientOptions) {
		o.connectTimeout = d
	}
}

// WithKeepaliveTime sets the keepalive time duration.
func WithKeepaliveTime(d time.Duration) Option {
	return func(o *ClientOptions) {
		o.keepaliveTime = d
	}
}

// WithKeepaliveTimeout sets the keepalive timeout duration.
func WithKeepaliveTimeout(d time.Duration) Option {
	return func(o *ClientOptions) {
		o.keepaliveTimeout = d
	}
}

// WithLogLevel sets the logging level.
func WithLogLevel(level LogLevel) Option {
	return func(o *ClientOptions) {
		o.logLevel = level
	}
}

// WithRetryPolicyJSON sets the retry policy JSON.
func WithRetryPolicyJSON(json string) Option {
	return func(o *ClientOptions) {
		o.retryPolicyJSON = json
		o.enableRetry = json != ""
	}
}

// defaultClientOptions returns the default options for the client.
func defaultClientOptions() *ClientOptions {
	return &ClientOptions{
		initialBackoff:    500 * time.Millisecond,
		maxBackoff:        10 * time.Second,
		backoffMultiplier: 1.5,
		maxRetries:        5,
		connectTimeout:    10 * time.Second,
		keepaliveTime:     10 * time.Second,
		keepaliveTimeout:  20 * time.Second,
		logLevel:          LogLevelDetailed,
		retryPolicyJSON:   `{
			"methodConfig": [{
				"name": [{"service": "stock.StockService"}],
				"waitForReady": true,
				"retryPolicy": {
					"MaxAttempts": 5,
					"InitialBackoff": ".5s",
					"MaxBackoff": "10s",
					"BackoffMultiplier": 1.5,
					"RetryableStatusCodes": ["UNKNOWN", "UNAVAILABLE", "DEADLINE_EXCEEDED", "ABORTED", "RESOURCE_EXHAUSTED"]
				}
			}]
		}`,
		enableRetry: true,
	}
}

