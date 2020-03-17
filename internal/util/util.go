package util

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/grpc/metadata"
)

// MapEnvWithDefaults returns a default value if the environment
// variable is empty
func MapEnvWithDefaults(envKey string, defaults string) string {
	v := os.Getenv(envKey)
	if v == "" {
		if defaults == "" {
			panic(fmt.Sprint(envKey, " defaults is not specified"))
		}
		return defaults
	}
	return v
}

// GetMetadataValue is a helper function that returns value stored in metadata
func GetMetadataValue(ctx context.Context, key string) string {
	mdC, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if len(mdC.Get(key)) > 0 {
			return mdC.Get(key)[0]
		}
	}
	return ""
}

func GetMetadata(ctx context.Context) map[string]string {
	m := make(map[string]string)
	mdC, ok := metadata.FromIncomingContext(ctx)
	if ok && mdC.Len() > 0 {
		for k, v := range mdC {
			if len(mdC.Get(k)) > 0 {
				m[k] = v[0]
			}
		}
	}
	return m
}
