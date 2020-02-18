package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

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

func mapEnvWithDefaults(envKey string, defaults string) string {
	v := os.Getenv(envKey)
	if v == "" {
		if defaults == "" {
			panic("defaults is not specified")
		}
		return defaults
	}
	return v
}

// LoadEnv loads environment variables for Application
func loadEnv() *EnvConfig {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env not loaded", err)
	}

	// Convert to int
	sec, err := strconv.ParseInt(mapEnvWithDefaults("DB_TIMEOUT", "3"), 10, 64)
	if err != nil {
		fmt.Printf("Error parsing DB_TIMEOUT: %v\nWill fallback to default value", err)
		sec = 3
	}

	dBTimeout := time.Duration(sec) * time.Second

	return &EnvConfig{
		AppEnv:     mapEnvWithDefaults("APP_ENV", "development"),
		Production: mapEnvWithDefaults("APP_ENV", "development") == "true",
		Address:    mapEnvWithDefaults("ADDRESS", "5000"),
		DBUri:      mapEnvWithDefaults("DB_URI", "mongodb://localhost:27017"),
		DBName:     mapEnvWithDefaults("DB_NAME", "auth"),
		DBUser:     mapEnvWithDefaults("DB_USER", ""),
		DBPassword: mapEnvWithDefaults("DB_PASSWORD", ""),
		DBTimeout:  dBTimeout,
	}
}
