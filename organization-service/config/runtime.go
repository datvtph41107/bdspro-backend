package config

import (
	"bytes"
	_ "embed"
	"log"
	"organization/env"
	"organization/infrastructure/server/grpc"
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

	viper.SetConfigType("yaml") // 👈 cần thiết!
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

	if chatServiceHost := cfg.GetString("GRPC.CHAT_SERVICE_HOST"); chatServiceHost != "" {
		env.CHAT_SERVICE_HOST = chatServiceHost
		configLogger.Infof("Chat service host specified: %s", chatServiceHost)
	} else {
		configLogger.Warnf("Chat service host not specified, using default: %s", env.CHAT_SERVICE_HOST)
	}
	if chatServicePort := cfg.GetString("GRPC.CHAT_SERVICE_PORT"); chatServicePort != "" {
		env.CHAT_SERVICE_PORT = chatServicePort
		configLogger.Infof("Chat service port specified: %s", chatServicePort)
	} else {
		configLogger.Warnf("Chat service port not specified, using default: %s", env.CHAT_SERVICE_PORT)
	}

	if workerQueueSize := cfg.GetInt("Queue.WORKER_QUEUE_SIZE"); workerQueueSize != 0 {
		env.LOG_WORKER_QUEUE_SIZE = workerQueueSize
		configLogger.Infof("Worker queue size specified: %d", workerQueueSize)
	} else {
		configLogger.Warnf("Worker queue size not specified, using default: %d", env.LOG_WORKER_QUEUE_SIZE)
	}

	if userServiceHost := cfg.GetString("GRPC.USER_SERVICE_HOST"); userServiceHost != "" {
		env.USER_SERVICE_HOST = userServiceHost
		configLogger.Infof("User service host specified: %s", userServiceHost)
	} else {
		configLogger.Warnf("User service host not specified, using default: %s", env.USER_SERVICE_HOST)
	}

	if userServicePort := cfg.GetString("GRPC.USER_SERVICE_PORT"); userServicePort != "" {
		env.USER_SERVICE_PORT = userServicePort
		configLogger.Infof("User service port specified: %s", userServicePort)
	} else {
		configLogger.Warnf("User service port not specified, using default: %s", env.USER_SERVICE_PORT)
	}

	if envRuntime := cfg.GetString("ENV_RUNTIME"); envRuntime != "" {
		env.ENV_RUNTIME = envRuntime
		configLogger.Infof("Environment runtime specified: %s", envRuntime)
	} else {
		configLogger.Warnf("Environment runtime not specified, using default: %s", env.ENV_RUNTIME)
	}

	return nil
}
