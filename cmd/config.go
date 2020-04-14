package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/isaiahwong/accounts-go/internal/common"
	"github.com/joho/godotenv"
)

// EnvConfig Application wide env configurations
//
// AppEnv specifies if the app is in `development` or `production`
// Host specifies host address or dns
// Port specifies the port the server will run on
// EnableStackdriver specifies if google stackdriver will be enabled
// StripeSecret specifies Stripe api production key
// StripeSecretDev specifies Stripe api key for development
// StripeEndpointSecret specifies Stripe api key for webhook verification
// PaypalClientIDDev specifies Paypal api key for development
// PaypalSecretDev specifies Paypal api key secret for development
// PaypalClientID
// PaypalSecret         string
// PaypalURL specifies Paypal api URL for request
// DBUri
// DBUriDev
// DBUriTest
// DBName
// DBUser
// DBPassword
type EnvConfig struct {
	AppEnv           string
	Production       bool
	Host             string
	Address          string
	DBUri            string
	DBName           string
	DBUser           string
	DBPassword       string
	DBTimeout        time.Duration
	DBInitialTimeout time.Duration
}

// LoadEnv loads environment variables for Application
func loadEnv() *EnvConfig {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf(".env not loaded: %v\n", err)
	}

	// Convert to int
	sec, err := strconv.ParseInt(common.MapEnvWithDefaults("DB_TIMEOUT", "10"), 10, 64)
	if err != nil {
		fmt.Printf("Error parsing DB_TIMEOUT: %v\nWill fallback to default value", err)
		sec = 10
	}
	dBTimeout := time.Duration(sec) * time.Second

	sec, err = strconv.ParseInt(common.MapEnvWithDefaults("DB_INITIAL_TIMEOUT", "10"), 10, 64)
	if err != nil {
		fmt.Printf("Error parsing DB_TIMEOUT: %v\nWill fallback to default value", err)
		sec = 10
	}
	initialTimeout := time.Duration(sec) * time.Second

	return &EnvConfig{
		AppEnv:           common.MapEnvWithDefaults("APP_ENV", "development"),
		Production:       common.MapEnvWithDefaults("APP_ENV", "development") == "production",
		Address:          common.MapEnvWithDefaults("ADDRESS", ":50051"),
		DBUri:            common.MapEnvWithDefaults("DB_URI", "mongodb://localhost:27017"),
		DBName:           common.MapEnvWithDefaults("DB_NAME", "accounts"),
		DBUser:           common.MapEnvWithDefaults("MONGO_USERNAME", ""),
		DBPassword:       common.MapEnvWithDefaults("MONGO_PASSWORD", ""),
		DBTimeout:        dBTimeout,
		DBInitialTimeout: initialTimeout,
	}
}
