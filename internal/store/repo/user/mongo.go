package user

import (
	"context"

	"github.com/isaiahwong/auth-go/internal/models"
	"github.com/isaiahwong/auth-go/internal/store/types/mongo"
)

type mongoUser struct {
	m *mongo.MongoStore
}

func (*mongoUser) Save(ctx context.Context, u *models.User) string {
	return ""
}
func (*mongoUser) Find(ctx context.Context, s string) []*models.User {
	return nil
}
func (*mongoUser) FindOne(ctx context.Context) *models.User {
	return nil
}

// NewMongoUserRepo returns a new Mongo Based Repo
func NewMongoUserRepo(m *mongo.MongoStore) Repo {
	return &mongoUser{m}
}
