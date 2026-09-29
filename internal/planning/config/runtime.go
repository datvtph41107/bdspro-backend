package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// RuntimeConfig is the process-owned technical configuration for the TQD gRPC
// process. Business/application packages receive typed values and capabilities;
// they never read process environment directly.
type RuntimeConfig struct {
	Properties *Properties

	DatabaseDSN string
	Redis       RedisRuntimeConfig
	RPC         RPCRuntimeConfig
	File        FileRuntimeConfig
	Classify    ClassifyRuntimeConfig
	Report      ReportRuntimeConfig
	Shutdown    time.Duration
}

type RedisRuntimeConfig struct {
	Address  string
	Password string
	DB       int
}

type RPCRuntimeConfig struct {
	User      string
	Auth      string
	Assistant string
}

type FileRuntimeConfig struct {
	ServiceName    string
	BaseURL        string
	PublicBaseURL  string
	ServiceAuthKey string
	Timeout        time.Duration
}

type ClassifyRuntimeConfig struct {
	Enabled        bool
	Interval       time.Duration
	BatchSize      int
	ReadContent    bool
	MaxInlineBytes int
}

type ReportRuntimeConfig struct {
	GeneratorMode    string
	WeasyPrintBinary string
	RenderTimeout    time.Duration
	File             FileRuntimeConfig
	WorkerID         string
}

func LoadRuntimeConfig() (*RuntimeConfig, error) {
	props, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	grpcPort, err := parsePortOverride(os.Getenv("TQD_GRPC_PORT"), props.Server.GrpcPort)
	if err != nil {
		return nil, err
	}
	props.Server.GrpcPort = grpcPort
	redisRuntime, err := loadRedisRuntimeConfig()
	if err != nil {
		return nil, err
	}

	cfg := &RuntimeConfig{
		Properties:  props,
		DatabaseDSN: firstNonEmpty(os.Getenv("TQD_DATABASE_URL"), props.Database.DSN),
		Redis:       redisRuntime,
		RPC: RPCRuntimeConfig{
			User:      firstNonEmpty(os.Getenv("QHPRO_USER_GRPC_ADDR"), viper.GetString("rpc.user.address")),
			Auth:      firstNonEmpty(os.Getenv("QHPRO_AUTH_GRPC_ADDR"), viper.GetString("rpc.auth.address")),
			Assistant: firstNonEmpty(os.Getenv("QHPRO_ASSISTANT_GRPC_ADDR"), assistantDefaultTarget()),
		},
		File: FileRuntimeConfig{
			ServiceName:    firstNonEmpty(os.Getenv("QHPRO_SERVICE_ID"), "tqd-service"),
			BaseURL:        firstNonEmpty(os.Getenv("QHPRO_FILE_HTTP_BASE_URL"), viper.GetString("file_service.base_url"), "http://127.0.0.1:8002"),
			PublicBaseURL:  firstNonEmpty(os.Getenv("QHPRO_FILE_PUBLIC_BASE_URL"), viper.GetString("file_service.public_base_url")),
			ServiceAuthKey: firstNonEmpty(os.Getenv("QHPRO_FILE_HTTP_SERVICE_AUTH_KEY"), os.Getenv("SERVICE_AUTH_KEY"), os.Getenv("FILE_SERVICE_AUTH_KEY"), viper.GetString("file_service.service_auth_key")),
			Timeout:        durationEnv("QHPRO_FILE_HTTP_TIMEOUT", durationSeconds(viper.GetInt("file_service.timeout_seconds"), 30*time.Second)),
		},
		Classify: ClassifyRuntimeConfig{
			Enabled:        boolSetting("TQD_CLASSIFY_ENABLED", "qh_planning_classify.enabled", true),
			Interval:       durationEnv("TQD_CLASSIFY_INTERVAL", durationSeconds(viper.GetInt("qh_planning_classify.interval_sec"), 15*time.Second)),
			BatchSize:      intSetting("TQD_CLASSIFY_BATCH_SIZE", "qh_planning_classify.batch_size", 5),
			ReadContent:    boolSetting("TQD_CLASSIFY_READ_CONTENT", "qh_planning_classify.read_content", true),
			MaxInlineBytes: intSetting("TQD_CLASSIFY_MAX_INLINE_BYTES", "qh_planning_classify.max_inline_bytes", 15<<20),
		},
		Shutdown: durationEnv("TQD_SHUTDOWN_TIMEOUT", 15*time.Second),
	}
	if cfg.File.PublicBaseURL == "" {
		cfg.File.PublicBaseURL = cfg.File.BaseURL
	}
	cfg.Report = ReportRuntimeConfig{
		GeneratorMode:    strings.ToLower(firstNonEmpty(os.Getenv("QHPRO_REPORT_GENERATOR_MODE"), "disabled")),
		WeasyPrintBinary: firstNonEmpty(os.Getenv("QHPRO_REPORT_WEASYPRINT_BINARY"), "weasyprint"),
		RenderTimeout:    durationEnv("QHPRO_REPORT_RENDER_TIMEOUT", 45*time.Second),
		File:             cfg.File,
		WorkerID:         firstNonEmpty(os.Getenv("QHPRO_REPORT_WORKER_ID"), "tqd-report-worker"),
	}

	if strings.TrimSpace(cfg.DatabaseDSN) == "" {
		return nil, fmt.Errorf("TQD database DSN is required")
	}
	if strings.TrimSpace(cfg.Redis.Address) == "" {
		return nil, fmt.Errorf("TQD Redis address is required")
	}
	if strings.TrimSpace(cfg.RPC.User) == "" {
		return nil, fmt.Errorf("TQD User gRPC target is required")
	}
	if strings.TrimSpace(cfg.RPC.Auth) == "" {
		return nil, fmt.Errorf("TQD Auth gRPC target is required")
	}
	if cfg.Classify.Enabled && strings.TrimSpace(cfg.RPC.Assistant) == "" {
		return nil, fmt.Errorf("TQD Assistant gRPC target is required while classify worker is enabled")
	}
	if cfg.Report.GeneratorMode != "disabled" && cfg.Report.GeneratorMode != "none" && strings.TrimSpace(cfg.File.ServiceAuthKey) == "" {
		return nil, fmt.Errorf("QHPRO_FILE_HTTP_SERVICE_AUTH_KEY is required when report generation is enabled")
	}
	return cfg, nil
}

