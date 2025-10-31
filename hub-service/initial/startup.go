package initial

import (
	"hub/infra/handler"
	"hub/internal/usecase"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/spf13/viper"
)

type InitialApp struct {
	EventQueueService   *handler.EventQueueService
	InternalHandler     *handler.InternalHandler
	LocationHandler     *handler.LocationHandler
	LocationV2Handler   *handler.LocationV2Handler
	UserGuideHandler    *handler.UserGuideHandler
	FAQHandler          *handler.FAQHandler
	SystemConfigHandler *handler.SystemConfigHandler
	SystemConfigUsecase *usecase.SystemConfigUsecase
	Logger              *flogging.FabricLogger
}

func NewInitialApp(
	eventQueueService *handler.EventQueueService,
	internalHandler *handler.InternalHandler,
	locationHandler *handler.LocationHandler,
	locationV2Handler *handler.LocationV2Handler,
	userGuideHandler *handler.UserGuideHandler,
	faqHandler *handler.FAQHandler,
	systemConfigHandler *handler.SystemConfigHandler,
	systemConfigUsecase *usecase.SystemConfigUsecase,
) *InitialApp {
	return &InitialApp{
		EventQueueService:   eventQueueService,
		InternalHandler:     internalHandler,
		LocationHandler:     locationHandler,
		LocationV2Handler:   locationV2Handler,
		UserGuideHandler:    userGuideHandler,
		FAQHandler:          faqHandler,
		SystemConfigHandler: systemConfigHandler,
		SystemConfigUsecase: systemConfigUsecase,
		Logger:              flogging.MustGetLogger(viper.GetString("server.name")),
	}
}
