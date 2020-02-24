package auth

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/isaiahwong/auth-go/internal/models"
	"github.com/isaiahwong/auth-go/internal/util/recaptcha"
	"github.com/isaiahwong/auth-go/internal/util/validator"
	pb "github.com/isaiahwong/auth-go/protogen/auth/v1"
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
	e := strings.ToLower(strings.TrimSpace(req.GetEmail()))
	p := strings.TrimSpace(req.GetPassword())
	cp := strings.TrimSpace(req.GetConfirmPassword())
	tok := strings.TrimSpace(req.GetCaptchaToken())
	ip := strings.TrimSpace(req.GetIp())
	md := metadata.Pairs()

	errors := validator.Val(
		s.validate,
		validator.Field{
			Param:   "email",
			Message: "Invalid email",
			Value:   e,
			Tag:     "required,email,max=64",
		},
		validator.Field{
			Param:   "password",
			Message: "Password required",
			Value:   p,
			Tag:     "required,max=64",
		},
		validator.Field{
			Param:      "confirm_password",
			Message:    "Passwords do not match",
			Value:      cp,
			OtherValue: p,
			Tag:        `eqfield`,
		},
		validator.Field{
			Param:   "captcha_token",
			Message: "Captcha token required",
			Value:   tok,
			Tag:     `required`,
		},
	)

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

	// rcpResp := recaptcha.Response{}
	if s.production {
		rcpResp, err := recaptcha.Verify(tok, ip)
		if err != nil || !rcpResp.Success {
			s.logger.Errorf("SignUp: recaptcha verify: %v", err)
			return nil, status.Error(codes.InvalidArgument, "Captcha Verification Failed")
		}
		// TODO: store challenge ts
		// 	t, err := time.Parse(time.RFC3339, r.ChallengeTSISO)
	}

	// Hash password

	u := &models.User{
		Auth: models.Auth{
			Email: e,
		},
	}
	s.userRepo.Save(nil, u)
	// oid, err := u.Save(nil, s.store)
	// if err != nil {
	// 	s.logger.Errorf("SignUp: Error saving user: %v", err)
	// 	return nil, status.Error(codes.Internal, "An Internal error has occurred")
	// }
	// // Check if email exists
	// return &pb.UserResponse{User: &pb.User{
	// 	Id:     oid.String(),
	// 	Object: u.Object,
	// }}, nil
	return nil, nil
}

// SignIn is a gRPC handler that allows user to authenticate
func (*Service) SignIn(ctx context.Context, in *pb.SignInRequest) (*pb.UserResponse, error) {
	return nil, nil
}
