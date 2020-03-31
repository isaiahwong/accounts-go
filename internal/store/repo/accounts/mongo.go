package accounts

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
var ErrUpdateDocuments = errors.New("Documents not updated")

type mongoAccountsRepo struct {
	m    *mt.MongoStore
	name string
}

func (r *mongoAccountsRepo) GetTimeout() time.Duration {
	return r.m.Timeout
}

func (r *mongoAccountsRepo) Save(ctx context.Context, u *models.Account) (string, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), r.m.Timeout)
		defer cancel()
	}
	coll := r.m.Client.Database(r.m.Database).Collection(r.name)
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
	return oid.Hex(), nil
}

func (r *mongoAccountsRepo) Update(ctx context.Context, f interface{}, up interface{}) (int, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), r.m.Timeout)
		defer cancel()
	}
	coll := r.m.Client.Database(r.m.Database).Collection(r.name)
	resp, err := coll.UpdateOne(ctx, f, up)
	if err != nil {
		return 0, err
	}

	if resp.MatchedCount < 1 {
		return int(resp.MatchedCount), mongo.ErrNoDocuments
	}

	if resp.ModifiedCount < 1 {
		return int(resp.ModifiedCount), ErrUpdateDocuments
	}

	return int(resp.ModifiedCount), nil
}

func (r *mongoAccountsRepo) Find(ctx context.Context, s interface{}, opts ...interface{}) ([]*models.Account, error) {
	return nil, nil
}

func (r *mongoAccountsRepo) FindOne(ctx context.Context, s interface{}, opts ...interface{}) (*models.Account, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), r.m.Timeout)
		defer cancel()
	}
	coll := r.m.Client.Database(r.m.Database).Collection(r.name)
	resp := coll.FindOne(ctx, s)
	account := &models.Account{}
	err := resp.Decode(account)

	switch err {
	case mongo.ErrNoDocuments:
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return account, nil
}

// NewMongoAccountsRepo returns a new Mongo Based Repo
func NewMongoAccountsRepo(m *mt.MongoStore) Repo {
	return &mongoAccountsRepo{m, "accounts"}
}
