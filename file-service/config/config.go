package config

import (
	"bytes"
	_ "embed"
	"file/env"
	"log"
	"os"

	"github.com/spf13/viper"
)

//go:embed config.yml
var EmbeddedConfig []byte

// LoadConfig đọc cấu hình từ file config.yml
func LoadConfig() {
	env.ENV_RUNTIME = os.Getenv("ENV_RUNTIME")
	viper.SetConfigType("yaml") // Loại file config

	// Đọc biến môi trường (nếu có)
	viper.AutomaticEnv()

	var err error
	if env.ENV_RUNTIME == "release" {
		config_release, err := os.ReadFile("config/product.yml")
		if err != nil {
			log.Fatal(err)
		}
		err = viper.ReadConfig(bytes.NewReader(config_release))
		if err != nil {
			log.Fatalf("Lỗi khi đọc file config: %v", err)
		}
	} else {
		// Đọc file config
		err = viper.ReadConfig(bytes.NewReader(EmbeddedConfig))
	}

	if err != nil {
		log.Fatalf("Lỗi khi đọc file config: %v", err)
	}
}
