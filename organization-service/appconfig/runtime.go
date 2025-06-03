package appconfig

import (
	"os"
	"path"
	"strings"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/spf13/viper"
)

const (
	USER_CONTEXT = "user_context"
	SERVICE_NAME = "organization-service"
)

type DBType string

const (
	DBTypePostgres DBType = "postgres"
	
	// Default values
	DefaultDBDNS    = "host=localhost user=postgres password=postgres dbname=organization port=5432 sslmode=disable"
	DefaultHTTPPort = "8080"
	DefaultGRPCPort = "50052"
	ADMIN_ROLE_KEY  = "admin"
)

var (
	cfg          *viper.Viper
	configLogger = flogging.MustGetLogger("config")
	PRIMARY_DB_TYPE               = DBTypePostgres
	DB_DNS                        = DefaultDBDNS
	HTTP_PORT                     = DefaultHTTPPort
	WHITELIST_METHODS_GET_PATH    []string
	WHITELIST_METHODS_POST_PATH   []string
	WHITELIST_METHODS_PUT_PATH    []string
	WHITELIST_METHODS_DELETE_PATH []string
	GRPC_PORT                     = DefaultGRPCPort
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

func validateWhitelistPaths(paths []string) []string {
	validPaths := make([]string, 0)
	for _, p := range paths {
		// Ensure path starts with /
		if !strings.HasPrefix(p, "/") {
			p = "/" + p
		}
		validPaths = append(validPaths, p)
	}
	return validPaths
}

func readConfig() error {
	// Read and validate database type
	dbType := cfg.GetString("Database.PRIMARY_DB_TYPE")
	PRIMARY_DB_TYPE = DBType(dbType)
	if PRIMARY_DB_TYPE != DBTypePostgres {
		configLogger.Warnf("Unsupported database type: %s, using default: postgres", dbType)
		PRIMARY_DB_TYPE = DBTypePostgres
	}

	// Read and validate database DNS
	if dns := cfg.GetString("Database.DB_DNS"); dns != "" {
		DB_DNS = dns
	} else {
		configLogger.Warnf("Database DNS not specified, using default: %s", DefaultDBDNS)
	}

	// Read and validate HTTP port
	if port := cfg.GetString("HTTP.PORT"); port != "" {
		HTTP_PORT = port
	} else {
		configLogger.Warnf("HTTP port not specified, using default: %s", DefaultHTTPPort)
	}

	// Read and validate whitelist methods paths
	WHITELIST_METHODS_GET_PATH = validateWhitelistPaths(cfg.GetStringSlice("WhitelistMethods.GET_PATH"))
	WHITELIST_METHODS_POST_PATH = validateWhitelistPaths(cfg.GetStringSlice("WhitelistMethods.POST_PATH"))
	WHITELIST_METHODS_PUT_PATH = validateWhitelistPaths(cfg.GetStringSlice("WhitelistMethods.PUT_PATH"))
	WHITELIST_METHODS_DELETE_PATH = validateWhitelistPaths(cfg.GetStringSlice("WhitelistMethods.DELETE_PATH"))

	// Log configuration summary
	configLogger.Infof("Configuration loaded successfully:")
	configLogger.Infof("- Database Type: %s", PRIMARY_DB_TYPE)
	configLogger.Infof("- HTTP Port: %s", HTTP_PORT)
	configLogger.Infof("- Whitelisted GET paths: %v", WHITELIST_METHODS_GET_PATH)
	configLogger.Infof("- Whitelisted POST paths: %v", WHITELIST_METHODS_POST_PATH)
	configLogger.Infof("- Whitelisted PUT paths: %v", WHITELIST_METHODS_PUT_PATH)
	configLogger.Infof("- Whitelisted DELETE paths: %v", WHITELIST_METHODS_DELETE_PATH)

	return nil
}