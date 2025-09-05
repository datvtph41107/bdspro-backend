package initial

import "noti/infra/service"

type InitialApp struct {
	InternalNotiService *service.InternalNotiService
	ServerHttp          *service.GatewayNotiService
	HistoryService      *service.HistoryService
}

func NewInitialApp(
	router *service.InternalNotiService,
	ServerHttp *service.GatewayNotiService,
	historyService *service.HistoryService,
) *InitialApp {
	return &InitialApp{
		InternalNotiService: router,
		ServerHttp:          ServerHttp,
		HistoryService:      historyService,
	}
}
