package initial

import (
	"utility/infra/handler"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/spf13/viper"
)

type InitialApp struct {
	EventQueueService        *handler.EventQueueService
	UtilityInternalService   *handler.UtilityInternalService
	Logger                   *flogging.FabricLogger
}

func NewInitialApp(
	eventQueueService *handler.EventQueueService,
	utilityInternalService *handler.UtilityInternalService,
) *InitialApp {
	return &InitialApp{
		EventQueueService:      eventQueueService,
		UtilityInternalService: utilityInternalService,
		Logger:                 flogging.MustGetLogger(viper.GetString("server.name")),
	}
}

