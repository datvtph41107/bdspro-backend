// shared/pkg/config/base.go
// Lần đầu tiên: định nghĩa những gì CHUNG cho tất cả service
package config

import "time"

type ReloadStrategy string

const (
	ReloadHot  ReloadStrategy = "hot"
	ReloadCold ReloadStrategy = "cold"
)

// BaseConfig: mọi service đều cần những thứ này.
// Đây là "ngôn ngữ chung" mà Developer + Operator thỏa thuận.
//
// Tại sao là struct lồng nhau thay vì flat?
// Vì struct lồng nhau map trực tiếp vào env var naming convention:
//
//	Server.Port   → APP_SERVER_PORT
//	Database.Host → APP_DATABASE_HOST
//
// Flat struct sẽ dẫn đến APP_PORT, APP_HOST — mơ hồ khi có nhiều "port"
type BaseConfig struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Feature  FeatureConfig  `mapstructure:"feature"  reload:"hot"`
}

type AppConfig struct {
	Name        string `mapstructure:"name"        validate:"required"`
	Environment string `mapstructure:"environment" validate:"required,oneof=development staging production"`
	Version     string `mapstructure:"version"`
}

type ServerConfig struct {
	// Port luôn có kiểu int, không phải string.
	// os.Getenv trả về string → phải convert → có thể lỗi.
	// Dùng struct + viper xử lý conversion tự động, type-safe.
	Port            int           `mapstructure:"port"              validate:"required,min=1024,max=65535"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"      validate:"min=1s,max=60s"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"     validate:"min=1s,max=60s"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"  validate:"min=1s,max=120s"`
}

type DatabaseConfig struct {
	// DSN là một string duy nhất thay vì host/port/user/pass riêng lẻ.
	// Lý do: DSN là unit nguyên tử — thay đổi host đồng thời với port/ssl
	// là chuyện bình thường khi chuyển môi trường.
	// Tách ra thành 4 field → Operator phải update 4 env vars thay vì 1.
	DSN             string        `mapstructure:"dsn"               validate:"required" sensitive:"true"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"    validate:"min=1,max=200"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"    validate:"min=1"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" validate:"min=30s"`
}

type LoggerConfig struct {
	Level  string `mapstructure:"level"  validate:"oneof=debug info warn error" reload:"hot"`
	Format string `mapstructure:"format" validate:"oneof=json text" reload:"cold"`
}

type FeatureConfig struct {
	EnableNewFlow   bool `mapstructure:"enable_new_flow"   reload:"hot"`
	MaintenanceMode bool `mapstructure:"maintenance_mode"  reload:"hot"`
}
