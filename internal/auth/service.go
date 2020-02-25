package auth

import (
	"github.com/go-playground/validator/v10"
	"github.com/isaiahwong/auth-go/internal"
	"github.com/isaiahwong/auth-go/internal/store"
	"github.com/isaiahwong/auth-go/internal/store/drivers/mongo"
	user "github.com/isaiahwong/auth-go/internal/store/repo/user"
	"github.com/isaiahwong/auth-go/internal/util/email"
	"github.com/isaiahwong/auth-go/internal/util/log"
	pb "github.com/isaiahwong/auth-go/protogen/auth/v1"
	"github.com/microcosm-cc/bluemonday"
)

// Service defines the logic for authentication
type Service struct {
	production bool
	logger     log.Logger
	policy     *bluemonday.Policy
	validate   *validator.Validate
	userRepo   user.Repo

	recaptchaURL    string
	recaptchaSecret string
}

func (svc *Service) initRepoWithMongo(s store.DataStore) error {
	m, ok := s.(*mongo.MongoStore)
	if !ok {
		return &internal.InvalidParam{S: "Invalid Type. Only MongoStore is supported at this time"}
	}
	svc.userRepo = user.NewMongoUserRepo(m)
	return nil
}

func (svc *Service) initValidator() {
	svc.validate = validator.New()
	svc.validate.RegisterValidation("emailMX", func(fl validator.FieldLevel) bool {
		f := fl.Field().String()
		return email.ValidateFormat(f) && email.ValidateHost(f)
	})
}

// RegisterService takes in arguments notably grpcServer which is needed to register for
// protobuf service
func RegisterService(opt ...ServiceOption) error {
	opts := defaultServiceOption
	for _, o := range opt {
		o(&opts)
	}
	if opts.grpcServer == nil {
		return &internal.InvalidParam{S: "grpcServer is nil, RegisterService requires type *grpc.Server"}
	}
	if opts.store == nil {
		return &internal.InvalidParam{S: "store is nil, RegisterService requires type *store.Datastore."}
	}
	svc := &Service{
		production: opts.production,
		logger:     opts.logger,
		policy:     bluemonday.StrictPolicy(),
	}
	svc.initValidator()
	// Initializes repositories
	if err := svc.initRepoWithMongo(opts.store); err != nil {
		return err
	}
	// Register AuthService
	pb.RegisterAuthServiceServer(opts.grpcServer, svc)
	return nil
}
