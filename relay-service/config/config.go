package config

import (
	"common"
	"os"
	"path"
	"time"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/spf13/viper"
)

// LoadConfig đọc cấu hình từ file config.yml
func LoadConfig() {
	ENV_RUNTIME := os.Getenv("ENV_RUNTIME")
	common.LoadYMLConfigFile(ENV_RUNTIME)
}

var (
	cfg          *viper.Viper
	configLogger = flogging.MustGetLogger("configs.config")

	DB_TYPE             string
	DNS                 string
	TRANSACTION_TIMEOUT time.Duration

	IDLE_TIMEOUT  time.Duration
	READ_TIMEOUT  time.Duration
	WRITE_TIMEOUT time.Duration

	REDIS_HOST     string
	REDIS_PASSWORD string

	SWAGGER_PATH     string
	NUMBER_OF_WORKER int = 10
	MAX_PARTICIPANT  int = 10
	ENV_RUNTIME      string
	PORT             string
)

func loadConfig() {
	file := "app.yml"
	if envFile := os.Getenv("CONFIG_FILE"); envFile != "" {
		file = envFile
	}
	flogging.ActivateSpec("debug")

	cfg = viper.New()
	cfg.SetConfigType("yaml")
	cfg.SetConfigName(path.Base(file)) // Chỉ lấy tên file, không có extension
	cfg.AddConfigPath(path.Dir(file))
	cfg.AddConfigPath("./config/")
	cfg.AddConfigPath("../config/")

	configLogger.Infof("Loading config file: %s", file)

	if err := cfg.ReadInConfig(); err != nil {
		configLogger.Fatalf("Error parsing config file: %s", err)
	}

	// Đọc giá trị từ file vào biến
	DB_TYPE = cfg.GetString("database.type")
	DNS = cfg.GetString("database.dns")
	TRANSACTION_TIMEOUT = cfg.GetDuration("timeout.transaction")
	IDLE_TIMEOUT = cfg.GetDuration("timeout.idle")
	READ_TIMEOUT = cfg.GetDuration("timeout.read")
	WRITE_TIMEOUT = cfg.GetDuration("timeout.write")

	REDIS_HOST = cfg.GetString("redis.host")
	REDIS_PASSWORD = cfg.GetString("redis.password")

	SWAGGER_PATH = cfg.GetString("swagger.path")

	NUMBER_OF_WORKER = cfg.GetInt("worker.number")
	MAX_PARTICIPANT = cfg.GetInt("max_participant")
	ENV_RUNTIME = os.Getenv("ENV_RUNTIME")
	PORT = cfg.GetString("server.port")
}

func InitConfig() {
	loadConfig()
}
