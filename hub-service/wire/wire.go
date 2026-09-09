//go:build wireinject
// +build wireinject

package wire

import (
	_db "common/db"
	common_injection "common/injection"
	common_utils "common/utils"
	"github.com/google/wire"
	"hub/config"
	hub_infra_client "hub/infra/client"
	hub_infra_db "hub/infra/db"
	hub_infra_handler "hub/infra/handler"
	hub_infra_mapper "hub/infra/mapper"
	hub_infra_postgre "hub/infra/postgre"
	hub_infra_provider "hub/infra/provider"
	hub_infra_redis "hub/infra/redis"
	hub_infra_rpc "hub/infra/rpc"
	hub_infra_setting "hub/infra/setting"
	hub_initial "hub/initial"
	hub_internal_interface "hub/internal/interface"
	hub_internal_repo "hub/internal/repo"
	hub_internal_usecase "hub/internal/usecase"
)

// Inject dependencies
var wireSet = wire.NewSet(
	hub_infra_db.NewDB,
	_db.NewTransactionRepo,
	wire.Bind(new(hub_internal_interface.INotificationClient), new(*hub_infra_client.NotificationClient)),
	wire.Bind(new(hub_internal_interface.IUserClient), new(*hub_infra_client.UserClient)),
	wire.Bind(new(hub_internal_repo.IEventQueueRepo), new(*hub_infra_postgre.EventQueueRepo)),
	wire.Bind(new(hub_internal_repo.IUpdateDataRepo), new(*hub_infra_postgre.UpdateDataRepo)),
	common_injection.NewHttpClient,
	common_utils.NewSyncUtil,
	hub_infra_rpc.NewAdminUserProfileClient,
	hub_infra_handler.NewApiKeyHandler,
	hub_infra_mapper.NewApiKeyMapper,
	hub_infra_postgre.NewApiKeyRepo,
	hub_internal_usecase.NewApiKeyUsecase,
	hub_infra_handler.NewApplinkHandler,
	hub_infra_postgre.NewApplinkPostgre,
	hub_internal_usecase.NewApplinkUsecase,
	hub_infra_client.NewBDSProClient,
	hub_infra_rpc.NewBdsproRPCClient,
	hub_infra_postgre.NewDistrictRepo,
	hub_infra_handler.NewErrorLogHandler,
	hub_infra_mapper.NewErrorLogMapper,
	hub_infra_postgre.NewErrorLogRepo,
	hub_internal_usecase.NewErrorLogUsecase,
	hub_infra_mapper.NewEventQueueMapper,
	hub_infra_postgre.NewEventQueueRepo,
	hub_infra_handler.NewEventQueueService,
	hub_internal_usecase.NewEventQueueUsecase,
	hub_infra_handler.NewFAQHandler,
	hub_infra_mapper.NewFAQMapper,
	hub_infra_postgre.NewFAQRepo,
	hub_internal_usecase.NewFAQUsecase,
	hub_initial.NewInitialApp,
	hub_infra_handler.NewInteractiveEventHandler,
	hub_infra_mapper.NewInteractiveEventMapper,
	hub_infra_postgre.NewInteractiveEventRepo,
	hub_internal_usecase.NewInteractiveEventUsecase,
	hub_infra_handler.NewInternalHandler,
	hub_infra_handler.NewLocationHandler,
	hub_infra_mapper.NewLocationMapper,
	hub_internal_usecase.NewLocationUsecase,
	hub_infra_handler.NewLocationV2Handler,
	hub_infra_mapper.NewLocationV2Mapper,
	hub_internal_usecase.NewLocationV2Usecase,
	hub_infra_client.NewNotificationClient,
	hub_infra_rpc.NewNotificationRPCClient,
	hub_infra_postgre.NewProvinceRepo,
	hub_infra_postgre.NewProvinceV2Repo,
	hub_infra_redis.NewRedisProvider,
	hub_infra_handler.NewSystemConfigHandler,
	hub_infra_setting.NewSystemConfigPersist,
	hub_infra_postgre.NewSystemConfigPostgres,
	hub_internal_usecase.NewSystemConfigUsecase,
	hub_infra_handler.NewUpdateDataHandler,
	hub_infra_provider.NewUpdateDataProviderImpl,
	hub_infra_postgre.NewUpdateDataRepo,
	hub_internal_usecase.NewUpdateDataUsecase,
	hub_infra_client.NewUserClient,
	hub_infra_handler.NewUserGuideHandler,
	hub_infra_mapper.NewUserGuideMapper,
	hub_infra_postgre.NewUserGuideRepo,
	hub_infra_postgre.NewUserGuideStepRepo,
	hub_internal_usecase.NewUserGuideUsecase,
	hub_infra_handler.NewVersionHandler,
	hub_infra_mapper.NewVersionMapper,
	hub_infra_postgre.NewVersionRepo,
	hub_internal_usecase.NewVersionUsecase,
	hub_infra_postgre.NewWardRepo,
	hub_infra_postgre.NewWardV2Repo,
)

func InitializeApp(runtime config.Runtime) (*hub_initial.InitialApp, func(), error) {
	wire.Build(wireSet)
	return nil, nil, nil
}
