package main

import (
	"context"
	"fmt"

	pb "github.com/isaiahwong/auth-go/protogen/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewAuthServiceClient(conn)
	ctx := context.Background()
	var header, trailer metadata.MD // variable to store header and trailer

	client.SignUp(ctx, &pb.SignUpRequest{
		Email:           "isaiah.wong@jirehsoho.com",
		Password:        "12345678",
		ConfirmPassword: "12345678",
		CaptchaToken:    "123",
		Ip:              "1234",
	},
		grpc.Header(&header),   // will retrieve header
		grpc.Trailer(&trailer)) // will retrieve trailer)
	fmt.Println(trailer.Get("errors-bin"))
}
