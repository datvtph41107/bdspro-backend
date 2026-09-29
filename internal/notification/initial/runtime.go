package initial

import service "notification/infra/handler"

type InitialApp struct {
	InternalHandler        *service.InternalHandler
	NotificationHandler    *service.NotificationHandler
	HistoryHandler         *service.HistoryHandler
	DealHistoryHandler     *service.DealHistoryHandler
	PropertyHistoryHandler *service.PropertyHistoryHandler
}

func NewInitialApp(
	internalHandler *service.InternalHandler,
	notificationHandler *service.NotificationHandler,
	historyHandler *service.HistoryHandler,
	dealHistoryHandler *service.DealHistoryHandler,
	propertyHistoryHandler *service.PropertyHistoryHandler,
) *InitialApp {
	return &InitialApp{
		InternalHandler:        internalHandler,
		NotificationHandler:    notificationHandler,
		HistoryHandler:         historyHandler,
		DealHistoryHandler:     dealHistoryHandler,
		PropertyHistoryHandler: propertyHistoryHandler,
	}
}
