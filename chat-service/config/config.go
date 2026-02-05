package config

import (
	"chat/env"
	"common"
	"log"
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
		Password string `mapstructure:"password"`
	} `mapstructure:"redis"`

	Swagger struct {
		Path string `mapstructure:"path"`
	} `mapstructure:"swagger"`

	Worker struct {
		Number int `mapstructure:"number"`
	} `mapstructure:"worker"`

	MaxParticipant       int `mapstructure:"max_participant"`
	RequiredParticipants int `mapstructure:"required_participants"`

	Services struct {
		User struct {
			Host string `mapstructure:"host"`
			Port string `mapstructure:"port"`
		} `mapstructure:"user"`
		Bdspro struct {
			Host string `mapstructure:"host"`
			Port string `mapstructure:"port"`
		} `mapstructure:"bdspro"`
	} `mapstructure:"services"`
}

var AppProperties ChatProperties

// LoadConfig đọc cấu hình từ file config.yml
func LoadConfig() {
	env.ENV_RUNTIME = os.Getenv("ENV_RUNTIME")
	common.LoadYMLConfigFile(env.ENV_RUNTIME)

	if err := viper.Unmarshal(&AppProperties); err != nil {
		log.Fatalf("Lỗi parse config: %v", err)
	}
}

// InitConfig - alias cho LoadConfig để maintain backward compatibility
func InitConfig() {
	LoadConfig()
}
