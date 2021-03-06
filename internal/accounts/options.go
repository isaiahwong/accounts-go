package accounts

import (
	"github.com/isaiahwong/accounts-go/internal/common/log"
	"github.com/isaiahwong/accounts-go/internal/store"
	"google.golang.org/grpc"
)

type serviceOption struct {
	logger        log.Logger
	store         store.DataStore
	grpcServer    *grpc.Server
	disableStore  bool
	disableServer bool
	production    bool
}

// ServiceOption sets options
type ServiceOption func(*serviceOption)

var defaultServiceOption = serviceOption{
	logger: log.NewLogrusLogger(),
}

// WithLogger returns a ServiceOption that will set the internal
// logging of the server
func WithLogger(l log.Logger) ServiceOption {
	return func(o *serviceOption) {
		o.logger = l
	}
}

// WithGrpc returns a ServiceOption that will set the gRPC server
func WithGrpc(g *grpc.Server) ServiceOption {
	return func(o *serviceOption) {
		o.grpcServer = g
	}
}

// WithStore returns a ServiceOption that sets the store the service will use
func WithStore(store store.DataStore) ServiceOption {
	return func(o *serviceOption) {
		o.store = store
	}
}

// SetEnvironment returns a ServiceOption that sets the service environment
func SetEnvironment(production bool) ServiceOption {
	return func(o *serviceOption) {
		o.production = production
	}
}

func DisableStore(s bool) ServiceOption {
	return func(o *serviceOption) {
		o.disableStore = s
	}
}

func DisableServer(s bool) ServiceOption {
	return func(o *serviceOption) {
		o.disableServer = s
	}
}
