package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/isaiahwong/auth-go/internal/util"
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
	AppEnv     string
	Production bool
	Host       string
	Address    string
	DBUri      string
	DBName     string
	DBUser     string
	DBPassword string
	DBTimeout  time.Duration
}

// LoadEnv loads environment variables for Application
func loadEnv() *EnvConfig {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf(".env not loaded: %v", err)
	}

	// Convert to int
	sec, err := strconv.ParseInt(util.MapEnvWithDefaults("DB_TIMEOUT", "10"), 10, 64)
	if err != nil {
		fmt.Printf("Error parsing DB_TIMEOUT: %v\nWill fallback to default value", err)
		sec = 10
	}

	dBTimeout := time.Duration(sec) * time.Second

	return &EnvConfig{
		AppEnv:     util.MapEnvWithDefaults("APP_ENV", "development"),
		Production: util.MapEnvWithDefaults("APP_ENV", "development") == "true",
		Address:    util.MapEnvWithDefaults("ADDRESS", "5000"),
		DBUri:      util.MapEnvWithDefaults("DB_URI", "mongodb://localhost:27017"),
		DBName:     util.MapEnvWithDefaults("DB_NAME", "auth"),
		DBUser:     util.MapEnvWithDefaults("DB_USER", ""),
		DBPassword: util.MapEnvWithDefaults("DB_PASSWORD", ""),
		DBTimeout:  dBTimeout,
	}
}
