package appconfig

import (
	"log"
	"os"
	"path"

	"github.com/spf13/viper"
)

var (
	cfg                       *viper.Viper
	Host                      string
	Port                      string
	DB_URL                    string
	Whitelist                 []string
	SepayAPIKey               string
	NOTIFICATION_SERVICE_HOST string
	NOTIFICATION_SERVICE_PORT string
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
	DB_URL = os.Getenv("PAYMENT_DATABASE_DSN")
	if DB_URL == "" {
		DB_URL = cfg.GetString("database.url")
	}
	if Whitelist = cfg.GetStringSlice("whitelist"); Whitelist == nil {
		Whitelist = []string{}
	}
	SepayAPIKey = os.Getenv("SEPAY_API_KEY")
	if SepayAPIKey == "" {
		SepayAPIKey = cfg.GetString("sepay.api_key")
	}
	if NOTIFICATION_SERVICE_HOST = cfg.GetString("notification.host"); NOTIFICATION_SERVICE_HOST == "" {
		NOTIFICATION_SERVICE_HOST = "localhost"
	}
	if NOTIFICATION_SERVICE_PORT = cfg.GetString("notification.port"); NOTIFICATION_SERVICE_PORT == "" {
		NOTIFICATION_SERVICE_PORT = "8204"
	}
}
