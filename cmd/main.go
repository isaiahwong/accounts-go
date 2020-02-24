package cmd

import (
	auth "github.com/isaiahwong/auth-go/internal/auth"
	"github.com/isaiahwong/auth-go/internal/server"
	"github.com/isaiahwong/auth-go/internal/store/types/mongo"
	"github.com/isaiahwong/auth-go/internal/util/log"
)

var s *server.Server

func init() {
	var err error

	config := loadEnv()

	// Initialize a new logger
	l := log.NewLogrusLogger()

	// Initialize a new Store
	m, err := mongo.NewMongoStore(
		mongo.WithDatabase(config.DBName),
		mongo.WithConnectionString(config.DBUri),
		mongo.WithTimeout(config.DBTimeout),
		mongo.WithAuth(mongo.MongoCredential{
			Username: config.DBUser,
			Password: config.DBPassword,
		}),
		mongo.WithHeartbeat(config.DBTimeout),
	)
	if err != nil {
		l.Fatalf("NewMongoStore: %v", err)
	}

	// Initialize a new Server
	s, err = server.New(
		server.WithAddress(":50051"),
		server.WithLogger(l),
		server.WithName("Auth Service"),
		server.WithDataStore(m),
	)
	s.Production = config.Production

	if err != nil {
		l.Fatalf("NewServer: %v", err)
	}

	// Register authentication service
	auth.RegisterService(
		auth.WithLogger(l),
		auth.WithGrpc(s.GRPCServer),
		auth.WithStore(m),
	)
}

// Execute starts application
func Execute() {
	err := s.Serve()
	panic(err)
}
