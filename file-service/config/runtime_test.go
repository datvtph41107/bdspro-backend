package config

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func validRuntimeConfig() RuntimeConfig {
	return RuntimeConfig{
		HTTPPort:    "8002",
		DatabaseDSN: "host=db.example user=file password=runtime-secret dbname=file",
		JWTKey:      "jwt-key-from-secret-store",
		Auth:        AuthConfig{GRPCAddress: "auth:8216", Timeout: 2 * time.Second},
		Hub:         HubConfig{GRPCAddress: "hub:8280", Timeout: 5 * time.Second},
		FileService: FileServiceConfig{
			StorageRoot:        "/var/lib/qhpro/files",
			CORSAllowedOrigins: []string{"https://qhpro.vn"},
			SignatureKey:       "signature-secret-from-secret-store",
			XorCryptKey:        "xor-secret-from-secret-store",
			SignatureExpire:    5 * time.Minute,
			ServiceAuthKey:     "service-auth-secret-from-secret-store",
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
		{name: "missing jwt key", mutate: func(c *RuntimeConfig) { c.JWTKey = "" }, want: "jwt.key-generate is required"},
		{name: "legacy jwt placeholder", mutate: func(c *RuntimeConfig) { c.JWTKey = "key-generate-jwt-token-authentication" }, want: "known insecure placeholder"},
		{name: "missing storage root", mutate: func(c *RuntimeConfig) { c.FileService.StorageRoot = "" }, want: "file.storage_root is required"},
		{name: "missing CORS origins", mutate: func(c *RuntimeConfig) { c.FileService.CORSAllowedOrigins = nil }, want: "file.cors_allowed_origins is required"},
		{name: "invalid CORS origin", mutate: func(c *RuntimeConfig) { c.FileService.CORSAllowedOrigins = []string{"*"} }, want: "invalid origin"},
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

func TestConfigStringListNormalizesCommaSeparatedRuntimeOverride(t *testing.T) {
	// Environment values are exposed to Viper as one string. The config owner
	// normalizes that representation so deployed FILE_CORS_ALLOWED_ORIGINS can
	// use the conventional comma-separated form.
	viper.Set("file.cors_allowed_origins", "https://qhpro.vn, https://admin.bdspro.com,https://qhpro.vn")
	t.Cleanup(func() { viper.Set("file.cors_allowed_origins", nil) })

	got := configStringList("file.cors_allowed_origins")
	want := []string{"https://qhpro.vn", "https://admin.bdspro.com"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("configStringList() = %#v, want %#v", got, want)
	}
}
