package store

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongoOptions a set of mongo options declared privately
type mongoOptions struct {
	connstr  string
	database string
	auth     *options.Credential
	timeout  time.Duration
}

// MongoOption sets options such as hostPort; parameters, etc.
type MongoOption func(*mongoOptions, *options.ClientOptions)

var defaultOptions = mongoOptions{
	connstr:  "mongodb://localhost:27017/",
	database: "auth",
	timeout:  15,
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

// WithTimeout specifies the timeout for an initial connection to a server.
// If a custom Dialer is used, this method won't be set and the user is
// responsible for setting the ConnectTimeout for connections on the dialer
// themselves.
func WithTimeout(t time.Duration) MongoOption {
	return func(o *mongoOptions, m *options.ClientOptions) {
		o.timeout = t
		m.SetConnectTimeout(1 * time.Second)
	}
}

// WithAuth Authentication for mongodb
func WithAuth(credential MongoCredential) MongoOption {
	return func(o *mongoOptions, m *options.ClientOptions) {
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
