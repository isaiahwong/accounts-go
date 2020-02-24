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
			Tag:     "required,email,max=64",
		},
		validator.Field{
			Param:   "password",
			Message: "Password required",
			Value:   password,
			Tag:     "required,max=64",
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
	if errors != nil {
		json, jerr := json.Marshal(errors)
		if jerr != nil {
			s.logger.Errorf("SignUp Handler: %v", jerr)
			return nil, status.Error(codes.Internal, "Unexpected error")
		}
		md.Append("errors-bin", string(json))
		grpc.SetTrailer(ctx, md)
		return nil, status.Error(codes.InvalidArgument, "Invalid arguments")
	}

	// Verify reCAPTCHA
	// rcpResp := recaptcha.Response{}
	if s.production {
		rcpResp, err := recaptcha.Verify(token, ip)
		if err != nil || !rcpResp.Success {
			s.logger.Errorf("SignUp: recaptcha verify: %v", err)
			return nil, status.Error(codes.InvalidArgument, "Captcha Verification Failed")
		}
		// TODO: store challenge ts
		// 	t, err := time.Parse(time.RFC3339, r.ChallengeTSISO)
	}

	// Check if email exists
	u := s.userRepo.FindOne(nil, bson.M{
		"$or": []interface{}{
			bson.M{"auth.email": email},
		},
	})

	// Reject is user exists
	if u != nil {
		json, jerr := json.Marshal(validator.Error{
			Param:   "email",
			Message: "Email is already in used",
			Value:   email,
		})
		if jerr != nil {
			s.logger.Errorf("SignUp Handler: %v", jerr)
			return nil, status.Error(codes.Internal, "An Internal error has occurred")
		}
		md.Append("errors-bin", string(json))
		grpc.SetTrailer(ctx, md)
		return nil, status.Error(codes.InvalidArgument, "Invalid arguments")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorf("bcrypt.GenerateFromPassword: %v", err)
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
		s.logger.Errorf("userRepo.Save: %v", err)
		return nil, status.Error(codes.Internal, "An Internal error has occurred")
	}
	return &pb.UserResponse{User: &pb.User{
		Id:     id,
		Object: u.Object,
	}}, nil
}

// SignIn is a gRPC handler that allows user to authenticate
func (*Service) SignIn(ctx context.Context, in *pb.SignInRequest) (*pb.UserResponse, error) {
	return nil, nil
}
