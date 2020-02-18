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

// ConnectionString returns MongoOption; sets
func ConnectionString(connstr string) MongoOption {
	return func(o *mongoOptions, m *options.ClientOptions) {
		o.connstr = connstr
		m.ApplyURI(connstr)
	}
}

// Database returns MongoOption; sets default database
func Database(db string) MongoOption {
	return func(o *mongoOptions, m *options.ClientOptions) {
		o.database = db
	}
}

// SetTimeout specifies the timeout for an initial connection to a server.
// If a custom Dialer is used, this method won't be set and the user is
// responsible for setting the ConnectTimeout for connections on the dialer
// themselves.
func SetTimeout(t time.Duration) MongoOption {
	return func(o *mongoOptions, m *options.ClientOptions) {
		o.timeout = t
		m.SetConnectTimeout(1 * time.Second)
	}
}

// SetAuth Authentication for mongodb
func SetAuth(credential MongoCredential) MongoOption {
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

// SetHeartbeat TODO
func SetHeartbeat(t time.Duration) MongoOption {
	return func(o *mongoOptions, m *options.ClientOptions) {
		m.SetHeartbeatInterval(t)
	}
}
