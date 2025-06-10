package config

import "noti/infra/service"

type ConfigApp struct {
	InternalNotiService *service.InternalNotiService
	ServerHttp          *service.GatewayNotiService
}

func NewConfigApp(
	router *service.InternalNotiService,
	ServerHttp *service.GatewayNotiService,
) *ConfigApp {
	return &ConfigApp{
		InternalNotiService: router,
		ServerHttp:          ServerHttp,
	}
}
