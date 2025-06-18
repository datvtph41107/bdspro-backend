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
	DB_URL          string
)

func InitConfig() {
	loadConfig()
	ReadConfig()
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

	if Host = cfg.GetString("host"); Host == "" {
		Host = "0.0.0.0"
	}
	if Port = cfg.GetString("port"); Port == "" {
		Port = "50055"
	}
	if DB_URL = cfg.GetString("database.url"); DB_URL == "" {
		DB_URL = "host=localhost user=postgres password=postgres dbname=payment port=5432 sslmode=disable"
	}
}
