package accounts

import (
	"context"
	"fmt"
	"strings"
	"time"

	accountsV1 "github.com/isaiahwong/accounts-go/api/accounts/v1"
	"github.com/isaiahwong/accounts-go/internal/common"
	"github.com/isaiahwong/accounts-go/internal/common/recaptcha"
	"github.com/isaiahwong/accounts-go/internal/common/validator"
	"github.com/isaiahwong/accounts-go/internal/models"
	"github.com/isaiahwong/accounts-go/internal/oauth"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Headers Enum types
const (
	XForwardedFor    = "x-forwarded-for"
	CaptchaResponse  = "captcha-response"
	LoginChallenge   = "login-challenge"
	ConsentChallenge = "consent-challenge"
)

func (s *Service) LoginWithChallenge(ctx context.Context, _ *accountsV1.Empty) (*accountsV1.HydraResponse, error) {
	api := "LoginWithChallenge"

	challenge := common.GetMetadataValue(ctx, LoginChallenge)
	ip := common.GetMetadataValue(ctx, XForwardedFor)

	errors := validator.Val(
		s.validate,
		validator.Field{
			Param:   LoginChallenge,
			Message: LoginChallenge + " header required",
			Value:   challenge,
			Tag:     `required`,
		},
		validator.Field{
			Param:   XForwardedFor,
			Message: XForwardedFor + " header required",
			Value:   ip,
			Tag:     `required`,
		},
	)

	// Validate
	if len(errors) > 0 {
		return nil, s.returnErrors(ctx, errors, codes.InvalidArgument, "Invalid arguments", api)
	}
	// Prepend IP for logging
	api = fmt.Sprintf("[%v] %v", ip, api)

	resp, err := s.oAuthClient.Login(challenge)
	if err != nil {
		s.logger.Errorf("%v: %v", api, err)
		// Cast  to hydra error
		if he, ok := err.(*oauth.HydraError); ok {
			return nil, s.returnHydraError(ctx, he, api)
		}
		return nil, status.Error(codes.Internal, "An Internal error has occurred")
	}

	hr := &accountsV1.HydraResponse{
		Challenge:  resp.Challenge,
		RequestUrl: resp.RequestURL,
		SessionId:  resp.SessionID,
		Subject:    resp.Subject,
		Skip:       resp.Skip,
	}
	if !resp.Skip {
		return hr, nil
	}

	r, err := s.oAuthClient.AcceptLogin(challenge, &oauth.HydraLoginAccept{
		Subject: resp.Subject,
	})
	if err != nil {
		s.logger.Errorf("%v: %v", api, err)
		// Cast  to hydra error
		if he, ok := err.(*oauth.HydraError); ok {
			return nil, s.returnHydraError(ctx, he, api)
		}
		return nil, status.Error(codes.Internal, "An Internal error has occurred")
	}
	hr.RedirectTo = r.RedirectTo
	return hr, nil
}

func (s *Service) ConsentWithChallenge(ctx context.Context, req *accountsV1.Empty) (*accountsV1.RedirectResponse, error) {
	api := "ConsentWithChallenge"
	ie := status.Error(codes.Internal, "An Internal error has occurred")
	challenge := common.GetMetadataValue(ctx, ConsentChallenge)
	ip := common.GetMetadataValue(ctx, XForwardedFor)

	errors := validator.Val(
		s.validate,
		validator.Field{
			Param:   ConsentChallenge,
			Message: ConsentChallenge + " header required",
			Value:   challenge,
			Tag:     `required`,
		},
		validator.Field{
			Param:   XForwardedFor,
			Message: XForwardedFor + " header required",
			Value:   ip,
			Tag:     `required`,
		},
	)
	// Validate
	if len(errors) > 0 {
		return nil, s.returnErrors(ctx, errors, codes.InvalidArgument, "Invalid arguments", api)
	}
	// Prepend IP for logging
	api = fmt.Sprintf("[%v] %v", ip, api)

	resp, err := s.oAuthClient.Consent(challenge)
	if err != nil {
		s.logger.Errorf("%v: %v", api, err)
		return nil, ie
	}

	// Retrieve account
	u, err := s.findAccountByID(ctx, resp.Subject)
	if err != nil {
		s.logger.Errorf("%v: %v", api, err)
		return nil, ie
	}
	if u == nil {
		// Invalidate
		s.logger.Errorf("%v: No such account initiated consent request", api)
		return nil, ie
	}

	// Accept the scope for the accounts as we're authenticating internally
	r, err := s.oAuthClient.AcceptConsent(challenge, &oauth.HydraConsentAccept{
		GrantScope:               resp.RequestedScope,
		GrantAccessTokenAudience: resp.RequestedAccessTokenAudience,
		Remember:                 true,
		RememberFor:              0,
		Session: oauth.Session{
			IDToken: map[string]string{
				"email":       u.Auth.Email,
				"given_name":  u.Auth.FirstName,
				"family_name": u.Auth.LastName,
				"name":        u.Auth.Name,
				"picture":     u.Auth.Picture,
				// "email_verified": true,
				// "gender": "string",
				// "locale": "string",
				// "middle_name": "string",
				// "nickname": "string",
				// "phone_number": "string",
				// "phone_number_verified": true,
				// "picture": "string",
				// "preferred_username": "string",
				// "profile": "string",
				// "sub": "string",
				// "updated_at": 0,
				// "website": "string",
				// "zoneinfo": "string"
			},
		},
	})
	if err != nil {
		s.logger.Errorf("%v: %v", api, err)
		return nil, ie
	}
	return &accountsV1.RedirectResponse{RedirectTo: r.RedirectTo}, nil
}

func (s *Service) IsAuthenticated(ctx context.Context, in *accountsV1.Empty) (*accountsV1.AuthenticateResponse, error) {
	return nil, nil
}

// SignUp is a gRPC handler allows account to register
func (s *Service) SignUp(ctx context.Context, req *accountsV1.SignUpRequest) (*accountsV1.RedirectResponse, error) {
	api := "SignUp"

	ip := common.GetMetadataValue(ctx, XForwardedFor)
	captchaResponse := common.GetMetadataValue(ctx, CaptchaResponse)
	challenge := common.GetMetadataValue(ctx, LoginChallenge)
	email := strings.ToLower(strings.TrimSpace(req.GetEmail()))
	firstname := req.GetFirstName()
	lastname := req.GetLastName()
	password := strings.TrimSpace(req.GetPassword())
	cpassword := strings.TrimSpace(req.GetConfirmPassword())

	errors := validator.Val(
		s.validate,
		validator.Field{
			Param:   "first_name",
			Message: "Invalid first name",
			Value:   firstname,
			Tag:     "required,alpha,max=64",
		},
		validator.Field{
			Param:   "last_name",
			Message: "Invalid last name",
			Value:   lastname,
			Tag:     "required,alpha,max=64",
		},
		validator.Field{
			Param:   "email",
			Message: "Invalid email",
			Value:   email,
			Tag:     "required,email,emailMX,max=64",
		},
		validator.Field{
			Param:   "password",
			Message: "Password invalid",
			Value:   password,
			Tag:     "required,min=8,max=64,containsany=\"!\"#$%&'()*+0x2C-./:;<=>?@[]^_`{0x7C}~\"", // Use the UTF-8 hex representation for pipe "|" is 0x7C and comma "," 0x2C
		},
		validator.Field{
			Param:      "confirm_password",
			Message:    "Passwords do not match",
			Value:      cpassword,
			OtherValue: password,
			Tag:        `eqfield`,
		},
		validator.Field{
			Param:   CaptchaResponse,
			Message: CaptchaResponse + " header required",
			Value:   captchaResponse,
			Tag:     `required`,
		},
		validator.Field{
			Param:   XForwardedFor,
			Message: XForwardedFor + " header required",
			Value:   ip,
			Tag:     `required`,
		},
		validator.Field{
			Param:   LoginChallenge,
			Message: LoginChallenge + " header required",
			Value:   challenge,
			Tag:     `required`,
		},
	)
	// Validate
	if len(errors) > 0 {
		return nil, s.returnErrors(ctx, errors, codes.InvalidArgument, "Invalid arguments", api)
	}
	// Prepend IP for logging
	api = fmt.Sprintf("[%v] %v", ip, api)

	// Verify reCAPTCHA
	if s.production {
		// Get ip via headers
		rcpResp, err := recaptcha.Verify(captchaResponse, ip)
		if err != nil || !rcpResp.Success {
			s.logger.Errorf("%v: recaptcha verify: %v", api, err)
			return nil, status.Error(codes.InvalidArgument, "Captcha Verification Failed")
		}
	}

	// Check if email exists
	u, err := s.findAccountByEmail(nil, email)
	if err != nil {
		s.logger.Errorf("%v: %v", api, err)
		return nil, status.Error(codes.Internal, "An Internal error has occurred")
	}
	// Reject is account exists
	if u != nil {
		return nil, s.returnErrors(ctx, []validator.Error{
			{
				Param:   "email",
				Message: "Email is already in used",
				Value:   email,
			},
		}, codes.AlreadyExists, "Email is already in used", api)
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorf("%v: %v", api, err)
		return nil, status.Error(codes.Internal, "An Internal error has occurred")
	}

	// Build Account model
	u = &models.Account{
		Auth: models.Auth{
			Email:     email,
			Password:  string(hash),
			FirstName: firstname,
			LastName:  lastname,
			Name:      firstname + " " + lastname,
		},
		Sessions: []models.Session{
			{
				IP:        ip,
				Timestamp: time.Now(),
			},
		},
		Object: "account",
	}
	id, err := s.accountsRepo.Save(nil, u)
	if err != nil {
		s.logger.Errorf("%v: account saving: %v", api, err)
		return nil, status.Error(codes.Internal, "An Internal error has occurred")
	}

	// Authenticate with Hydra
	r, err := s.oAuthClient.AcceptLogin(challenge, &oauth.HydraLoginAccept{
		Subject:     id,
		Remember:    true,
		RememberFor: 0, // TODO: Change with env variable
	})
	if err != nil {
		s.logger.Errorf("%v: oAuthClient AcceptLogin: %v", api, err)
		return nil, status.Error(codes.Internal, "An Internal error has occurred")
	}

	return &accountsV1.RedirectResponse{RedirectTo: r.RedirectTo}, nil
}

// Authenticate is a gRPC handler that allows account to authenticate
func (s *Service) Authenticate(ctx context.Context, req *accountsV1.AuthenticateRequest) (*accountsV1.RedirectResponse, error) {
	api := "Authenticate"

	ip := common.GetMetadataValue(ctx, "x-forwarded-for")
	captchaResponse := common.GetMetadataValue(ctx, "captcha-response")

	challenge := common.GetMetadataValue(ctx, LoginChallenge)
	email := strings.ToLower(strings.TrimSpace(req.GetEmail()))
	password := strings.TrimSpace(req.GetPassword())

	errors := validator.Val(
		s.validate,
		validator.Field{
			Param:   "email",
			Message: "Wrong email or password",
			Value:   email,
			Tag:     "required,email,emailMX,max=64",
		},
		validator.Field{
			Param:          "password",
			Message:        "Wrong email or password",
			Value:          password,
			Tag:            "required",
			OmitParamValue: true,
		},
		validator.Field{
			Param:   "captcha-response",
			Message: "captcha-response header required",
			Value:   captchaResponse,
			Tag:     `required`,
		},
		validator.Field{
			Param:   "x-forward-for",
			Message: "x-forward-for header required",
			Value:   ip,
			Tag:     `required`,
		},
	)
	// Validate
	if len(errors) > 0 {
		return nil, s.returnErrors(ctx, errors, codes.PermissionDenied, "Invalid arguments", api)
	}
	// Prepend IP for logging
	api = fmt.Sprintf("[%v] %v", ip, api)

	// Verify reCAPTCHA
	if s.production {
		// Get ip via headers
		rcpResp, err := recaptcha.Verify(captchaResponse, ip)
		if err != nil || !rcpResp.Success {
			s.logger.Errorf("%v: recaptcha verify: %v", api, err)
			return nil, status.Error(codes.InvalidArgument, "Captcha Verification Failed")
		}
	}

	// Retrieve account
	u, err := s.findAccountByEmail(nil, email)
	if err != nil {
		s.logger.Errorf("%v: %v", api, err)
		return nil, status.Error(codes.Internal, "An Internal error has occurred")
	}
	if u == nil {
		s.logger.Warnf("%v: %v", api, err)
		return nil, s.returnErrors(ctx, []validator.Error{
			{
				Param:   "password",
				Message: "Wrong email or password",
			},
		}, codes.PermissionDenied, "Wrong email or password", api)
	}

	// Add time constant for comparing hash
	if err := bcrypt.CompareHashAndPassword([]byte(u.Auth.Password), []byte(password)); err != nil {
		return nil, s.returnErrors(ctx, []validator.Error{
			{
				Param:   "password",
				Message: "Wrong email or password",
			},
		}, codes.PermissionDenied, "Wrong email or password", api)
	}

	// Authenticate via Hydra
	r, err := s.oAuthClient.AcceptLogin(challenge, &oauth.HydraLoginAccept{
		Subject:     u.ID.Hex(),
		Remember:    true,
		RememberFor: 0, // TODO: Change with env variable
	})
	if err != nil {
		s.logger.Errorf("%v: oAuthClient AcceptLogin: %v", api, err)
		return nil, status.Error(codes.Internal, "An Internal error has occurred")
	}

	// Add Device Session
	u.LoggedIn = time.Now()
	_, err = s.accountsRepo.Update(
		nil,
		bson.M{"_id": u.ID},
		bson.M{
			"$set": bson.M{
				"auth.first_name": "Edited",
				"updated_at":      time.Now(),
				"logged_in":       time.Now(),
			},
		},
	)
	if err != nil {
		s.logger.Errorf("%v: %v", api, err)
		return nil, status.Error(codes.Internal, "An Internal error has occurred")
	}

	return &accountsV1.RedirectResponse{RedirectTo: r.RedirectTo}, nil
}

func (s *Service) EmailExists(ctx context.Context, req *accountsV1.EmailExistsRequest) (*accountsV1.EmailExistsResponse, error) {
	api := "EmailExists"

	email := strings.ToLower(strings.TrimSpace(req.GetEmail()))
	ip := common.GetMetadataValue(ctx, "x-forwarded-for")
	captchaResponse := common.GetMetadataValue(ctx, "captcha-response")

	errors := validator.Val(
		s.validate,
		validator.Field{
			Param:   "email",
			Message: "Invalid email",
			Value:   email,
			Tag:     "required,email,emailMX,max=64",
		},
		validator.Field{
			Param:   "captcha-response",
			Message: "captcha-response header required",
			Value:   captchaResponse,
			Tag:     `required`,
		},
		validator.Field{
			Param:   "ip",
			Message: "x-forward-for header required",
			Value:   ip,
			Tag:     `required`,
		},
	)

	// Validate
	if len(errors) > 0 {
		return nil, s.returnErrors(ctx, errors, codes.InvalidArgument, "Invalid arguments", api)
	}

	// Verify reCAPTCHA
	if s.production {
		// Get ip via headers
		rcpResp, err := recaptcha.Verify(captchaResponse, ip)
		if err != nil || !rcpResp.Success {
			s.logger.Errorf("%v: recaptcha verify: %v", api, err)
			return nil, status.Error(codes.InvalidArgument, "Captcha Verification Failed")
		}
	}

	// Check if email exists
	u, err := s.findAccountByEmail(nil, email)
	if err != nil {
		s.logger.Errorf("%v: %v", api, err)
		return nil, status.Error(codes.Internal, "An Internal error has occurred")
	}

	return &accountsV1.EmailExistsResponse{Exist: u != nil}, nil
}
