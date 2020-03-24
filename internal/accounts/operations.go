package accounts

import (
	"context"
	"encoding/json"

	"github.com/isaiahwong/accounts-go/internal/common/validator"
	"github.com/isaiahwong/accounts-go/internal/models"
	"github.com/isaiahwong/accounts-go/internal/oauth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s *Service) returnErrors(ctx context.Context, errors []validator.Error, code codes.Code, msg string, prefix string) error {
	if len(errors) < 0 {
		s.logger.Errorf("%v: %v", prefix, "returnErrors: errors is empty")
		return status.Error(codes.Internal, "An Internal error has occurred")
	}
	md := metadata.Pairs()
	json, jerr := json.Marshal(errors)
	if jerr != nil {
		s.logger.Errorf("%v: %v", prefix, jerr)
		return status.Error(codes.Internal, "An Internal error has occurred")
	}
	s.logger.Warnf("%v: %v", prefix, string(json))
	md.Append("errors-bin", string(json))
	grpc.SetTrailer(ctx, md)
	return status.Error(code, msg)
}

func (s *Service) returnHydraError(ctx context.Context, he *oauth.HydraError, prefix string) error {
	if he == nil {
		s.logger.Errorf("%v: %v", prefix, "returnHydraError: HydraError is nil")
		return status.Error(codes.Internal, "An Internal error has occurred")
	}
	md := metadata.Pairs()
	json, jerr := json.Marshal(he)
	if jerr != nil {
		s.logger.Errorf("%v: %v", prefix, jerr)
		return status.Error(codes.Internal, "An Internal error has occurred")
	}
	s.logger.Warnf("%v: %v", prefix, string(json))
	md.Append("errors-bin", string(json))
	grpc.SetTrailer(ctx, md)

	var c codes.Code
	desc := he.ErrorDescription

	switch he.StatusCode {
	case 401:
		c = codes.PermissionDenied
	case 404:
		c = codes.NotFound
	case 409:
		c = codes.AlreadyExists
	case 500:
		c = codes.Internal
	default:
		c = codes.Internal
	}
	if c == codes.Internal {
		desc = "An Internal error has occurred"
	}
	return status.Error(c, desc)
}

func (s *Service) findAccountByEmail(ctx context.Context, email string) (*models.Account, error) {
	u, err := s.accountsRepo.FindOne(ctx, bson.M{
		"$or": []interface{}{
			bson.M{"auth.email": email},
		},
	})
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) findAccountByID(ctx context.Context, id string) (*models.Account, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	u, err := s.accountsRepo.FindOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return nil, err
	}
	return u, nil
}
