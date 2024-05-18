package server

import "time"

// Options ServerOptions holds configuration options for the gRPC server.
type Options struct {
	maxConnectionIdle time.Duration
	keepaliveTime     time.Duration
	keepaliveTimeout  time.Duration
}

// Option is a function that configures a ServerOptions.
type Option func(*Options)

// WithMaxConnectionIdle sets the maximum connection idle time.
func WithMaxConnectionIdle(d time.Duration) Option {
	return func(o *Options) {
		o.maxConnectionIdle = d
	}
}

// WithKeepaliveTime sets the keepalive time duration.
func WithKeepaliveTime(d time.Duration) Option {
	return func(o *Options) {
		o.keepaliveTime = d
	}
}

// WithKeepaliveTimeout sets the keepalive timeout duration.
func WithKeepaliveTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.keepaliveTimeout = d
	}
}

// defaultServerOptions returns the default options for the server.
func defaultServerOptions() *Options {
	return &Options{
		maxConnectionIdle: 5 * time.Minute,
		keepaliveTime:     10 * time.Second,
		keepaliveTimeout:  20 * time.Second,
	}
}