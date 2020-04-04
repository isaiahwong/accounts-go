package server

import (
	"context"
	"errors"
	"net"
	"reflect"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	healthV1 "github.com/isaiahwong/accounts-go/api/health/v1"
	"github.com/isaiahwong/accounts-go/internal/common/log"
	"github.com/isaiahwong/accounts-go/internal/store"
)

// Server encapsulates authentication and account operations
type Server struct {
	GRPCServer *grpc.Server
	Name       string
	Production bool
	listener   net.Listener
	logger     log.Logger
	store      store.DataStore
}

// Serve starts gRPC server as well as other dependencies such as connect to store
func (s *Server) Serve() error {
	if s.store == nil {
		return errors.New("Serve: requires a Datastore to start server")
	}
	if err := s.store.Connect(nil); err != nil {
		return err
	}

	s.logger.Infof("Serving %v on %v %v", s.Name, s.listener.Addr().Network(), s.listener.Addr().String())
	s.logger.Infof("Production: %v", s.Production)
	if err := s.GRPCServer.Serve(s.listener); err != nil {
		return err
	}
	return nil
}

func LoggerDecider(method string, err error) bool {
	exclude := []string{
		"/grpc.health.v1.Health/Check",
	}
	if err != nil {
		return true
	}
	for _, e := range exclude {
		if strings.Contains(method, e) {
			return false
		}
	}
	return true
}

// New returns a new Server
func New(opt ...Option) (*Server, error) {
	opts := defaultServerOptions
	// Apply passed in options
	for _, o := range opt {
		o(&opts)
	}
	l := &logrus.Logger{}

	// We set the value of opts.logger to interface type logrus.Logger
	// In this case, we only support logrus.Logger thus no error checks
	// are performed.
	// TODO: Refactor to check type of logger before setting to the
	// respective interface
	func(i, o interface{}) {
		reflect.ValueOf(o).Elem().Set(reflect.ValueOf(i))
	}(opts.logger, &l)

	// Create a new gRPC server
	gs := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(l), grpc_logrus.WithDecider(LoggerDecider)),
		)),
	)

	// Create a new network listener
	lis, err := net.Listen("tcp", opts.address)
	if err != nil {
		return nil, err
	}

	server := &Server{
		GRPCServer: gs,
		listener:   lis,
		store:      opts.store,
		logger:     opts.logger,
		Production: opts.production,
		Name:       opts.name,
	}

	// Register HealthService
	healthV1.RegisterHealthServer(gs, server)

	return server, nil
}

// Check ensures all services the server is running is healthy
func (*Server) Check(ctx context.Context, req *healthV1.HealthCheckRequest) (*healthV1.HealthCheckResponse, error) {
	return &healthV1.HealthCheckResponse{Status: healthV1.HealthCheckResponse_SERVING}, nil
}
