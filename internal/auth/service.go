package auth

import (
	"github.com/go-playground/validator/v10"
	"github.com/isaiahwong/auth-go/internal/store"
	"github.com/isaiahwong/auth-go/internal/util/log"
	pb "github.com/isaiahwong/auth-go/protogen/auth/v1"
	"github.com/microcosm-cc/bluemonday"
)

// Service defines the logic for authentication
type Service struct {
	production bool
	logger     log.Logger
	store      *store.MongoStore
	policy     *bluemonday.Policy
	validate   *validator.Validate

	recaptchaURL    string
	recaptchaSecret string
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
		production: opts.production,
		logger:     opts.logger,
		policy:     bluemonday.StrictPolicy(),
		validate:   validator.New(),
	}
	// Register AuthService
	pb.RegisterAuthServiceServer(opts.grpcServer, svc)
	return nil
}
