package server

import (
	"errors"
	"net"
	"reflect"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/isaiahwong/accounts-go/internal/store"
	"github.com/isaiahwong/accounts-go/internal/util/log"
)

// Server encapsulates authentication and user operations
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
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(l)),
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

	return server, nil
}
