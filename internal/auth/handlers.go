package auth

import (
	"context"

	pb "github.com/isaiahwong/auth-go/protogen/auth/v1"
)

func (*Service) IsAuthenticated(ctx context.Context, in *pb.Empty) (*pb.AuthenticateResponse, error) {
	return nil, nil
}
func (*Service) SignUp(ctx context.Context, in *pb.SignUpRequest) (*pb.UserResponse, error) {
	return nil, nil
}
func (*Service) SignIn(ctx context.Context, in *pb.SignInRequest) (*pb.UserResponse, error) {
	return nil, nil
}
