package config

import (
	"common/configloader"
	qhprorpc "common/rpc"
	"common/rpcenv"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type DatabaseRuntime struct {
	DSN string
}

type RedisRuntime struct {
	Address  string
	Password string
	DB       int
}

type SecurityRuntime struct {
	ProtectedMethods []string
}

type Runtime struct {
	GRPCPort              int
	ServerName            string
	UserRPCTarget         string
	BDSProRPCTarget       string
	NotificationRPCTarget string
	Database              DatabaseRuntime
	Redis                 RedisRuntime
	Security              SecurityRuntime
	Transport             qhprorpc.TransportConfig
}

// NewRuntime materializes Hub process configuration from the registry prepared
// by LoadConfig plus process environment values. Call it once from the gRPC
// composition root and inject the immutable snapshot into resource providers.
func NewRuntime() (Runtime, error) {
	port := viper.GetInt("server.tcp_port")
	if port <= 0 || port > 65535 {
		return Runtime{}, fmt.Errorf("invalid hub gRPC port %d", port)
	}

	serverName, err := configloader.RequiredString("server.name")
	if err != nil {
		return Runtime{}, err
	}
	userTarget, err := configloader.RequiredString("rpc.user.address")
	if err != nil {
		return Runtime{}, err
	}
	bdsproTarget, err := configloader.RequiredString("rpc.bdspro.address")
	if err != nil {
		return Runtime{}, err
	}
	notificationTarget, err := configloader.RequiredString("rpc.notification.address")
	if err != nil {
		return Runtime{}, err
	}
	databaseDSN, err := configloader.RequiredString("database.dsn")
	if err != nil {
		return Runtime{}, err
	}
	redisAddress, err := configloader.RequiredString("redis.host")
	if err != nil {
		return Runtime{}, err
	}
	redisDB := viper.GetInt("redis.db")
	if redisDB < 0 {
		return Runtime{}, fmt.Errorf("invalid hub redis db %d", redisDB)
	}

	protectedMethods := append(
		[]string(nil),
		viper.GetStringSlice("security.protected_methods")...,
	)

	return Runtime{
		GRPCPort:              port,
		ServerName:            strings.TrimSpace(serverName),
		UserRPCTarget:         strings.TrimSpace(userTarget),
		BDSProRPCTarget:       strings.TrimSpace(bdsproTarget),
		NotificationRPCTarget: strings.TrimSpace(notificationTarget),
		Database: DatabaseRuntime{
			DSN: strings.TrimSpace(databaseDSN),
		},
		Redis: RedisRuntime{
			Address:  strings.TrimSpace(redisAddress),
			Password: viper.GetString("redis.pass"),
			DB:       redisDB,
		},
		Security: SecurityRuntime{
			ProtectedMethods: protectedMethods,
		},
		Transport: rpcenv.LoadTransportConfig(),
	}, nil
}
