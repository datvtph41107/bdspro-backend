package config

import (
	_foundationconfig "common/config"
	"common/trust"

	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Runtime is the typed Gateway composition configuration. It deliberately
// keeps service topology in Gateway instead of shared infrastructure.
type Runtime struct {
	Environment   string
	UseServiceDNS bool
	Server        ServerRuntime  `yaml:"server"`
	Service       ServiceRuntime `yaml:"service"`
	JWT           JWTRuntime     `yaml:"jwt"`
	HTTP          HTTPRuntime    `yaml:"http"`
	Trust         trust.Config   `yaml:"-"`
}

type ServerRuntime struct {
	Port int `yaml:"port"`
}

type ServiceRuntime struct {
	Domain string         `yaml:"domain"`
	Port   map[string]int `yaml:"port"`
}

// JWTRuntime contains only the JWT input owned by Gateway: signature verification.
// Token issuance/expiry policy remains owned by the services that issue tokens.
type JWTRuntime struct {
	VerificationKey string `yaml:"key-generate"`
}

// HTTPRuntime contains Gateway-owned HTTP edge policy.
type HTTPRuntime struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
}

// LoadRuntime loads the caller-owned typed Gateway runtime without installing
// process configuration into a shared/global registry.
func LoadRuntime() (Runtime, error) {
	runtimeEnv := strings.TrimSpace(os.Getenv("ENV_RUNTIME"))
	fileEnv := runtimeEnv
	if fileEnv == "" {
		fileEnv = "local"
	}

	path := "config/" + fileEnv + ".yml"
	var cfg Runtime
	if err := _foundationconfig.LoadYAML(path, &cfg); err != nil {
		return Runtime{}, err
	}
	cfg.Environment = fileEnv
	cfg.UseServiceDNS = runtimeEnv != ""
	cfg.Trust = loadTrustConfig()

	// Preserve the historical Viper environment override for jwt.key-generate
	// without installing Gateway configuration into a process-global registry.
	if value := strings.TrimSpace(os.Getenv("JWT_KEY_GENERATE")); value != "" {
		cfg.JWT.VerificationKey = value
	}

	if err := cfg.Validate(); err != nil {
		return Runtime{}, err
	}
	return cfg, nil
}

func (c Runtime) Validate() error {
	if c.Server.Port <= 0 {
		return fmt.Errorf("gateway HTTP port is required")
	}
	if strings.TrimSpace(c.JWT.VerificationKey) == "" {
		return fmt.Errorf("gateway JWT verification key is required")
	}
	for _, name := range []string{"notification", "payment", "bdspro", "crm", "organization", "social", "user", "chat", "auth", "assistant", "hub", "tqd"} {
		if _, err := c.Endpoint(name); err != nil {
			return err
		}
	}
	if len(c.HTTP.AllowedOrigins) == 0 {
		return fmt.Errorf("gateway HTTP allowed_origins is required")
	}
	for index, origin := range c.HTTP.AllowedOrigins {
		if strings.TrimSpace(origin) == "" {
			return fmt.Errorf("gateway HTTP allowed_origins[%d] is empty", index)
		}
	}

	return nil
}

func (c Runtime) Endpoint(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", fmt.Errorf("service name is required")
	}

	port := c.Service.Port[name]
	if port <= 0 {
		return "", fmt.Errorf("gateway service port %q is required", name)
	}

	host := strings.TrimSpace(c.Service.Domain)
	if c.UseServiceDNS {
		host = name
	}
	if host == "" {
		host = "localhost"
	}

	return host + ":" + strconv.Itoa(port), nil
}

func loadTrustConfig() trust.Config {
	mode := trust.Mode(strings.ToLower(strings.TrimSpace(os.Getenv("QHPRO_TRUSTED_METADATA_MODE"))))
	if mode != trust.ModeEnforce {
		mode = trust.ModeAudit
	}

	// 	QHPRO_INTERNAL_METADATA_VERIFICATION_KEYS='{
	//   "auth-service":"auth-inbound-secret",
	//   "user-service":"user-inbound-secret",
	//   "tqd-service":"tqd-inbound-secret"
	// }'
	verificationSecrets, configured := verificationSecretsFromEnv("QHPRO_INTERNAL_METADATA_VERIFICATION_KEYS")
	return trust.Config{
		ServiceID:                  firstNonEmpty(os.Getenv("QHPRO_SERVICE_ID"), os.Getenv("SERVICE_NAME"), "gateway-service"),
		Secret:                     strings.TrimSpace(os.Getenv("QHPRO_INTERNAL_METADATA_SECRET")),
		VerificationSecrets:        verificationSecrets,
		VerificationKeysConfigured: configured,
		Mode:                       mode,
		MaxAge:                     durationFromEnv("QHPRO_INTERNAL_METADATA_MAX_AGE", 30*time.Second),
		ClockSkew:                  durationFromEnv("QHPRO_INTERNAL_METADATA_CLOCK_SKEW", 5*time.Second),
	}
}

func durationFromEnv(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil || parsed <= 0 {
		return fallback
	}

	return parsed
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}

func verificationSecretsFromEnv(key string) (map[string]string, bool) {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return nil, false
	}

	var values map[string]string
	if err := json.Unmarshal([]byte(raw), &values); err != nil {
		return nil, true
	}

	result := make(map[string]string, len(values))
	for serviceID, secret := range values {
		serviceID = strings.TrimSpace(serviceID)
		secret = strings.TrimSpace(secret)
		if serviceID != "" && secret != "" {
			result[serviceID] = secret
		}
	}
	if len(result) == 0 {
		return nil, true
	}

	return result, true
}
