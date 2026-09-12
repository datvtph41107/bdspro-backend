package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
)

// Runtime is the validated process configuration required before the User
// serving process opens resources. Feature-specific business settings remain in
// Properties until their owners are migrated to typed config.
type Runtime struct {
	GRPCAddress string
	Database    DatabaseRuntime
	Redis       RedisRuntime
}

type DatabaseRuntime struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisRuntime struct {
	Address  string
	Password string
	DB       int
}

// LoadRuntime establishes config file/env precedence once and returns the
// process-owned subset needed for resource composition.
func LoadRuntime() (Runtime, error) {
	if err := LoadProperties(); err != nil {
		return Runtime{}, err
	}

	grpcAddress := strings.TrimSpace(os.Getenv("GRPC_ADDRESS"))
	if grpcAddress == "" {
		if Properties.Server.TCPPort <= 0 {
			return Runtime{}, errors.New("grpc address is required")
		}
		grpcAddress = fmt.Sprintf(":%d", Properties.Server.TCPPort)
	}

	database, err := databaseRuntimeFromLoadedProperties()
	if err != nil {
		return Runtime{}, err
	}

	redisAddress := firstNonEmpty(
		os.Getenv("USER_REDIS_ADDRESS"),
		os.Getenv("REDIS_ADDRESS"),
		Properties.Redis.Host,
	)
	if redisAddress == "" || strings.HasPrefix(redisAddress, "CHANGE_ME") {
		return Runtime{}, errors.New("redis address is required")
	}

	jwtKey := strings.TrimSpace(Properties.JWT.SecretKey)
	if jwtKey == "" || strings.HasPrefix(jwtKey, "CHANGE_ME") || len(jwtKey) < 32 {
		return Runtime{}, errors.New("JWT signing key must be provided and at least 32 bytes")
	}

	return Runtime{
		GRPCAddress: grpcAddress,
		Database:    database,
		Redis: RedisRuntime{
			Address:  redisAddress,
			Password: firstNonEmpty(os.Getenv("USER_REDIS_PASSWORD"), os.Getenv("REDIS_PASSWORD"), Properties.Redis.Password),
			DB:       Properties.Redis.DB,
		},
	}, nil
}

// LoadDatabaseRuntime loads only the dependency needed by an explicit
// database command. It deliberately does not require Redis, JWT or a listening
// address because bootstrap/migration processes do not own those resources.
func LoadDatabaseRuntime() (DatabaseRuntime, error) {
	if err := LoadProperties(); err != nil {
		return DatabaseRuntime{}, err
	}
	return databaseRuntimeFromLoadedProperties()
}

func databaseRuntimeFromLoadedProperties() (DatabaseRuntime, error) {
	// Database DSN is owned by the environment YAML:
	// - workstation: config/runtime.yml merged with config/local.yml
	// - deploy: CONFIG_FILE / QHPRO_CONFIG_FILE (develop.yml, release.yml, …)
	dsn := strings.TrimSpace(Properties.Database.DSN)
	if dsn == "" || strings.HasPrefix(dsn, "CHANGE_ME") {
		return DatabaseRuntime{}, errors.New("database.dsn is required in the environment config YAML")
	}
	return DatabaseRuntime{
		DSN:             dsn,
		MaxOpenConns:    Properties.Database.MaxOpenConns,
		MaxIdleConns:    Properties.Database.MaxIdleConns,
		ConnMaxLifetime: Properties.Database.ConnMaxLifetime,
	}, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

// DescribeDatabaseDSN returns a password-free host/db summary for startup logs.
func DescribeDatabaseDSN(dsn string) string {
	dsn = strings.TrimSpace(dsn)
	if dsn == "" {
		return "dsn=<empty>"
	}

	host, port, user, dbname := "", "", "", ""
	for _, part := range strings.Fields(dsn) {
		key, value, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}
		switch strings.ToLower(strings.TrimSpace(key)) {
		case "host":
			host = strings.TrimSpace(value)
		case "port":
			port = strings.TrimSpace(value)
		case "user":
			user = strings.TrimSpace(value)
		case "dbname", "database":
			dbname = strings.TrimSpace(value)
		}
	}
	if host == "" && strings.Contains(dsn, "://") {
		if u, err := url.Parse(dsn); err == nil {
			host = u.Hostname()
			port = u.Port()
			user = u.User.Username()
			dbname = strings.TrimPrefix(u.Path, "/")
		}
	}
	if port == "" {
		port = "5432"
	}
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s", host, port, user, dbname)
}
