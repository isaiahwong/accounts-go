package accounts

import (
	"context"
	"encoding/json"

	"github.com/isaiahwong/accounts-go/internal/models"
	"github.com/isaiahwong/accounts-go/internal/util/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s *Service) returnErrors(ctx context.Context, errors []validator.Error, code codes.Code, msg string, prefix string) error {
	if len(errors) < 0 {
		return nil
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

func (s *Service) findUserByEmail(ctx context.Context, email string) (*models.User, error) {
	u, err := s.userRepo.FindOne(ctx, bson.M{
		"$or": []interface{}{
			bson.M{"auth.email": email},
		},
	})
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) findUserByID(ctx context.Context, id string) (*models.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	u, err := s.userRepo.FindOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return nil, err
	}
	return u, nil
}
