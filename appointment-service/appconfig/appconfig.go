package appconfig

import (
	"os"
	"path"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/spf13/viper"
)

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

var (
	cfg          *viper.Viper
	configLogger = flogging.MustGetLogger("config")
)

func InitConfig() {
	loadConfig()
	if err := readConfig(); err != nil {
		configLogger.Fatalf("Error reading config values: %v", err)
	}
}

func loadConfig() {
	file := "config.yaml"
	if envFile := os.Getenv("CONFIG_FILE"); envFile != "" {
		file = envFile
	}
	flogging.ActivateSpec("debug")

	cfg = viper.New()
	cfg.SetConfigType("yaml")
	cfg.SetConfigName(path.Base(file))
	cfg.AddConfigPath(path.Dir(file))
	cfg.AddConfigPath("./config/")
	cfg.AddConfigPath("../config/")

	configLogger.Infof("Loading config file: %s", file)

	if err := cfg.ReadInConfig(); err != nil {
		configLogger.Fatalf("Error parsing config file: %s", err)
	}
}

func readConfig() error {
	if port := cfg.GetString("GRPC.PORT"); port != "" {
		GRPC_PORT = port
		configLogger.Infof("GRPC port specified: %s", GRPC_PORT)
	} else {
		configLogger.Warnf("GRPC port not specified, using default: %s", GRPC_PORT)
	}

	if dns := cfg.GetString("Database.DB_DNS"); dns != "" {
		DB_DNS = dns
		configLogger.Infof("Database DNS specified: %s", DB_DNS)
	} else {
		configLogger.Warnf("Database DNS not specified, using default: %s", DB_DNS)
	}

	if emailHost := cfg.GetString("Email.HOST"); emailHost != "" {
		EMAIL_HOST = emailHost
		configLogger.Infof("Email host specified: %s", EMAIL_HOST)
	} else {
		configLogger.Warnf("Email host not specified, using default: %s", EMAIL_HOST)
	}

	if emailPort := cfg.GetInt("Email.PORT"); emailPort != 0 {
		EMAIL_PORT = emailPort
		configLogger.Infof("Email port specified: %d", EMAIL_PORT)
	} else {
		configLogger.Warnf("Email port not specified, using default: %d", EMAIL_PORT)
	}

	if emailUsername := cfg.GetString("Email.USERNAME"); emailUsername != "" {
		EMAIL_USERNAME = emailUsername
		configLogger.Infof("Email username specified: %s", EMAIL_USERNAME)
	} else {
		configLogger.Warnf("Email username not specified, using default: %s", EMAIL_USERNAME)
	}

	if emailPassword := cfg.GetString("Email.PASSWORD"); emailPassword != "" {
		EMAIL_PASSWORD = emailPassword
		configLogger.Infof("Email password specified: %s", EMAIL_PASSWORD)
	} else {
		configLogger.Warnf("Email password not specified, using default: %s", EMAIL_PASSWORD)
	}

	if emailFrom := cfg.GetString("Email.FROM"); emailFrom != "" {
		EMAIL_FROM = emailFrom
		configLogger.Infof("Email from specified: %s", EMAIL_FROM)
	} else {
		configLogger.Warnf("Email from not specified, using default: %s", EMAIL_FROM)
	}

	if emailName := cfg.GetString("Email.NAME"); emailName != "" {
		EMAIL_NAME = emailName
		configLogger.Infof("Email name specified: %s", EMAIL_NAME)
	} else {
		configLogger.Warnf("Email name not specified, using default: %s", EMAIL_NAME)
	}

	return nil
}
