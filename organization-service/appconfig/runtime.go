package appconfig

import (
	"os"
	"path"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/spf13/viper"
)

const (
	SERVICE_NAME = "organization-service"
)

type DBType string

const (
	DBTypePostgres DBType = "postgres"

	DefaultDBDNS    = "host=localhost user=postgres password=postgres dbname=organization port=5432 sslmode=disable"
	DefaultGRPCPort = "50052"
)

var (
	cfg             *viper.Viper
	configLogger    = flogging.MustGetLogger("config")
	PRIMARY_DB_TYPE = DBTypePostgres
	DB_DNS          = DefaultDBDNS
	GRPC_PORT       = DefaultGRPCPort

	LOG_WORKER_QUEUE_SIZE = 1000

	CHAT_SERVICE_HOST = "localhost"
	CHAT_SERVICE_PORT = "50053"

	USER_SERVICE_HOST = "localhost"
	USER_SERVICE_PORT = "8201"

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
	dbType := cfg.GetString("Database.PRIMARY_DB_TYPE")
	PRIMARY_DB_TYPE = DBType(dbType)
	if PRIMARY_DB_TYPE != DBTypePostgres {
		configLogger.Warnf("Unsupported database type: %s, using default: postgres", dbType)
		PRIMARY_DB_TYPE = DBTypePostgres
	}

	if dns := cfg.GetString("Database.DB_DNS"); dns != "" {
		DB_DNS = dns
	} else {
		configLogger.Warnf("Database DNS not specified, using default: %s", DefaultDBDNS)
	}

	if port := cfg.GetString("GRPC.PORT"); port != "" {
		GRPC_PORT = port
	} else {
		configLogger.Warnf("GRPC port not specified, using default: %s", DefaultGRPCPort)
	}

	if chatServiceHost := cfg.GetString("GRPC.CHAT_SERVICE_HOST"); chatServiceHost != "" {
		CHAT_SERVICE_HOST = chatServiceHost
		configLogger.Infof("Chat service host specified: %s", chatServiceHost)
	} else {
		configLogger.Warnf("Chat service host not specified, using default: %s", CHAT_SERVICE_HOST)
	}
	if chatServicePort := cfg.GetString("GRPC.CHAT_SERVICE_PORT"); chatServicePort != "" {
		CHAT_SERVICE_PORT = chatServicePort
		configLogger.Infof("Chat service port specified: %s", chatServicePort)
	} else {
		configLogger.Warnf("Chat service port not specified, using default: %s", CHAT_SERVICE_PORT)
	}

	if workerQueueSize := cfg.GetInt("Queue.WORKER_QUEUE_SIZE"); workerQueueSize != 0 {
		LOG_WORKER_QUEUE_SIZE = workerQueueSize
		configLogger.Infof("Worker queue size specified: %d", workerQueueSize)
	} else {
		configLogger.Warnf("Worker queue size not specified, using default: %d", LOG_WORKER_QUEUE_SIZE)
	}

	
	if userServiceHost := cfg.GetString("GRPC.USER_SERVICE_HOST"); userServiceHost != "" {
		USER_SERVICE_HOST = userServiceHost
		configLogger.Infof("User service host specified: %s", userServiceHost)
	} else {
		configLogger.Warnf("User service host not specified, using default: %s", USER_SERVICE_HOST)
	}


	if userServicePort := cfg.GetString("GRPC.USER_SERVICE_PORT"); userServicePort != "" {
		USER_SERVICE_PORT = userServicePort
		configLogger.Infof("User service port specified: %s", userServicePort)
	} else {
		configLogger.Warnf("User service port not specified, using default: %s", USER_SERVICE_PORT)
	}

	return nil
}
