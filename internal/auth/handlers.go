package auth

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/isaiahwong/auth-go/internal/models"
	"github.com/isaiahwong/auth-go/internal/util/recaptcha"
	"github.com/isaiahwong/auth-go/internal/util/validator"
	pb "github.com/isaiahwong/auth-go/protogen/auth/v1"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s *Service) IsAuthenticated(ctx context.Context, in *pb.Empty) (*pb.AuthenticateResponse, error) {
	return nil, nil
}

// SignUp is a gRPC handler allows user to register
func (s *Service) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.UserResponse, error) {
	api := "/v1/auth/signup"

	email := strings.ToLower(strings.TrimSpace(req.GetEmail()))
	password := strings.TrimSpace(req.GetPassword())
	cpassword := strings.TrimSpace(req.GetConfirmPassword())
	token := strings.TrimSpace(req.GetCaptchaToken())
	ip := strings.TrimSpace(req.GetIp())
	md := metadata.Pairs()

	errors := validator.Val(
		s.validate,
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
			Tag:     "required,min=8,max=32,containsany=\"!\"#$%&'()*+0x2C-./:;<=>?@[]^_`{0x7C}~\"", // Use the UTF-8 hex representation for pipe "|" is 0x7C and comma "," 0x2C
		},
		validator.Field{
			Param:      "confirm_password",
			Message:    "Passwords do not match",
			Value:      cpassword,
			OtherValue: password,
			Tag:        `eqfield`,
		},
		validator.Field{
			Param:   "captcha_token",
			Message: "Captcha token required",
			Value:   token,
			Tag:     `required`,
		},
	)

	// Validate
	if len(errors) > 0 {
		json, jerr := json.Marshal(errors)
		if jerr != nil {
			s.logger.Errorf("%b: %v", api, jerr)
			return nil, status.Error(codes.Internal, "Unexpected error")
		}
		s.logger.Warnf("%v: %v", api, string(json))
		md.Append("errors-bin", string(json))
		grpc.SetTrailer(ctx, md)
		return nil, status.Error(codes.InvalidArgument, "Invalid arguments")
	}

	// Verify reCAPTCHA
	if s.production {
		rcpResp, err := recaptcha.Verify(token, ip)
		if err != nil || !rcpResp.Success {
			s.logger.Errorf("%v: recaptcha verify: %v", api, err)
			return nil, status.Error(codes.InvalidArgument, "Captcha Verification Failed")
		}
	}

	// Check if email exists
	u, err := s.userRepo.FindOne(nil, bson.M{
		"$or": []interface{}{
			bson.M{"auth.email": email},
		},
	})
	if err != nil {
		s.logger.Errorf("%v: %v", api, err)
		return nil, status.Error(codes.Internal, "An Internal error has occurred")
	}
	// Reject is user exists
	if u != nil {
		json, jerr := json.Marshal(validator.Error{
			Param:   "email",
			Message: "Email is already in used",
			Value:   email,
		})
		if jerr != nil {
			s.logger.Errorf("%v: %v", api, jerr)
			return nil, status.Error(codes.Internal, "An Internal error has occurred")
		}
		s.logger.Warnf("%v: %v", api, string(json))
		md.Append("errors-bin", string(json))
		grpc.SetTrailer(ctx, md)
		return nil, status.Error(codes.InvalidArgument, "Invalid arguments")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorf("%v: %v", api, err)
		return nil, status.Error(codes.Internal, "An Internal error has occurred")
	}

	u = &models.User{
		Auth: models.Auth{
			Email:    email,
			Password: string(hash),
		},
		Object: "user",
	}
	id, err := s.userRepo.Save(nil, u, "")
	if err != nil {
		s.logger.Errorf("%v: user saving: %v", api, err)
		return nil, status.Error(codes.Internal, "An Internal error has occurred")
	}
	return &pb.UserResponse{User: &pb.User{
		Id:     id,
		Object: u.Object,
		Auth: &pb.Auth{
			Email: u.Auth.Email,
		},
	}}, nil
}

// SignIn is a gRPC handler that allows user to authenticate
func (*Service) SignIn(ctx context.Context, in *pb.SignInRequest) (*pb.UserResponse, error) {
	return nil, nil
}
