package accounts

import (
	"errors"

	"github.com/go-playground/validator/v10"
	accountsV1 "github.com/isaiahwong/accounts-go/api/accounts/v1"
	"github.com/isaiahwong/accounts-go/api/client"
	mailV1 "github.com/isaiahwong/accounts-go/api/mail/v1"
	"github.com/isaiahwong/accounts-go/internal/oauth"
	"github.com/isaiahwong/accounts-go/internal/store"
	"github.com/isaiahwong/accounts-go/internal/store/drivers/mongo"
	user "github.com/isaiahwong/accounts-go/internal/store/repo/user"
	"github.com/isaiahwong/accounts-go/internal/util"
	"github.com/isaiahwong/accounts-go/internal/util/email"
	"github.com/isaiahwong/accounts-go/internal/util/log"
	"github.com/microcosm-cc/bluemonday"
)

// Service defines the logic for authentication
type Service struct {
	production      bool
	test            bool
	logger          log.Logger
	policy          *bluemonday.Policy
	validate        *validator.Validate
	recaptchaURL    string
	recaptchaSecret string
	userRepo        user.Repo
	oauthClient     *oauth.Hydra
	mailSVC         mailV1.MailServiceClient
}

func (svc *Service) initRepoWithMongo(s store.DataStore) error {
	m, ok := s.(*mongo.MongoStore)
	if !ok {
		return errors.New("Invalid Type. Only MongoStore is supported at this time")
	}
	svc.userRepo = user.NewMongoUserRepo(m)
	return nil
}

func (svc *Service) initValidator() {
	svc.validate = validator.New()
	svc.validate.RegisterValidation("emailMX", func(fl validator.FieldLevel) bool {
		var valid bool
		f := fl.Field().String()
		if svc.production {
			valid = email.ValidateFormat(f) && email.ValidateHost(f)
		} else {
			valid = email.ValidateFormat(f)
		}
		return valid
	})
}

func (svc *Service) initServices() error {
	mailSVC, err := mailV1.NewMailClient(
		client.WithAddress(util.MapEnvWithDefaults("MAIL_SERVICE", ":50051")),
	)
	if err != nil {
		return err
	}
	svc.mailSVC = mailSVC
	return nil
}

func initServices() error {
	return nil
}

// RegisterService takes in arguments notably grpcServer which is needed to register for
// protobuf Service
func RegisterService(opt ...ServiceOption) error {
	opts := defaultServiceOption
	for _, o := range opt {
		o(&opts)
	}
	if !opts.disableServer && opts.grpcServer == nil {
		return errors.New("auth: grpcServer is nil. RegisterService requires type *grpc.Server")
	}
	if !opts.disableStore && opts.store == nil {
		return errors.New("auth: store is nil. RegisterService requires type *store.Datastore")
	}
	svc := &Service{
		production:  opts.production,
		logger:      opts.logger,
		policy:      bluemonday.StrictPolicy(),
		oauthClient: oauth.NewHydraClient(),
	}
	svc.initValidator()
	svc.initServices()

	// Initializes repositories
	if err := svc.initRepoWithMongo(opts.store); err != nil {
		return err
	}
	// Register AuthService
	accountsV1.RegisterAccountsServiceServer(opts.grpcServer, svc)
	return nil
}
