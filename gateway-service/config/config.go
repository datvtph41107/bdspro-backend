package config

import (
	"bytes"
	_ "embed"
	"gateway/env"
	"log"

	"github.com/spf13/viper"
)

//go:embed config.yml
var ConfigFile []byte

// LoadConfig đọc cấu hình từ file config.yml
func LoadConfig() {
	// viper.SetConfigName("config")   // Tên file không có phần mở rộng
	viper.SetConfigType("yaml") // Loại file config
	// viper.AddConfigPath("./config") // Đọc từ thư mục gốc

	// Đọc biến môi trường (nếu có)
	viper.AutomaticEnv()

	// Đọc file config
	if env.ENV_RUNTIME == "release" {
		err := viper.ReadConfig(bytes.NewReader(ConfigFile))
		if err != nil {
			log.Fatalf("Lỗi khi đọc file config: %v", err)
		}
	} else {
		err := viper.ReadConfig(bytes.NewReader(ConfigFile))
		if err != nil {
			log.Fatalf("Lỗi khi đọc file config: %v", err)
		}
	}
}
