package c_config

import (
	"log"

	"github.com/spf13/viper"
)

// Config struct chứa toàn bộ thông tin từ file config.yml
type Config struct {
	Server struct {
		Port string
	}
	Eureka struct {
		Url  string
		Port string
	}
	Database struct {
		Host     string
		User     string
		Password string
	}
}

// Cấu hình toàn cục
var AppConfig Config

// LoadConfig đọc cấu hình từ file config.yml
func LoadConfig() {
	// Gán giá trị vào struct
	err := viper.Unmarshal(&AppConfig)
	if err != nil {
		log.Fatalf("Lỗi khi parse file config: %v", err)
	}
}
