package appconfig

import (
	"log"
	"os"
	"path"

	"github.com/spf13/viper"
)

var (
	cfg             *viper.Viper
	Host            string
	Port            string
	UserServiceHost string
	UserServicePort string
)

func InitConfig() {
	loadConfig()
}

func loadConfig() {
	file := "config.yaml"
	if envFile := os.Getenv("CONFIG_FILE"); envFile != "" {
		file = envFile
	}

	cfg = viper.New()
	cfg.SetConfigType("yaml")
	cfg.SetConfigName(path.Base(file))
	cfg.AddConfigPath(path.Dir(file))
	cfg.AddConfigPath("./config/")
	cfg.AddConfigPath("../config/")

	if err := cfg.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
}

func ReadConfig() {
	viper.ReadInConfig()

	if Host = viper.GetString("host"); Host == "" {
		Host = "0.0.0.0"
	}
	if Port = viper.GetString("port"); Port == "" {
		Port = "8080"
	}
	if UserServiceHost = viper.GetString("user_service_host"); UserServiceHost == "" {
		UserServiceHost = "localhost"
	}
	if UserServicePort = viper.GetString("user_service_port"); UserServicePort == "" {
		UserServicePort = "50051"
	}
}