func LoadRedisRuntimeConfig() (RedisRuntimeConfig, error) {
	if _, err := LoadConfig(); err != nil {
		return RedisRuntimeConfig{}, err
	}
	return loadRedisRuntimeConfig()
}

func loadRedisRuntimeConfig() (RedisRuntimeConfig, error) {
	redisDB, err := parseRedisDBOverride(os.Getenv("TQD_REDIS_DB"), viper.GetInt("redis.db"))
	if err != nil {
		return RedisRuntimeConfig{}, err
	}
	cfg := RedisRuntimeConfig{
		Address:  firstNonEmpty(os.Getenv("TQD_REDIS_ADDR"), viper.GetString("redis.host")),
		Password: firstNonEmpty(os.Getenv("TQD_REDIS_PASSWORD"), viper.GetString("redis.pass")),
		DB:       redisDB,
	}
	if strings.TrimSpace(cfg.Address) == "" {
		return RedisRuntimeConfig{}, fmt.Errorf("TQD Redis address is required")
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
		return 0, fmt.Errorf("invalid TQD_REDIS_DB %q: must be between 0 and 15", raw)
	}
	return db, nil
}

func parsePortOverride(raw string, fallback int) (int, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fallback, nil
	}
	port, err := strconv.Atoi(raw)
	if err != nil || port <= 0 || port > 65535 {
		return 0, fmt.Errorf("invalid TQD_GRPC_PORT %q", raw)
	}
	return port, nil
}

func assistantDefaultTarget() string {
	return "127.0.0.1:8218"
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}

func durationSeconds(seconds int, fallback time.Duration) time.Duration {
	if seconds <= 0 {
		return fallback
	}
	return time.Duration(seconds) * time.Second
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

func intSetting(envKey, viperKey string, fallback int) int {
	if raw := strings.TrimSpace(os.Getenv(envKey)); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil && n > 0 {
			return n
		}
	}
	if n := viper.GetInt(viperKey); n > 0 {
		return n
	}
	return fallback
}

func boolSetting(envKey, viperKey string, fallback bool) bool {
	if raw := strings.TrimSpace(os.Getenv(envKey)); raw != "" {
		if value, err := strconv.ParseBool(raw); err == nil {
			return value
		}
	}
	if viper.IsSet(viperKey) {
		return viper.GetBool(viperKey)
	}
	return fallback
}
