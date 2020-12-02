package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type ErrConnect struct {
	S string
}

func (e *ErrConnect) Error() string {
	return "Mongo Connection: " + e.S
}

// MongoStore struct wrapper
type MongoStore struct {
	Client   *mongo.Client
	opts     *mongoOptions
	Database string
	Timeout  time.Duration
}

// mongoOptions a set of mongo options declared privately
type mongoOptions struct {
	connstr        string
	database       string
	auth           *options.Credential
	timeout        time.Duration
	initialTimeout time.Duration
}

// MongoOption sets options such as hostPort; parameters, etc.
type MongoOption func(*mongoOptions, *options.ClientOptions)

var defaultOptions = mongoOptions{
	connstr:        "mongodb://localhost:27017",
	database:       "auth",
	timeout:        10,
	initialTimeout: 10,
}

// MongoCredential holds auth options.
// taken from go.mongodb.org/mongo-driver/mongo/options
type MongoCredential struct {
	AuthMechanism           string
	AuthMechanismProperties map[string]string
	AuthSource              string
	Username                string
	Password                string
	PasswordSet             bool
}

// WithConnectionString returns MongoOption; sets
func WithConnectionString(connstr string) MongoOption {
	return func(o *mongoOptions, m *options.ClientOptions) {
		o.connstr = connstr
		m.ApplyURI(connstr)
	}
}

// WithDatabase returns MongoOption; sets default database
func WithDatabase(db string) MongoOption {
	return func(o *mongoOptions, m *options.ClientOptions) {
		o.database = db
	}
}

// WithTimeout specifies the timeout for requests to the server.
func WithTimeout(t time.Duration) MongoOption {
	return func(o *mongoOptions, m *options.ClientOptions) {
		o.timeout = t
	}
}

// WithInitialTimeout specifies the timeout for an initial connection to a server.
// If a custom Dialer is used, this method won't be set and the user is
// responsible for setting the ConnectTimeout for connections on the dialer
// themselves.
func WithInitialTimeout(t time.Duration) MongoOption {
	return func(o *mongoOptions, m *options.ClientOptions) {
		o.initialTimeout = t
		m.SetConnectTimeout(t)
	}
}

// WithAuth Authentication for mongodb
func WithAuth(credential MongoCredential) MongoOption {
	return func(o *mongoOptions, m *options.ClientOptions) {
		if credential.Username == "" || credential.Password == "" {
			return
		}
		m.SetAuth(options.Credential{
			AuthMechanism: credential.AuthMechanism,
			AuthSource:    credential.AuthSource,
			Username:      credential.Username,
			Password:      credential.Password,
			PasswordSet:   credential.PasswordSet,
		})
	}
}

// WithHeartbeat TODO
func WithHeartbeat(t time.Duration) MongoOption {
	return func(o *mongoOptions, m *options.ClientOptions) {
		m.SetHeartbeatInterval(t)
	}
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
		ctx, cancel = context.WithTimeout(context.Background(), m.opts.initialTimeout)
		defer cancel()
	}
	err := m.Client.Connect(ctx)
	if err != nil {
		return &ErrConnect{S: fmt.Sprint(err)}
	}

	// Test Connectivity
	pe := m.Ping(&m.opts.initialTimeout)
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
func (m *MongoStore) Ping(t *time.Duration) error {
	if t == nil {
		t = &m.opts.timeout
	}
	ctx, cancel := context.WithTimeout(context.Background(), *t)
	defer cancel()
	// Test connection
	err := m.Client.Ping(ctx, readpref.Primary())
	if err != nil {
		return &ErrConnect{S: fmt.Sprint(err)}
	}
	return nil
}
