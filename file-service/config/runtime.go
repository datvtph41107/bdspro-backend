package config

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var insecureSecretLiterals = map[string]struct{}{
	"key-signature":         {},
	"key-encode-decode":     {},
	"default-signature-key": {},
	"default-xor-key":       {},
	"service-internal-key":  {},
	"internal-service-secret-key-change-in-production": {},
	"key-generate-jwt-token-authentication":            {},
}

// RuntimeConfig is the validated production composition input for the
// File-service runtime. Environment/YAML lookup ends here; business services
// receive typed values and never read Viper directly.
type RuntimeConfig struct {
	HTTPPort    string
	DatabaseDSN string
	JWTKey      string
	Auth        AuthConfig
	Hub         HubConfig
	FileService FileServiceConfig
}

type AuthConfig struct {
	GRPCAddress string
	Timeout     time.Duration
}

type HubConfig struct {
	GRPCAddress string
	Timeout     time.Duration
}

type FileServiceConfig struct {
	StorageRoot        string
	CORSAllowedOrigins []string
	SignatureKey       string
	XorCryptKey        string
	SignatureExpire    time.Duration
	ServiceAuthKey     string
}

// CurrentRuntimeConfig resolves Viper only after LoadConfig has loaded the
// environment-specific YAML. Secret values may be supplied by environment
// variables through Viper's dot-to-underscore mapping.
func CurrentRuntimeConfig() (RuntimeConfig, error) {
	cfg := RuntimeConfig{
		HTTPPort:    strings.TrimSpace(viper.GetString("server.port")),
		DatabaseDSN: strings.TrimSpace(viper.GetString("database.dsn")),
		JWTKey:      strings.TrimSpace(viper.GetString("jwt.key-generate")),
		Auth: AuthConfig{
			GRPCAddress: strings.TrimSpace(viper.GetString("auth.grpc_address")),
			Timeout:     viper.GetDuration("auth.timeout"),
		},
		Hub: HubConfig{
			GRPCAddress: strings.TrimSpace(viper.GetString("hub.grpc_address")),
			Timeout:     viper.GetDuration("hub.timeout"),
		},
		FileService: FileServiceConfig{
			StorageRoot:        strings.TrimSpace(viper.GetString("file.storage_root")),
			CORSAllowedOrigins: configStringList("file.cors_allowed_origins"),
			SignatureKey:       strings.TrimSpace(viper.GetString("file.signature_key")),
			XorCryptKey:        strings.TrimSpace(viper.GetString("file.xor_crypt_key")),
			SignatureExpire:    viper.GetDuration("file.signature_expire"),
			ServiceAuthKey:     strings.TrimSpace(viper.GetString("service.auth_key")),
		},
	}
	if err := cfg.Validate(); err != nil {
		return RuntimeConfig{}, err
	}
	return cfg, nil
}

func (c RuntimeConfig) Validate() error {
	if c.HTTPPort == "" {
		return errors.New("file-service server.port is required")
	}
	if strings.TrimSpace(c.DatabaseDSN) == "" {
		return errors.New("file-service database.dsn is required")
	}
	if err := validateSecret("jwt.key-generate", c.JWTKey); err != nil {
		return err
	}
	if strings.TrimSpace(c.Auth.GRPCAddress) == "" {
		return errors.New("file-service auth.grpc_address is required")
	}
	if c.Auth.Timeout <= 0 {
		return errors.New("file-service auth.timeout must be positive")
	}
	if strings.TrimSpace(c.Hub.GRPCAddress) == "" {
		return errors.New("file-service hub.grpc_address is required")
	}
	if c.Hub.Timeout <= 0 {
		return errors.New("file-service hub.timeout must be positive")
	}
	if strings.TrimSpace(c.FileService.StorageRoot) == "" {
		return errors.New("file.storage_root is required")
	}
	if err := validateCORSOrigins(c.FileService.CORSAllowedOrigins); err != nil {
		return err
	}
	if err := validateSecret("file.signature_key", c.FileService.SignatureKey); err != nil {
		return err
	}
	if err := validateSecret("file.xor_crypt_key", c.FileService.XorCryptKey); err != nil {
		return err
	}
	if err := validateSecret("service.auth_key", c.FileService.ServiceAuthKey); err != nil {
		return err
	}
	if c.FileService.SignatureExpire <= 0 {
		return errors.New("file.signature_expire must be positive")
	}
	if c.FileService.SignatureKey == c.FileService.XorCryptKey ||
		c.FileService.SignatureKey == c.FileService.ServiceAuthKey ||
		c.FileService.XorCryptKey == c.FileService.ServiceAuthKey {
		return errors.New("file-service security keys must be distinct")
	}
	return nil
}

func validateSecret(name, value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("%s is required", name)
	}
	if _, insecure := insecureSecretLiterals[value]; insecure {
		return fmt.Errorf("%s uses a known insecure placeholder", name)
	}
	return nil
}

func configStringList(key string) []string {
	var values []string
	switch raw := viper.Get(key).(type) {
	case []string:
		values = append(values, raw...)
	case []interface{}:
		for _, item := range raw {
			if value, ok := item.(string); ok {
				values = append(values, value)
			}
		}
	case string:
		values = append(values, strings.Split(raw, ",")...)
	}

	if len(values) == 0 {
		values = append(values, viper.GetStringSlice(key)...)
	}

	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for _, value := range values {
		for _, part := range strings.Split(value, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			if _, exists := seen[part]; exists {
				continue
			}
			seen[part] = struct{}{}
			out = append(out, part)
		}
	}
	return out
}

func validateCORSOrigins(origins []string) error {
	if len(origins) == 0 {
		return errors.New("file.cors_allowed_origins is required")
	}
	for _, origin := range origins {
		parsed, err := url.ParseRequestURI(origin)
		if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.RawQuery != "" || parsed.Fragment != "" {
			return fmt.Errorf("file.cors_allowed_origins contains invalid origin %q", origin)
		}
		if parsed.Path != "" && parsed.Path != "/" {
			return fmt.Errorf("file.cors_allowed_origins contains non-origin path %q", origin)
		}
	}
	return nil
}
