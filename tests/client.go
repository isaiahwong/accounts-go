package main

import (
	"context"
	"fmt"

	pb "github.com/isaiahwong/accounts-go/api/auth/v1"
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
	var trailer metadata.MD // variable to store header and trailer

	resp, err := client.SignUp(ctx, &pb.SignUpRequest{
		Email:           "sada@mpillow.com",
		Password:        "12345678UF020|",
		ConfirmPassword: "12345678UF020|",
		CaptchaToken:    "123",
		Ip:              "1234",
	}, grpc.Trailer(&trailer)) // will retrieve trailer)
	fmt.Println(err)
	fmt.Println(trailer.Get("errors-bin"))
	fmt.Println(resp)
}
