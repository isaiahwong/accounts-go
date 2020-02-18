package auth

import (
	"github.com/isaiahwong/auth-go/internal/store"
	"github.com/isaiahwong/auth-go/internal/util/log"
	pb "github.com/isaiahwong/auth-go/protogen/auth/v1"
	"google.golang.org/grpc"
)

// Service defines the logic for authentication
type Service struct {
	logger log.Logger
	store  *store.MongoStore
}

type serviceOption struct {
	logger     log.Logger
	store      *store.MongoStore
	grpcServer *grpc.Server
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

// RegisterService takes in arguments notably grpcServer which is needed to register for
// protobuf service
func RegisterService(opt ...ServiceOption) error {
	opts := defaultServiceOption
	for _, o := range opt {
		o(&opts)
	}
	if opts.grpcServer == nil {
		return &InvalidParam{"grpcServer is nil, RegisterService requires type *grpc.Server"}
	}
	svc := &Service{
		logger: opts.logger,
	}
	// Register AuthService
	pb.RegisterAuthServiceServer(opts.grpcServer, svc)
	return nil
}
