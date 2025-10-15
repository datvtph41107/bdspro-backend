package config

import (
	_ "embed"
	"os"
	"path"

	"marketing/infrastructure/server/grpc"

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
}

func loadConfig() {
	file := "develop.yaml"
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

	// Đồng bộ với global viper
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	// Copy tất cả config từ local viper sang global viper
	for key, value := range cfg.AllSettings() {
		viper.Set(key, value)
	}
}
