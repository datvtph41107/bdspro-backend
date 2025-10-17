package config

import (
	"common"
	"os"
	"user/env"
)

// LoadProperties đọc cấu hình từ file config.yml
func LoadProperties() {
	env.ENV_RUNTIME = os.Getenv("ENV_RUNTIME")
	common.LoadYMLConfigFile(env.ENV_RUNTIME)
	// viper.SetConfigType("yaml") // Loại file config

	// // Đọc biến môi trường (nếu có)
	// viper.AutomaticEnv()

	// var err error
	// if env.ENV_RUNTIME == "release" {
	// 	config_release, err := os.ReadFile("config/product.yml")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	err = viper.ReadConfig(bytes.NewReader(config_release))
	// 	if err != nil {
	// 		log.Fatalf("Lỗi khi đọc file config: %v", err)
	// 	}
	// } else {
	// 	// Đọc file config
	// 	err = viper.ReadConfig(bytes.NewReader(EmbeddedConfig))
	// }

	// if err != nil {
	// 	log.Fatalf("Lỗi khi đọc file config: %v", err)
	// }

	// if err := viper.Unmarshal(&property.AppProperties); err != nil {
	// 	log.Fatalf("Lỗi parse config: %v", err)
	// }
}
