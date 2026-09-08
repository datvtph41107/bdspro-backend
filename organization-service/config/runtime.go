package config

import (
	"common/configloader"
	"fmt"
	"organization/env"
	"organization/infrastructure/server/grpc"
	"strings"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/spf13/viper"
)

type ConfigApp struct {
	GRPCServer *grpc.GRPCServer
}

func NewConfigApp(grpcServer *grpc.GRPCServer) *ConfigApp {
	return &ConfigApp{
		GRPCServer: grpcServer,
	}
}

var (
	cfg          *viper.Viper
	configLogger = flogging.MustGetLogger("config")
)

func InitConfig() error {
	_, err := configloader.LoadRuntimeYML()
	if err != nil {
		return fmt.Errorf("load organization config: %w", err)
	}

	// YAML sở hữu topology/default; process environment sở hữu DSN
	// và endpoint thay đổi theo lần deploy. Bind tường minh để Viper
	// Unmarshal và common/db cùng đọc một registry duy nhất.
	for key, names := range map[string][]string{
		"database.dsn":             {"ORGANIZATION_DATABASE_URL", "DATABASE_DSN"},
		"rpc.user.address":         {"ORGANIZATION_USER_GRPC_ADDRESS"},
		"rpc.bdspro.address":       {"ORGANIZATION_BDSPRO_GRPC_ADDRESS"},
		"rpc.notification.address": {"ORGANIZATION_NOTIFICATION_GRPC_ADDRESS"},
		"rpc.payment.address":      {"ORGANIZATION_PAYMENT_GRPC_ADDRESS"},
		"rpc.chat.address":         {"ORGANIZATION_CHAT_GRPC_ADDRESS"},
		"rpc.transaction.address":  {"ORGANIZATION_TRANSACTION_GRPC_ADDRESS"},
		"rpc.auth.address":         {"ORGANIZATION_AUTH_GRPC_ADDRESS"},
		"grpc.port":                {"ORGANIZATION_GRPC_PORT"},
	} {
		args := append([]string{key}, names...)
		if err := viper.BindEnv(args...); err != nil {
			return fmt.Errorf("bind organization config %s: %w", key, err)
		}
	}
	cfg = viper.GetViper()
	if err := readConfig(); err != nil {
		return fmt.Errorf("read organization config: %w", err)
	}
	return nil
}

func readConfig() error {
	dbType := cfg.GetString("Database.PRIMARY_DB_TYPE")
	env.PRIMARY_DB_TYPE = env.DBType(dbType)
	if env.PRIMARY_DB_TYPE != env.DBTypePostgres {
		configLogger.Warnf("Unsupported database type: %s, using default: postgres", dbType)
		env.PRIMARY_DB_TYPE = env.DBTypePostgres
	}

	if dns := cfg.GetString("database.dsn"); dns != "" {
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

	return nil
}

func RPCAddress(service string) (string, error) {
	if cfg == nil {
		return "", fmt.Errorf("organization config is not initialized")
	}
	key := "rpc." + strings.TrimSpace(service) + ".address"
	value := strings.TrimSpace(cfg.GetString(key))
	if value == "" {
		return "", fmt.Errorf("config %s is required", key)
	}
	return value, nil
}
