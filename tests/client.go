package main

import (
	"context"
	"fmt"
	"time"

	pb "github.com/isaiahwong/accounts-go/api/accounts/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var client pb.AccountsServiceClient
var conn *grpc.ClientConn

func init() {
	var err error
	conn, err = grpc.Dial("api-accounts-v1-accountsservice.default.svc.cluster.local:5000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client = pb.NewAccountsServiceClient(conn)
}

func SignUp() {
	ctx := context.Background()
	var trailer metadata.MD // variable to store header and trailer

	resp, err := client.SignUp(ctx, &pb.SignUpRequest{
		Email:           "sada@mpillow.com",
		Password:        "12345678UF020|",
		ConfirmPassword: "12345678UF020|",
	}, grpc.Trailer(&trailer)) // will retrieve trailer)
	fmt.Println(err)
	fmt.Println(trailer.Get("errors-bin"))
	fmt.Println(resp)
}

func Login() {
	ctx := context.Background()
	var trailer metadata.MD // variable to store header and trailer
	ctx = metadata.AppendToOutgoingContext(ctx, "x-forwarded-for", "localhost", "captcha-response", "1234")

	resp, err := client.Authenticate(ctx, &pb.AuthenticateRequest{
		Email:    "isaiah.wong@jirehsoho.com",
		Password: "password!",
	}, grpc.Trailer(&trailer)) // will retrieve trailer)
	fmt.Println(err)
	fmt.Println(trailer.Get("errors-bin"))
	fmt.Println(resp)
}

func Loop() {
	c1 := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second)
		c1 <- conn.GetState().String()
	}()

	select {
	case res := <-c1:
		fmt.Println(res)
	}
}

func main() {
	defer conn.Close()
	Login()
}
