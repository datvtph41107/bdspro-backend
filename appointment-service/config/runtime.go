package config

import (
	"bytes"
	_ "embed"
	"log"
	"os"

	"github.com/spf13/viper"
)

//go:embed app.yml
var EmbeddedConfig []byte

const (
	SERVICE_NAME = "appointment-service"
	USER_CONTEXT = "user"
)

var (
	DB_DNS    = "host=localhost user=postgres password=postgres dbname=appointment port=5432 sslmode=disable"
	GRPC_PORT = "8213"
)

var (
	EMAIL_HOST     = "smtp.gmail.com"
	EMAIL_PORT     = 587
	EMAIL_USERNAME = "tuandq2112@gmail.com"
	EMAIL_PASSWORD = "tuandq2112"
	EMAIL_FROM     = "tuandq2112@gmail.com"
	EMAIL_NAME     = "Tuấn Đặng"
)

func InitConfig() {
	LoadConfig()
	if err := readConfig(); err != nil {
		log.Fatalf("Error reading config values: %v", err)
	}
}

// LoadConfig đọc cấu hình từ file config.yml
func LoadConfig() {
	envRuntime := os.Getenv("ENV_RUNTIME")
	viper.SetConfigType("yaml") // Loại file config

	// Đọc biến môi trường (nếu có)
	viper.AutomaticEnv()

	var err error
	if envRuntime == "release" {
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

func readConfig() error {
	if port := viper.GetString("GRPC.PORT"); port != "" {
		GRPC_PORT = port
		log.Printf("GRPC port specified: %s", GRPC_PORT)
	} else {
		log.Printf("GRPC port not specified, using default: %s", GRPC_PORT)
	}

	if dns := viper.GetString("Database.DB_DNS"); dns != "" {
		DB_DNS = dns
		log.Printf("Database DNS specified: %s", DB_DNS)
	} else {
		log.Printf("Database DNS not specified, using default: %s", DB_DNS)
	}

	if emailHost := viper.GetString("Email.HOST"); emailHost != "" {
		EMAIL_HOST = emailHost
		log.Printf("Email host specified: %s", EMAIL_HOST)
	} else {
		log.Printf("Email host not specified, using default: %s", EMAIL_HOST)
	}

	if emailPort := viper.GetInt("Email.PORT"); emailPort != 0 {
		EMAIL_PORT = emailPort
		log.Printf("Email port specified: %d", EMAIL_PORT)
	} else {
		log.Printf("Email port not specified, using default: %d", EMAIL_PORT)
	}

	if emailUsername := viper.GetString("Email.USERNAME"); emailUsername != "" {
		EMAIL_USERNAME = emailUsername
		log.Printf("Email username specified: %s", EMAIL_USERNAME)
	} else {
		log.Printf("Email username not specified, using default: %s", EMAIL_USERNAME)
	}

	if emailPassword := viper.GetString("Email.PASSWORD"); emailPassword != "" {
		EMAIL_PASSWORD = emailPassword
		log.Printf("Email password specified: %s", EMAIL_PASSWORD)
	} else {
		log.Printf("Email password not specified, using default: %s", EMAIL_PASSWORD)
	}

	if emailFrom := viper.GetString("Email.FROM"); emailFrom != "" {
		EMAIL_FROM = emailFrom
		log.Printf("Email from specified: %s", EMAIL_FROM)
	} else {
		log.Printf("Email from not specified, using default: %s", EMAIL_FROM)
	}

	if emailName := viper.GetString("Email.NAME"); emailName != "" {
		EMAIL_NAME = emailName
		log.Printf("Email name specified: %s", EMAIL_NAME)
	} else {
		log.Printf("Email name not specified, using default: %s", EMAIL_NAME)
	}

	return nil
}
