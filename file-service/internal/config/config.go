package internal_config

import (
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	Storage   StorageConfig
	JWT       JWTConfig
	Security  SecurityConfig
	Queue     QueueConfig
	RateLimit RateLimitConfig
}

type RateLimitConfig struct {
	RequestsPerSecond float64
	Burst             int
	PerUser           bool
	Redis             bool
}

type ServerConfig struct {
	Port         string
	Mode         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	PoolSize int
	Enabled  bool
}

type StorageConfig struct {
	Type             string // "local", "s3", "gcs"
	LocalPath        string
	S3Bucket         string
	S3Region         string
	S3Endpoint       string
	MaxUploadSize    int64
	AllowedMimeTypes []string
	CDNBaseURL       string
}

type JWTConfig struct {
	SecretKey          string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	ServiceTokenExpiry time.Duration
}

type SecurityConfig struct {
	XorCryptKey     string
	SignatureKey    string
	SignatureExpiry time.Duration
	ServiceAuthKey  string
}

type QueueConfig struct {
	WorkerCount int
	QueueSize   int
	RedisQueue  bool
}

var GlobalConfig *Config

func LoadConfig() error {
	// Load .env file
	// Setup viper
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	GlobalConfig = &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8002"),
			Mode:         getEnv("SERVER_MODE", "release"),
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:  getEnvAsDuration("SERVER_IDLE_TIMEOUT", 120*time.Second),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", ""),
			DBName:          getEnv("DB_NAME", "file_service"),
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 100),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
			PoolSize: getEnvAsInt("REDIS_POOL_SIZE", 100),
		},
		Storage: StorageConfig{
			Type:          getEnv("STORAGE_TYPE", "local"),
			LocalPath:     getEnv("STORAGE_LOCAL_PATH", "./files"),
			S3Bucket:      getEnv("S3_BUCKET", ""),
			S3Region:      getEnv("S3_REGION", ""),
			S3Endpoint:    getEnv("S3_ENDPOINT", ""),
			MaxUploadSize: getEnvAsInt64("MAX_UPLOAD_SIZE", 100*1024*1024), // 100MB
			AllowedMimeTypes: []string{
				"image/jpeg",
				"image/png",
				"image/gif",
				"video/mp4",
				"video/quicktime",
				"application/pdf",
			},
			CDNBaseURL: getEnv("CDN_BASE_URL", ""),
		},
		JWT: JWTConfig{
			SecretKey:          getEnv("JWT_SECRET_KEY", "your-secret-key"),
			AccessTokenExpiry:  getEnvAsDuration("JWT_ACCESS_EXPIRY", 15*time.Minute),
			RefreshTokenExpiry: getEnvAsDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour),
			ServiceTokenExpiry: getEnvAsDuration("JWT_SERVICE_EXPIRY", 24*time.Hour),
		},
		Security: SecurityConfig{
			XorCryptKey:     getEnv("XOR_CRYPT_KEY", "default-xor-key"),
			SignatureKey:    getEnv("SIGNATURE_KEY", "default-signature-key"),
			SignatureExpiry: getEnvAsDuration("SIGNATURE_EXPIRY", 5*time.Minute),
			ServiceAuthKey:  getEnv("SERVICE_AUTH_KEY", "service-internal-key"),
		},
		Queue: QueueConfig{
			WorkerCount: getEnvAsInt("QUEUE_WORKER_COUNT", 10),
			QueueSize:   getEnvAsInt("QUEUE_SIZE", 100),
			RedisQueue:  getEnvAsBool("REDIS_QUEUE_ENABLED", false),
		},
	}

	return nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}
