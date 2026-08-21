package config

import (
	"strings"
	"testing"
	"time"
)

func validRuntimeConfig() RuntimeConfig {
	return RuntimeConfig{
		HTTPPort:     "8002",
		DatabaseDSN:  "host=db.example user=file password=runtime-secret dbname=file",
		JWTKey:       "jwt-key-from-secret-store",
		Auth:         AuthConfig{GRPCAddress: "auth:8216", Timeout: 2 * time.Second},
		Hub:          HubConfig{GRPCAddress: "hub:8280", Timeout: 5 * time.Second},
		InternalGRPC: InternalGRPCConfig{Enabled: false, Port: "8203"},
		FileService: FileServiceConfig{
			SignatureKey:    "signature-secret-from-secret-store",
			XorCryptKey:     "xor-secret-from-secret-store",
			SignatureExpire: 5 * time.Minute,
			ServiceAuthKey:  "service-auth-secret-from-secret-store",
		},
	}
}

func TestRuntimeConfigValidateAcceptsExplicitDistinctSecrets(t *testing.T) {
	if err := validRuntimeConfig().Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestRuntimeConfigValidateFailsClosedForMissingOrPlaceholderSecrets(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*RuntimeConfig)
		want   string
	}{
		{name: "missing database dsn", mutate: func(c *RuntimeConfig) { c.DatabaseDSN = "" }, want: "database.dsn is required"},
		{name: "missing auth address", mutate: func(c *RuntimeConfig) { c.Auth.GRPCAddress = "" }, want: "auth.grpc_address is required"},
		{name: "missing auth timeout", mutate: func(c *RuntimeConfig) { c.Auth.Timeout = 0 }, want: "auth.timeout must be positive"},
		{name: "missing hub address", mutate: func(c *RuntimeConfig) { c.Hub.GRPCAddress = "" }, want: "hub.grpc_address is required"},
		{name: "missing hub timeout", mutate: func(c *RuntimeConfig) { c.Hub.Timeout = 0 }, want: "hub.timeout must be positive"},
		{name: "enabled grpc missing port", mutate: func(c *RuntimeConfig) { c.InternalGRPC = InternalGRPCConfig{Enabled: true} }, want: "internal_grpc.port is required"},
		{name: "invalid grpc port", mutate: func(c *RuntimeConfig) { c.InternalGRPC.Port = "70000" }, want: "between 1 and 65535"},
		{name: "missing jwt key", mutate: func(c *RuntimeConfig) { c.JWTKey = "" }, want: "jwt.key-generate is required"},
		{name: "legacy jwt placeholder", mutate: func(c *RuntimeConfig) { c.JWTKey = "key-generate-jwt-token-authentication" }, want: "known insecure placeholder"},
		{name: "missing signature", mutate: func(c *RuntimeConfig) { c.FileService.SignatureKey = "" }, want: "file.signature_key is required"},
		{name: "legacy signature placeholder", mutate: func(c *RuntimeConfig) { c.FileService.SignatureKey = "key-signature" }, want: "known insecure placeholder"},
		{name: "legacy xor placeholder", mutate: func(c *RuntimeConfig) { c.FileService.XorCryptKey = "key-encode-decode" }, want: "known insecure placeholder"},
		{name: "tracked service auth placeholder", mutate: func(c *RuntimeConfig) {
			c.FileService.ServiceAuthKey = "internal-service-secret-key-change-in-production"
		}, want: "known insecure placeholder"},
		{name: "missing expiry", mutate: func(c *RuntimeConfig) { c.FileService.SignatureExpire = 0 }, want: "must be positive"},
		{name: "shared secret", mutate: func(c *RuntimeConfig) { c.FileService.XorCryptKey = c.FileService.SignatureKey }, want: "must be distinct"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := validRuntimeConfig()
			tc.mutate(&cfg)
			err := cfg.Validate()
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("Validate() error = %v, want containing %q", err, tc.want)
			}
		})
	}
}
