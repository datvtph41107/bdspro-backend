package config

import (
	"payment/internal/interface/handler"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/spf13/viper"
)

type ConfigApp struct {
	Logger      *flogging.FabricLogger
	BankHandler *handler.BankHandler
}

func NewConfigApp(
	bankHandler *handler.BankHandler,
) *ConfigApp {
	return &ConfigApp{
		BankHandler: bankHandler,
		Logger:      flogging.MustGetLogger(viper.GetString("server.name")),
	}
}
