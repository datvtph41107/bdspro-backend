package initial

import (
	_redis "common/redis"

	"feedback/infra/handler"
	"feedback/internal/usecase"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/spf13/viper"
)

type InitialApp struct {
	Logger *flogging.FabricLogger

	InternalHandler     *handler.InternalHandler
	RateHandler         *handler.RateHandler
	ReportHandler       *handler.ReportHandler
	ReportAdminHandler  *handler.ReportAdminHandler
	ReportReasonHandler *handler.ReportReasonHandler
	RegistryHandler     *handler.RegistryHandler
	RegistryUsecase     *usecase.RegistryUsecase
}

func NewApp(internalHandler *handler.InternalHandler,
	rateHandler *handler.RateHandler,
	reportHandler *handler.ReportHandler,
	reportAdminHandler *handler.ReportAdminHandler,
	reportReasonHandler *handler.ReportReasonHandler,
	registryHandler *handler.RegistryHandler,
	registryUsecase *usecase.RegistryUsecase) *InitialApp {
	return &InitialApp{
		RegistryUsecase:     registryUsecase,
		InternalHandler:     internalHandler,
		RateHandler:         rateHandler,
		ReportHandler:       reportHandler,
		ReportAdminHandler:  reportAdminHandler,
		ReportReasonHandler: reportReasonHandler,
		RegistryHandler:     registryHandler,
		Logger:              flogging.MustGetLogger(viper.GetString("server.name")),
	}
}

func NewRedisService() *_redis.RedisService {
	return _redis.NewRedisService()
}
