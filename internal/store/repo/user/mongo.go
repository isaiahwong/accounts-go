package user

import (
	"context"
	"errors"
	"time"

	"github.com/isaiahwong/accounts-go/internal/models"
	mt "github.com/isaiahwong/accounts-go/internal/store/drivers/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo "go.mongodb.org/mongo-driver/mongo"
)

// ErrOIDType defines and invalid mongo object id
var ErrOIDType = errors.New("Invalid OID")

type mongoUserRepo struct {
	m *mt.MongoStore
}

func (r *mongoUserRepo) Save(ctx context.Context, u *models.User) (string, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), r.m.Timeout)
		defer cancel()
	}
	coll := r.m.Client.Database(r.m.Database).Collection("user")
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	res, err := coll.InsertOne(ctx, u)
	if err != nil {
		return "", err
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", ErrOIDType
	}

	return oid.String(), nil
}

func (r *mongoUserRepo) Find(ctx context.Context, s interface{}, opts ...interface{}) ([]*models.User, error) {
	return nil, nil
}

func (r *mongoUserRepo) FindOne(ctx context.Context, s interface{}, opts ...interface{}) (*models.User, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), r.m.Timeout)
		defer cancel()
	}
	coll := r.m.Client.Database(r.m.Database).Collection("user")
	resp := coll.FindOne(ctx, s)
	user := &models.User{}
	err := resp.Decode(user)

	switch err {
	case mongo.ErrNoDocuments:
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

// NewMongoUserRepo returns a new Mongo Based Repo
func NewMongoUserRepo(m *mt.MongoStore) Repo {
	return &mongoUserRepo{m}
}
