package config

import "noti/infra/service"

type ConfigApp struct {
	InternalNotiService *service.InternalNotiService
	ServerHttp          *service.GatewayNotiService
	HistoryService      *service.HistoryService
}

func NewConfigApp(
	router *service.InternalNotiService,
	ServerHttp *service.GatewayNotiService,
	historyService *service.HistoryService,
) *ConfigApp {
	return &ConfigApp{
		InternalNotiService: router,
		ServerHttp:          ServerHttp,
		HistoryService:      historyService,
	}
}
