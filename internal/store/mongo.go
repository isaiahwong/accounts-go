package store

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoStore struct wrapper
type MongoStore struct {
	Client   *mongo.Client
	opts     *mongoOptions
	Database string
	Timeout  time.Duration
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
		Client:   c,
		opts:     &opts,
		Database: opts.database,
		Timeout:  opts.timeout,
	}, nil
}

// Connect connects to mongodb
func (m *MongoStore) Connect(ctx context.Context) error {
	var cancel context.CancelFunc
	if ctx == nil {
		ctx, cancel = context.WithTimeout(context.Background(), m.opts.timeout)
		defer cancel()
	}

	err := m.Client.Connect(ctx)
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
	return m.Client.Disconnect(ctx)
}

// Ping verifies that the client can connect to the topology.
func (m *MongoStore) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.opts.timeout)
	defer cancel()
	// Test connection
	err := m.Client.Ping(ctx, readpref.Primary())
	if err != nil {
		return &connectError{fmt.Sprintf("Mongo Connection %v", err)}
	}
	return nil
}
