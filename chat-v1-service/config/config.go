package config

import (
	"common/configloader"
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

const C_SESSION_ID = "C_SESSION_ID"

// ChatProperties chứa cấu hình của ứng dụng
type ChatProperties struct {
	Server struct {
		Port    string `mapstructure:"port"`
		TcpPort string `mapstructure:"tcp_port"`
	} `mapstructure:"server"`

	Eureka struct {
		Url  string `mapstructure:"url"`
		Port string `mapstructure:"port"`
	} `mapstructure:"eureka"`

	Database struct {
		Type string `mapstructure:"type"`
		DSN  string `mapstructure:"dsn"`
	} `mapstructure:"database"`

	Timeout struct {
		Transaction string `mapstructure:"transaction"`
		Idle        string `mapstructure:"idle"`
		Read        string `mapstructure:"read"`
		Write       string `mapstructure:"write"`
	} `mapstructure:"timeout"`

	Redis struct {
		Host     string `mapstructure:"host"`
		Password string `mapstructure:"pass"`
	} `mapstructure:"redis"`

	Swagger struct {
		Path string `mapstructure:"path"`
	} `mapstructure:"swagger"`

	Worker struct {
		Number int `mapstructure:"number"`
	} `mapstructure:"worker"`

	MaxParticipant int `mapstructure:"max_participant"`
}

var AppProperties ChatProperties

// LoadConfig đọc cấu hình từ file config.yml
func LoadConfig() {
	if _, err := configloader.LoadRuntimeYML(); err != nil {
		slog.Error(
			"load Chat V1 runtime config",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	if err := viper.Unmarshal(&AppProperties); err != nil {
		slog.Error(
			"Lỗi parse config",
			slog.Any("error", err),
		)
		os.Exit(1)
	}
}

// InitConfig - alias cho LoadConfig để maintain backward compatibility
func InitConfig() {
	LoadConfig()
}
