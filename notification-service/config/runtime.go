package config

import "noti/infra/service"

type ConfigApp struct {
	Noti       *service.InternalNotiService
	ServerHttp *service.GatewayNotiService
}

func NewConfigApp(
	router *service.InternalNotiService,
	ServerHttp *service.GatewayNotiService,
) *ConfigApp {
	return &ConfigApp{
		Noti:       router,
		ServerHttp: ServerHttp,
	}
}
