package config

import (
	_foundationconfig "common/config"
	"common/configloader"
	"common/rpc"

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
	Server        ServerRuntime       `yaml:"server"`
	Service       ServiceRuntime      `yaml:"service"`
	JWT           JWTRuntime          `yaml:"jwt"`
	HTTP          HTTPRuntime         `yaml:"http"`
	RPCTransport  rpc.TransportConfig `yaml:"-"`
}

type ServerRuntime struct {
	Port int `yaml:"port"`
}

type ServiceRuntime struct {
	Domain   string            `yaml:"domain"`
	Host     map[string]string `yaml:"host"`
	Port     map[string]int    `yaml:"port"`
	HTTPHost map[string]string `yaml:"http_host"`
	HTTPPort map[string]int    `yaml:"http_port"`
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
	selection, err := configloader.ResolveRuntimeSelection()
	if err != nil {
		return Runtime{}, err
	}
	path := strings.TrimSpace(os.Getenv("QHPRO_CONFIG_FILE"))
	if path == "" {
		path = strings.TrimSpace(os.Getenv("CONFIG_FILE"))
	}
	if path == "" {
		path = "config/runtime.yml"
	}
	var cfg Runtime
	if err := _foundationconfig.LoadYAML(path, &cfg); err != nil {
		return Runtime{}, err
	}
	cfg.Environment = selection.Environment
	// Execution mode controls only address wiring. It never selects a second
	// source configuration file.
	cfg.UseServiceDNS = usesServiceDNS(selection.ExecutionMode)
	cfg.RPCTransport = loadRPCTransportConfig()

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

func usesServiceDNS(executionMode string) bool {
	return strings.EqualFold(strings.TrimSpace(executionMode), "container")
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
	return c.serviceEndpoint(name, c.Service.Host, c.Service.Port, "gRPC")
}

// HTTPEndpoint returns only endpoints whose serving process actually speaks
// HTTP. It deliberately does not fall back to the gRPC port map.
func (c Runtime) HTTPEndpoint(name string) (string, error) {
	return c.serviceEndpoint(name, c.Service.HTTPHost, c.Service.HTTPPort, "HTTP")
}

func (c Runtime) serviceEndpoint(name string, hosts map[string]string, ports map[string]int, protocol string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", fmt.Errorf("service name is required")
	}

	port := ports[name]
	if port <= 0 {
		return "", fmt.Errorf("gateway service %s port %q is required", protocol, name)
	}

	host := strings.TrimSpace(c.Service.Domain)
	if c.UseServiceDNS {
		// Host alias diễn đạt đúng serving owner. Ví dụ API Auth public do User
		// phục vụ, trong khi auth-service:8216 chỉ giữ AuthInternal.
		host = strings.TrimSpace(hosts[name])
		if host == "" {
			host = name
		}
	}
	if host == "" {
		host = "localhost"
	}

	return host + ":" + strconv.Itoa(port), nil
}

func loadRPCTransportConfig() rpc.TransportConfig {
	verificationSecrets, configured :=
		verificationSecretsFromEnv("QHPRO_INTERNAL_METADATA_VERIFICATION_KEYS")

	return rpc.TransportConfig{
		ServiceAssertion: rpc.ServiceAssertionConfig{
			ServiceID: firstNonEmpty(
				os.Getenv("QHPRO_SERVICE_ID"),
				os.Getenv("SERVICE_NAME"),
				"gateway-service",
			),
			Secret: strings.TrimSpace(
				os.Getenv("QHPRO_INTERNAL_METADATA_SECRET"),
			),
			VerificationSecrets:        verificationSecrets,
			VerificationKeysConfigured: configured,
			MaxAge: durationFromEnv(
				"QHPRO_INTERNAL_METADATA_MAX_AGE",
				30*time.Second,
			),
			ClockSkew: durationFromEnv(
				"QHPRO_INTERNAL_METADATA_CLOCK_SKEW",
				5*time.Second,
			),
		},
		RequireServiceAssertion: strings.EqualFold(
			strings.TrimSpace(
				os.Getenv("QHPRO_TRUSTED_METADATA_MODE"),
			),
			"enforce",
		),
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
