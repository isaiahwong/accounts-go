package client

import (
	"time"

	"google.golang.org/grpc"
)

type clientOption struct {
	address string
	timeout time.Duration
}

var defaultOptions clientOption = clientOption{
	address: ":50051",
	timeout: 10 * time.Second,
}

// Option defines the APIClient options
type Option func(*clientOption)

// WithAddress returns an Option which sets the APIClient remote address
func WithAddress(a string) Option {
	return func(o *clientOption) {
		o.address = a
	}
}

// WithTimeout returns an option which sets the connection timeout on connection initially
func WithTimeout(t time.Duration) Option {
	return func(o *clientOption) {
		o.timeout = t
	}
}

// CreateClient returns a Client Connection
func CreateClient(opt ...Option) (*grpc.ClientConn, error) {
	var opts = &defaultOptions
	for _, o := range opt {
		o(opts)
	}

	conn, err := grpc.Dial(
		opts.address,
		grpc.WithInsecure(),
		grpc.WithTimeout(opts.timeout),
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
