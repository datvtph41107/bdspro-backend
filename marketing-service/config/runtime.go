package config

import (
	"bytes"
	_ "embed"
	"log"
	"marketing/env"
	"marketing/infrastructure/server/grpc"
	"os"
	"path"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/spf13/viper"
)

type ConfigApp struct {
	GRPCServer *grpc.GRPCServer
}

//go:embed config.yaml
var ConfigFile []byte

func NewConfigApp(grpcServer *grpc.GRPCServer) *ConfigApp {
	return &ConfigApp{
		GRPCServer: grpcServer,
	}
}

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

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	// Đọc file config
	err := viper.ReadConfig(bytes.NewReader(ConfigFile))
	if err != nil {
		log.Fatalf("Lỗi khi đọc file config: %v", err)
	}
}

func readConfig() error {
	dbType := cfg.GetString("Database.PRIMARY_DB_TYPE")
	env.PRIMARY_DB_TYPE = env.DBType(dbType)
	if env.PRIMARY_DB_TYPE != env.DBTypePostgres {
		configLogger.Warnf("Unsupported database type: %s, using default: postgres", dbType)
		env.PRIMARY_DB_TYPE = env.DBTypePostgres
	}

	if dns := cfg.GetString("Database.DB_DNS"); dns != "" {
		env.DB_DNS = dns
	} else {
		configLogger.Warnf("Database DNS not specified, using default: %s", env.DefaultDBDNS)
	}

	if port := cfg.GetString("GRPC.PORT"); port != "" {
		env.GRPC_PORT = port
	} else {
		configLogger.Warnf("GRPC port not specified, using default: %s", env.DefaultGRPCPort)
	}

	if paymentServiceHost := cfg.GetString("Payment.SERVICE_HOST"); paymentServiceHost != "" {
		env.PAYMENT_SERVICE_HOST = paymentServiceHost
		configLogger.Infof("Payment service host specified: %s", paymentServiceHost)
	} else {
		configLogger.Warnf("Payment service host not specified, using default: %s", env.PAYMENT_SERVICE_HOST)
	}

	if paymentServicePort := cfg.GetString("Payment.SERVICE_PORT"); paymentServicePort != "" {
		env.PAYMENT_SERVICE_PORT = paymentServicePort
		configLogger.Infof("Payment service port specified: %s", paymentServicePort)
	} else {
		configLogger.Warnf("Payment service port not specified, using default: %s", env.PAYMENT_SERVICE_PORT)
	}

	if notificationServiceHost := cfg.GetString("Notification.SERVICE_HOST"); notificationServiceHost != "" {
		env.NOTIFICATION_SERVICE_HOST = notificationServiceHost
		configLogger.Infof("Notification service host specified: %s", notificationServiceHost)
	} else {
		configLogger.Warnf("Notification service host not specified, using default: %s", env.NOTIFICATION_SERVICE_HOST)
	}

	if notificationServicePort := cfg.GetString("Notification.SERVICE_PORT"); notificationServicePort != "" {
		env.NOTIFICATION_SERVICE_PORT = notificationServicePort
		configLogger.Infof("Notification service port specified: %s", notificationServicePort)
	} else {
		configLogger.Warnf("Notification service port not specified, using default: %s", env.NOTIFICATION_SERVICE_PORT)
	}

	if envRuntime := cfg.GetString("ENV_RUNTIME"); envRuntime != "" {
		env.ENV_RUNTIME = envRuntime
		configLogger.Infof("Environment runtime specified: %s", envRuntime)
	} else {
		configLogger.Warnf("Environment runtime not specified, using default: %s", env.ENV_RUNTIME)
	}

	return nil
} 