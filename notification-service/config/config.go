package config

import (
	"common/configloader"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	C_SESSION_ID                  = "C_SESSION_ID"
	defaultBusinessEventsExchange = "qhpro.payment.events"
	defaultDeliveryWorkerID       = "notification-delivery"
)

// UserProperties — process config từ config/<profile>.yml (local/develop/…).
type UserProperties struct {
	Server struct {
		TCPPort int `mapstructure:"tcp_port"`
	} `mapstructure:"server"`
	Database struct {
		DSN      string `mapstructure:"dsn"`
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
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`
	RPC struct {
		User struct {
			Address string `mapstructure:"address"`
		} `mapstructure:"user"`
	} `mapstructure:"rpc"`
	Rabbit struct {
		URL      string `mapstructure:"url"`
		Exchange string `mapstructure:"exchange"`
	} `mapstructure:"rabbit"`
	Delivery struct {
		WorkerID     string `mapstructure:"worker_id"`
		PollInterval string `mapstructure:"poll_interval"`
		Lease        string `mapstructure:"lease"`
	} `mapstructure:"delivery"`
	ShutdownTimeout string `mapstructure:"shutdown_timeout"`
}

var AppProperties UserProperties

// LoadConfig đọc đúng một file config/<profile>.yml vào Viper + AppProperties.
func LoadConfig() error {
	if _, err := configloader.LoadRuntimeYML(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&AppProperties); err != nil {
		return fmt.Errorf("unmarshal notification config: %w", err)
	}
	return nil
}

// RuntimeConfig — wiring kỹ thuật của gRPC process (không đọc os.Getenv).
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
		RedisAddress:           strings.TrimSpace(AppProperties.Redis.Host),
		RedisPassword:          AppProperties.Redis.Password,
		RedisDB:                AppProperties.Redis.DB,
		UserGRPCTarget:         strings.TrimSpace(AppProperties.RPC.User.Address),
		RabbitURL:              strings.TrimSpace(AppProperties.Rabbit.URL),
		BusinessEventsExchange: firstNonBlank(strings.TrimSpace(AppProperties.Rabbit.Exchange), defaultBusinessEventsExchange),
		DeliveryWorkerID:       firstNonBlank(strings.TrimSpace(AppProperties.Delivery.WorkerID), defaultDeliveryWorkerID),
		DeliveryPollInterval:   parseDurationOr(AppProperties.Delivery.PollInterval, time.Second),
		DeliveryLease:          parseDurationOr(AppProperties.Delivery.Lease, 30*time.Second),
		ShutdownTimeout:        parseDurationOr(AppProperties.ShutdownTimeout, 10*time.Second),
	}
	if err := validateRedisDB(cfg.RedisDB); err != nil {
		return RuntimeConfig{}, err
	}
	if strings.TrimSpace(cfg.DatabaseDSN) == "" {
		return RuntimeConfig{}, fmt.Errorf("notification database.dsn is required in config YAML")
	}
	if strings.TrimSpace(cfg.RedisAddress) == "" {
		return RuntimeConfig{}, fmt.Errorf("notification redis.host is required in config YAML")
	}
	if strings.TrimSpace(cfg.UserGRPCTarget) == "" {
		return RuntimeConfig{}, fmt.Errorf("notification rpc.user.address is required in config YAML")
	}
	if cfg.ServerPort <= 0 || cfg.ServerPort > 65535 {
		return RuntimeConfig{}, fmt.Errorf("notification server.tcp_port is required in config YAML")
	}
	return cfg, nil
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
		RabbitURL:       strings.TrimSpace(AppProperties.Rabbit.URL),
		Exchange:        firstNonBlank(strings.TrimSpace(AppProperties.Rabbit.Exchange), defaultBusinessEventsExchange),
		ShutdownTimeout: parseDurationOr(AppProperties.ShutdownTimeout, 10*time.Second),
	}
	if cfg.DatabaseDSN == "" {
		return EventingConfig{}, fmt.Errorf("notification database.dsn is required in config YAML")
	}
	if cfg.RabbitURL == "" {
		return EventingConfig{}, fmt.Errorf("notification rabbit.url is required in config YAML")
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
		UserGRPCTarget: strings.TrimSpace(AppProperties.RPC.User.Address),
		WorkerID:       firstNonBlank(strings.TrimSpace(AppProperties.Delivery.WorkerID), defaultDeliveryWorkerID),
		PollInterval:   parseDurationOr(AppProperties.Delivery.PollInterval, time.Second),
		Lease:          parseDurationOr(AppProperties.Delivery.Lease, 30*time.Second),
	}
	if cfg.DatabaseDSN == "" {
		return DeliveryConfig{}, fmt.Errorf("notification database.dsn is required in config YAML")
	}
	if cfg.UserGRPCTarget == "" {
		return DeliveryConfig{}, fmt.Errorf("notification rpc.user.address is required in config YAML")
	}
	return cfg, nil
}

func databaseDSN() string {
	if dsn := strings.TrimSpace(AppProperties.Database.DSN); dsn != "" {
		return dsn
	}
	if dsn := strings.TrimSpace(viper.GetString("database.dsn")); dsn != "" {
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

func parseDurationOr(raw string, fallback time.Duration) time.Duration {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fallback
	}
	value, err := time.ParseDuration(raw)
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}

func validateRedisDB(db int) error {
	if db < 0 || db > 15 {
		return fmt.Errorf("invalid redis.db %d: must be between 0 and 15", db)
	}
	return nil
}
