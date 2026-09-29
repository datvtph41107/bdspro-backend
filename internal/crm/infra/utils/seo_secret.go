package utils

import (
	"context"
	"crm/config"
	"crypto/subtle"

	"google.golang.org/grpc/metadata"
)

func ValidateSeoInternalSecret(ctx context.Context) bool {
	expected := config.AppProperties.Seo.InternalSecretKey
	if expected == "" {
		return false
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false
	}

	values := md.Get("x-seo-secret-key")
	if len(values) == 0 {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(values[0]), []byte(expected)) == 1
}