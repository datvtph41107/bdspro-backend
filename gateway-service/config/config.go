package config

import (
	"bytes"
	_ "embed"
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

//go:embed config.yml
var configFile []byte

// LoadConfig đọc cấu hình từ file config.yml
func LoadConfig() {
	// viper.SetConfigName("config")   // Tên file không có phần mở rộng
	viper.SetConfigType("yaml") // Loại file config
	// viper.AddConfigPath("./config") // Đọc từ thư mục gốc

	// Đọc biến môi trường (nếu có)
	viper.AutomaticEnv()

	// Đọc file config
	err := viper.ReadConfig(bytes.NewReader(configFile))
	if err != nil {
		log.Fatalf("Lỗi khi đọc file config: %v", err)
	}

	// Gán giá trị vào struct
	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		log.Fatalf("Lỗi khi parse file config: %v", err)
	}
}
