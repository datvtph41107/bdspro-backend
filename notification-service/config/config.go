package config

import (
	"common/configloader"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	C_SESSION_ID                  = "C_SESSION_ID"
	defaultBusinessEventsExchange = "qhpro.payment.events"
)

// UserProperties chứa cấu hình của ứng dụng
type UserProperties struct {
	Server struct {
		TCPPort int `mapstructure:"tcp_port"`
	} `mapstructure:"server"`
	Database struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Name     string `mapstructure:"name"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
	} `mapstructure:"database"`
	Jwt struct {
		SecretKey           string   `mapstructure:"key-generate"`
		TokenPrefix         string   `mapstructure:"tokenPrefix"`
		TokenExpirationDays int      `mapstructure:"tokenExpirationAfterDays"`
		AccessExpMinutes    int      `mapstructure:"accessExpAfterMinutes"`
		RefreshExpMinutes   uint64   `mapstructure:"refreshExpAfterMinutes"`
		ListPermit          []string `mapstructure:"listPermit"`
		AuthorizationHeader string   `mapstructure:"authorizationHeader"`
	} `mapstructure:"jwt"`

	Firebase struct {
		File    string `mapstructure:"file"`
		FileURL string `mapstructure:"fileUrl"`
	} `mapstructure:"firebase"`
	Redis struct {
		Host     string `mapstructure:"host"`
		Password string `mapstructure:"pass"`
	} `mapstructure:"redis"`
}

var AppProperties UserProperties

// LoadConfig đọc cấu hình từ file config.yml
func LoadConfig() error {
	_, err := configloader.LoadRuntimeYML()
	if err != nil {
		return err
	}

	if err := viper.Unmarshal(&AppProperties); err != nil {
		return fmt.Errorf("unmarshal notification config: %w", err)
	}

	return nil
}

// RuntimeConfig contains process-owned technical configuration that must not be
// read from business packages. Values may come from the existing YAML/Viper
// configuration or explicit environment overrides used by containers.
type RuntimeConfig struct {
	ServerPort             int
	DatabaseDSN            string
	RedisAddress           string
	RedisPassword          string
	RedisDB                int
	UserGRPCTarget         string
	RabbitURL              string
	BusinessEventsExchange string
	DeliveryWorkerID       string
	DeliveryPollInterval   time.Duration
	DeliveryLease          time.Duration
	ShutdownTimeout        time.Duration
}

func LoadRuntimeConfig() (RuntimeConfig, error) {
	if err := LoadConfig(); err != nil {
		return RuntimeConfig{}, err
	}
	cfg := RuntimeConfig{
		ServerPort:             AppProperties.Server.TCPPort,
		DatabaseDSN:            databaseDSN(),
		RedisAddress:           strings.TrimSpace(os.Getenv("NOTIFICATION_REDIS_ADDR")),
		RedisPassword:          os.Getenv("NOTIFICATION_REDIS_PASSWORD"),
		RedisDB:                viper.GetInt("redis.db"),
		UserGRPCTarget:         strings.TrimSpace(os.Getenv("QHPRO_USER_GRPC_ADDR")),
		RabbitURL:              strings.TrimSpace(os.Getenv("QHPRO_RABBITMQ_URL")),
		BusinessEventsExchange: strings.TrimSpace(os.Getenv("QHPRO_BUSINESS_EVENTS_EXCHANGE")),
		DeliveryWorkerID:       strings.TrimSpace(os.Getenv("NOTIFICATION_DELIVERY_WORKER_ID")),
		DeliveryPollInterval:   time.Second,
		DeliveryLease:          30 * time.Second,
		ShutdownTimeout:        10 * time.Second,
	}
	if raw := strings.TrimSpace(os.Getenv("NOTIFICATION_GRPC_PORT")); raw != "" {
		port, err := strconv.Atoi(raw)
		if err != nil || port <= 0 || port > 65535 {
			return RuntimeConfig{}, fmt.Errorf("invalid NOTIFICATION_GRPC_PORT %q", raw)
		}
		cfg.ServerPort = port
	}
	if cfg.RedisAddress == "" {
		cfg.RedisAddress = strings.TrimSpace(AppProperties.Redis.Host)
	}
	if cfg.RedisPassword == "" {
		cfg.RedisPassword = AppProperties.Redis.Password
	}
	redisDB, err := parseRedisDBOverride(os.Getenv("NOTIFICATION_REDIS_DB"), cfg.RedisDB)
	if err != nil {
		return RuntimeConfig{}, err
	}
	cfg.RedisDB = redisDB
	if cfg.UserGRPCTarget == "" {
		cfg.UserGRPCTarget = strings.TrimSpace(viper.GetString("rpc.user.address"))
	}
	if cfg.BusinessEventsExchange == "" {
		cfg.BusinessEventsExchange = defaultBusinessEventsExchange
	}
	if cfg.DeliveryWorkerID == "" {
		cfg.DeliveryWorkerID = "notification-delivery"
	}
	if raw := strings.TrimSpace(os.Getenv("NOTIFICATION_DELIVERY_POLL_INTERVAL")); raw != "" {
		d, err := time.ParseDuration(raw)
		if err != nil || d <= 0 {
			return RuntimeConfig{}, fmt.Errorf("invalid NOTIFICATION_DELIVERY_POLL_INTERVAL %q", raw)
		}
		cfg.DeliveryPollInterval = d
	}
	if raw := strings.TrimSpace(os.Getenv("NOTIFICATION_DELIVERY_LEASE")); raw != "" {
		d, err := time.ParseDuration(raw)
		if err != nil || d <= 0 {
			return RuntimeConfig{}, fmt.Errorf("invalid NOTIFICATION_DELIVERY_LEASE %q", raw)
		}
		cfg.DeliveryLease = d
	}
	if raw := strings.TrimSpace(os.Getenv("NOTIFICATION_SHUTDOWN_TIMEOUT")); raw != "" {
		d, err := time.ParseDuration(raw)
		if err != nil || d <= 0 {
			return RuntimeConfig{}, fmt.Errorf("invalid NOTIFICATION_SHUTDOWN_TIMEOUT %q", raw)
		}
		cfg.ShutdownTimeout = d
	}
	if strings.TrimSpace(cfg.DatabaseDSN) == "" {
		return RuntimeConfig{}, fmt.Errorf("notification database DSN is required")
	}
	if strings.TrimSpace(cfg.RedisAddress) == "" {
		return RuntimeConfig{}, fmt.Errorf("notification Redis address is required")
	}
	if strings.TrimSpace(cfg.UserGRPCTarget) == "" {
		return RuntimeConfig{}, fmt.Errorf("notification User gRPC target is required")
	}
	if cfg.ServerPort <= 0 {
		return RuntimeConfig{}, fmt.Errorf("notification gRPC server port is required")
	}
	return cfg, nil
}

func parseRedisDBOverride(raw string, fallback int) (int, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		if fallback < 0 || fallback > 15 {
			return 0, fmt.Errorf("invalid redis.db %d: must be between 0 and 15", fallback)
		}
		return fallback, nil
	}
	db, err := strconv.Atoi(raw)
	if err != nil || db < 0 || db > 15 {
		return 0, fmt.Errorf("invalid NOTIFICATION_REDIS_DB %q: must be between 0 and 15", raw)
	}
	return db, nil
}

type EventingConfig struct {
	DatabaseDSN     string
	RabbitURL       string
	Exchange        string
	ShutdownTimeout time.Duration
}

func LoadEventingConfig() (EventingConfig, error) {
	if err := LoadConfig(); err != nil {
		return EventingConfig{}, err
	}
	cfg := EventingConfig{
		DatabaseDSN:     databaseDSN(),
		RabbitURL:       strings.TrimSpace(os.Getenv("QHPRO_RABBITMQ_URL")),
		Exchange:        strings.TrimSpace(os.Getenv("QHPRO_BUSINESS_EVENTS_EXCHANGE")),
		ShutdownTimeout: durationEnv("NOTIFICATION_SHUTDOWN_TIMEOUT", 10*time.Second),
	}
	if cfg.Exchange == "" {
		cfg.Exchange = defaultBusinessEventsExchange
	}
	if cfg.DatabaseDSN == "" {
		return EventingConfig{}, fmt.Errorf("notification database DSN is required")
	}
	if cfg.RabbitURL == "" {
		return EventingConfig{}, fmt.Errorf("QHPRO_RABBITMQ_URL is required")
	}
	return cfg, nil
}

type DeliveryConfig struct {
	DatabaseDSN    string
	UserGRPCTarget string
	WorkerID       string
	PollInterval   time.Duration
	Lease          time.Duration
}

func LoadDeliveryConfig() (DeliveryConfig, error) {
	if err := LoadConfig(); err != nil {
		return DeliveryConfig{}, err
	}
	cfg := DeliveryConfig{
		DatabaseDSN:    databaseDSN(),
		UserGRPCTarget: firstNonBlank(strings.TrimSpace(os.Getenv("QHPRO_USER_GRPC_ADDR")), strings.TrimSpace(viper.GetString("rpc.user.address"))),
		WorkerID:       firstNonBlank(strings.TrimSpace(os.Getenv("NOTIFICATION_DELIVERY_WORKER_ID")), "notification-delivery"),
		PollInterval:   durationEnv("NOTIFICATION_DELIVERY_POLL_INTERVAL", time.Second),
		Lease:          durationEnv("NOTIFICATION_DELIVERY_LEASE", 30*time.Second),
	}
	if cfg.DatabaseDSN == "" {
		return DeliveryConfig{}, fmt.Errorf("notification database DSN is required")
	}
	if cfg.UserGRPCTarget == "" {
		return DeliveryConfig{}, fmt.Errorf("notification User gRPC target is required")
	}
	return cfg, nil
}

func databaseDSN() string {
	if dsn := strings.TrimSpace(os.Getenv("NOTIFICATION_DATABASE_DSN")); dsn != "" {
		return dsn
	}
	if strings.TrimSpace(AppProperties.Database.Host) == "" || AppProperties.Database.Port <= 0 || strings.TrimSpace(AppProperties.Database.Name) == "" {
		return ""
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		AppProperties.Database.Host,
		AppProperties.Database.Port,
		AppProperties.Database.User,
		AppProperties.Database.Password,
		AppProperties.Database.Name,
	)
}

func firstNonBlank(values ...string) string {
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}

func durationEnv(key string, fallback time.Duration) time.Duration {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	value, err := time.ParseDuration(raw)
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}
