package store

import (
	"context"
	"fmt"
	"time"

	"github.com/isaiahwong/go-services/src/payment/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoStore struct wrapper
type MongoStore struct {
	client *mongo.Client
	opts   *mongoOptions
}

// NewMongoStore provides a new MongoStore
func NewMongoStore(opt ...MongoOption) (*MongoStore, error) {
	opts := defaultOptions
	// Get the initial mongo settings
	mongoOpt := options.Client()

	// Apply options
	for _, o := range opt {
		o(&opts, mongoOpt)
	}

	// Create a new mongo client
	c, err := mongo.NewClient(
		mongoOpt,
	)
	if err != nil {
		return nil, err
	}
	return &MongoStore{
		client: c,
		opts:   &opts,
	}, nil
}

// Connect connects to mongodb
func (m *MongoStore) Connect(ctx context.Context) error {
	err := m.client.Connect(ctx)
	if err != nil {
		return &connectError{fmt.Sprintf("Mongo Connection %v", err)}
	}

	pe := m.Ping()
	if pe != nil {
		return pe
	}

	return nil
}

// Disconnect Disconnects Mongo client
func (m *MongoStore) Disconnect(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

// Ping verifies that the client can connect to the topology.
func (m *MongoStore) Ping() error {
	ctx, _ := context.WithTimeout(context.Background(), m.opts.timeout)
	// Test connection
	err := m.client.Ping(ctx, readpref.Primary())
	if err != nil {
		return &connectError{fmt.Sprintf("Mongo Connection %v", err)}
	}
	return nil
}

// Create creates a new payment object
func (m *MongoStore) Create(ctx context.Context, payment *model.Payment) (*primitive.ObjectID, error) {
	coll := m.client.Database(m.opts.database).Collection("payment")
	payment.ID = primitive.NewObjectID()
	payment.Updated = time.Now()
	payment.Created = time.Now()

	res, err := coll.InsertOne(ctx, payment)
	if err != nil {
		return nil, err
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, &idError{"Invalid OID"}
	}
	return &oid, nil
}
