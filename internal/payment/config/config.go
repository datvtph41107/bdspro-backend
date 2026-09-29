package config

import (
	"common/configloader"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Address         string
	ShutdownTimeout time.Duration
}

type DatabaseConfig struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type ClientConfig struct {
	Address string
	Timeout time.Duration
}

type CommerceConfig struct {
	OrderTTL         time.Duration
	FulfillmentRetry time.Duration
	FulfillmentLease time.Duration
	FulfillmentPoll  time.Duration
}

type RabbitConfig struct {
	URL             string
	Exchange        string
	PublisherLease  time.Duration
	PublisherPoll   time.Duration
	OutboxRetryBase time.Duration
	OutboxRetryMax  time.Duration
	ReconnectBase   time.Duration
	ReconnectMax    time.Duration
}

type ProviderConfig struct {
	SepayAPIKey string
}

type Config struct {
	Environment  string
	Server       ServerConfig
	Database     DatabaseConfig
	User         ClientConfig
	Auth         ClientConfig
	Notification ClientConfig
	Commerce     CommerceConfig
	Rabbit       RabbitConfig
	Provider     ProviderConfig
}

func Load() (Config, error) {
	selection, err := configloader.ResolveRuntimeSelection()
	if err != nil {
		return Config{}, err
	}
	v := viper.New()
	v.SetConfigType("yaml")
	file := resolveConfigFile(os.Getenv("CONFIG_FILE"))
	v.SetConfigFile(file)
	// YAML là default không nhạy cảm; environment là override/secret. Thiếu hoặc
	// lỗi YAML phải fail sớm để dev và container dùng cùng một contract.
	if err := v.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("load Payment config %s: %w", file, err)
	}

	host := envOr("PAYMENT_GRPC_HOST", v.GetString("host"), "0.0.0.0")
	port := envOr("PAYMENT_GRPC_PORT", v.GetString("port"), "8205")
	address := envOr("PAYMENT_GRPC_ADDRESS", "", host+":"+port)
	dsn := firstNonBlank(os.Getenv("PAYMENT_DATABASE_URL"), os.Getenv("PAYMENT_DATABASE_DSN"), v.GetString("database.url"))
	if strings.TrimSpace(dsn) == "" || strings.Contains(strings.ToUpper(dsn), "CHANGE_ME") {
		return Config{}, errors.New("PAYMENT_DATABASE_URL is required")
	}

	outboxRetryBase, outboxRetryMax, reconnectBase, reconnectMax := rabbitTimingConfig()

	cfg := Config{
		Environment: selection.Environment,
		Server: ServerConfig{
			Address:         address,
			ShutdownTimeout: durationEnv("PAYMENT_SHUTDOWN_TIMEOUT", 15*time.Second),
		},
		Database: DatabaseConfig{
			DSN:             dsn,
			MaxOpenConns:    intEnv("PAYMENT_DATABASE_MAX_OPEN_CONNS", 30),
			MaxIdleConns:    intEnv("PAYMENT_DATABASE_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: durationEnv("PAYMENT_DATABASE_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		User: ClientConfig{
			Address: envOr("QHPRO_USER_GRPC_ADDRESS", "", "localhost:8201"),
			Timeout: durationEnv("PAYMENT_USER_RPC_TIMEOUT", 5*time.Second),
		},
		Auth: ClientConfig{
			Address: envOr("QHPRO_AUTH_GRPC_ADDRESS", "", "localhost:8216"),
			Timeout: durationEnv("PAYMENT_AUTH_RPC_TIMEOUT", 5*time.Second),
		},
		Notification: ClientConfig{
			Address: envOr("QHPRO_NOTIFICATION_GRPC_ADDRESS", "", "localhost:8204"),
			Timeout: durationEnv("PAYMENT_NOTIFICATION_RPC_TIMEOUT", 5*time.Second),
		},
		Commerce: CommerceConfig{
			OrderTTL:         durationEnv("PAYMENT_ORDER_TTL", 30*time.Minute),
			FulfillmentRetry: durationEnv("PAYMENT_FULFILLMENT_RETRY_DELAY", 5*time.Second),
			FulfillmentLease: durationEnv("PAYMENT_FULFILLMENT_CLAIM_LEASE", 30*time.Second),
			FulfillmentPoll:  durationEnv("PAYMENT_FULFILLMENT_POLL_INTERVAL", time.Second),
		},
		Rabbit: RabbitConfig{
			URL:             strings.TrimSpace(os.Getenv("PAYMENT_RABBIT_URL")),
			Exchange:        envOr("PAYMENT_EVENT_EXCHANGE", "", "qhpro.payment.events"),
			PublisherLease:  durationEnv("PAYMENT_OUTBOX_CLAIM_LEASE", 30*time.Second),
			PublisherPoll:   durationEnv("PAYMENT_OUTBOX_POLL_INTERVAL", time.Second),
			OutboxRetryBase: outboxRetryBase,
			OutboxRetryMax:  outboxRetryMax,
			ReconnectBase:   reconnectBase,
			ReconnectMax:    reconnectMax,
		},
		Provider: ProviderConfig{SepayAPIKey: strings.TrimSpace(os.Getenv("SEPAY_API_KEY"))},
	}
	if cfg.Server.Address == "" || cfg.Database.MaxOpenConns <= 0 || cfg.Database.MaxIdleConns < 0 ||
		cfg.Commerce.OrderTTL <= 0 || cfg.Commerce.FulfillmentLease <= 0 || cfg.Commerce.FulfillmentPoll <= 0 ||
		cfg.Rabbit.PublisherLease <= 0 || cfg.Rabbit.PublisherPoll <= 0 ||
		cfg.Rabbit.OutboxRetryBase <= 0 || cfg.Rabbit.OutboxRetryMax < cfg.Rabbit.OutboxRetryBase ||
		cfg.Rabbit.ReconnectBase <= 0 || cfg.Rabbit.ReconnectMax < cfg.Rabbit.ReconnectBase {
		return Config{}, fmt.Errorf("invalid Payment runtime configuration")
	}
	return cfg, nil
}

// resolveConfigFile uses one committed runtime default. CONFIG_FILE is the
// explicit operator escape hatch for an externally mounted configuration.
func rabbitTimingConfig() (outboxRetryBase, outboxRetryMax, reconnectBase, reconnectMax time.Duration) {
	legacyRetry := durationEnv("PAYMENT_OUTBOX_RETRY_DELAY", 5*time.Second)
	outboxRetryBase = durationEnv("PAYMENT_OUTBOX_RETRY_BASE", legacyRetry)
	outboxRetryMax = durationEnv("PAYMENT_OUTBOX_RETRY_MAX", time.Minute)
	reconnectBase = durationEnv("PAYMENT_RABBIT_RECONNECT_BASE", legacyRetry)
	reconnectMax = durationEnv("PAYMENT_RABBIT_RECONNECT_MAX", 30*time.Second)
	return outboxRetryBase, outboxRetryMax, reconnectBase, reconnectMax
}

func resolveConfigFile(override string) string {
	if file := strings.TrimSpace(override); file != "" {
		return file
	}
	return "config/runtime.yml"
}

func firstNonBlank(values ...string) string {
	for _, value := range values {
		if v := strings.TrimSpace(value); v != "" {
			return v
		}
	}
	return ""
}

func envOr(key, compatibility, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	if value := strings.TrimSpace(compatibility); value != "" && !strings.Contains(strings.ToUpper(value), "CHANGE_ME") {
		return value
	}
	return fallback
}

func durationEnv(key string, fallback time.Duration) time.Duration {
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

func intEnv(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return fallback
	}
	return parsed
}
