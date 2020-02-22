package util

import (
	"os"
)

// MapEnvWithDefaults returns a default value if the environment
// variable is empty
func MapEnvWithDefaults(envKey string, defaults string) string {
	v := os.Getenv(envKey)
	if v == "" {
		if defaults == "" {
			panic("defaults is not specified")
		}
		return defaults
	}
	return v
}
