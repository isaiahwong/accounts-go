package accounts

import (
	"context"
	"testing"

	pb "github.com/isaiahwong/accounts-go/api/accounts/v1"
	"github.com/isaiahwong/accounts-go/internal/util/log"
	"github.com/isaiahwong/accounts-go/tests/mocks"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var logger *logrus.Logger = log.NewLogrusLogger()

var invalidEmails = []string{
	"",
	"plainaddress",
	"#@%^%#$@#$@#.com",
	"@example.com",
	"Joe Smith <email@example.com>",
	"email.example.com",
	"email@example@example.com",
	".email@example.com",
	"email.@example.com",
	"email..email@example.com",
	"あいうえお@example.com",
	"email@example.com (Joe Smith)",
	"email@example",
	"email@-example.com",
	"email@example.web",
	"email@111.222.333.44444",
	"email@example..com",
	"Abc..123@example.com",
}

var invalidPasswords = []string{
	"",
	"aA4567%",                           // 7 length too short
	"aaaaaaaA",                          // require symbols
	"12345678",                          // require symbols
	"121314151617119****************&&", // 33 length too long
}

var validPasswords = []string{
	"hello1234^",
	"hello12345^",
}

func SetupRepo() *mocks.Repo {
	r := new(mocks.Repo)
	r.On("FindOne", nil, mock.Anything).Return(nil, nil)
	return r
}

func TestSignUp(t *testing.T) {
	validReq := &pb.SignUpRequest{
		Email:           "isaiah@example.com",
		Password:        "12345678UF020|",
		ConfirmPassword: "12345678UF020|",
	}
	repo := SetupRepo()

	svc := &Service{
		production: true,
		logger:     logger,
		policy:     bluemonday.StrictPolicy(),
		userRepo:   repo,
	}
	svc.initValidator()

	req := &pb.SignUpRequest{}
	t.Run("Invalid Email", func(t *testing.T) {
		*req = *validReq
		for _, e := range invalidEmails {
			req.Email = e
			_, err := svc.SignUp(context.Background(), req)
			st, ok := status.FromError(err)
			if !ok {
				t.Errorf("Error parsing grpc error code")
			}

			assert.Equal(t, codes.InvalidArgument, st.Code())
		}
	})

	t.Run("Invalid Password", func(t *testing.T) {
		*req = *validReq
		for _, p := range invalidPasswords {
			req.Password = p
			req.ConfirmPassword = p
			_, err := svc.SignUp(context.Background(), req)
			st, ok := status.FromError(err)
			if !ok {
				t.Errorf("Error parsing grpc error code")
			}

			assert.Equal(t, codes.InvalidArgument, st.Code())
		}
	})

	t.Run("Invalid Password do not match", func(t *testing.T) {
		*req = *validReq
		req.Password = validPasswords[0]
		req.ConfirmPassword = validPasswords[1]
		_, err := svc.SignUp(context.Background(), req)
		st, ok := status.FromError(err)
		if !ok {
			t.Errorf("Error parsing grpc error code")
		}

		assert.Equal(t, codes.InvalidArgument, st.Code())
	})
}
