package auth

import (
	"github.com/isaiahwong/auth-go/internal/store"
	"github.com/isaiahwong/auth-go/internal/util/log"
	"google.golang.org/grpc"
)

type serviceOption struct {
	logger     log.Logger
	store      *store.MongoStore
	grpcServer *grpc.Server
	production bool
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

func SetEnvironment(p bool) ServiceOption {
	return func(o *serviceOption) {
		o.production = p
	}
}
