package user

import (
	"context"
	"time"

	"github.com/isaiahwong/auth-go/internal/models"
	"github.com/isaiahwong/auth-go/internal/store/types/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mongoUserRepo struct {
	m *mongo.MongoStore
}

func (r *mongoUserRepo) Save(ctx context.Context, u *models.User, id string) (string, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), r.m.Timeout)
		defer cancel()
	}
	coll := r.m.Client.Database(r.m.Database).Collection("user")
	u.ID = primitive.NewObjectID()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	res, err := coll.InsertOne(ctx, u)
	if err != nil {
		return "", err
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "nil", &mongo.OIDTypeError{}
	}

	return oid.String(), nil
}

func (r *mongoUserRepo) Find(ctx context.Context, s interface{}, opts ...interface{}) []*models.User {
	return nil
}

func (r *mongoUserRepo) FindOne(ctx context.Context, s interface{}, opts ...interface{}) *models.User {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), r.m.Timeout)
		defer cancel()
	}
	coll := r.m.Client.Database(r.m.Database).Collection("user")
	resp := coll.FindOne(ctx, s)
	if resp == nil {
		return nil
	}
	user := &models.User{}
	err := resp.Decode(user)
	if err != nil {
		// TODO assert error type
		return nil
	}
	return user
}

// NewMongoUserRepo returns a new Mongo Based Repo
func NewMongoUserRepo(m *mongo.MongoStore) Repo {
	return &mongoUserRepo{m}
}
