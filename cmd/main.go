package cmd

import (
	accounts "github.com/isaiahwong/accounts-go/internal/accounts"
	"github.com/isaiahwong/accounts-go/internal/server"
	"github.com/isaiahwong/accounts-go/internal/store/drivers/mongo"
	"github.com/isaiahwong/accounts-go/internal/util/log"
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
		mongo.WithInitialTimeout(config.DBInitialTimeout),
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
		server.WithAddress(config.Address),
		server.WithLogger(l),
		server.WithName("Accounts Service"),
		server.WithDataStore(m),
	)
	s.Production = config.Production

	if err != nil {
		l.Fatalf("NewServer: %v", err)
	}

	// Register authentication service
	accounts.RegisterService(
		accounts.WithLogger(l),
		accounts.WithGrpc(s.GRPCServer),
		accounts.WithStore(m),
	)
}

// Execute starts application
func Execute() {
	err := s.Serve()
	panic(err)
}
