package server

import (
	"github.com/isaiahwong/accounts-go/internal/store"
	"github.com/isaiahwong/accounts-go/internal/util/log"
)

type serverOptions struct {
	name       string
	address    string
	production bool
	logger     log.Logger
	store      store.DataStore
}

// Option is an option that can be given to a Server on construction.
type Option func(*serverOptions)

var defaultServerOptions = serverOptions{
	address:    ":50051",
	name:       "Server",
	production: false,
	logger:     log.NewLogrusLogger(),
}

// WithLogger an Option that will set the internal
// logging of the server
func WithLogger(l log.Logger) Option {
	return func(o *serverOptions) {
		o.logger = l
	}
}

// WithName an Option that will set the server's name.
// This is merely for cosmetics purposes
func WithName(s string) Option {
	return func(o *serverOptions) {
		o.name = s
	}
}

// WithAddress an Option that will set the grpc server's address
func WithAddress(a string) Option {
	return func(o *serverOptions) {
		o.address = a
	}
}

// WithAppEnvironment an Option that will set determine if the application
// run in production
func WithAppEnvironment(production bool) Option {
	return func(o *serverOptions) {
		o.production = production
	}
}

// WithDataStore an Option which sets the store the server will use
// Initially the server will assist in overseeing the connection between
// the data store
func WithDataStore(s store.DataStore) Option {
	return func(o *serverOptions) {
		o.store = s
	}
}
