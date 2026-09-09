//go:build wireinject
// +build wireinject

package wire

import (
	assistant_config "assistant/config"
	assistant_infra_client "assistant/infra/client"
	assistant_infra_handler "assistant/infra/handler"
	assistant_infra_rpc "assistant/infra/rpc"
	assistant_initial "assistant/initial"
	assistant_internal_interface_provider "assistant/internal/interface/provider"
	assistant_internal_usecases "assistant/internal/usecases"
	_db "common/db"
	common_injection "common/injection"
	_provider "common/provider"
	common_utils "common/utils"
	"github.com/google/wire"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	_db.NewTransactionRepo,
	_provider.SyncProviderSet,
	wire.Bind(new(assistant_internal_interface_provider.BdsproInternalProvider), new(*assistant_infra_client.BdsproInternalClient)),
	wire.Bind(new(assistant_internal_interface_provider.DeepseekProvider), new(*assistant_infra_client.DeepseekClient)),
	wire.Bind(new(assistant_internal_interface_provider.GeminiProvider), new(*assistant_infra_client.GeminiClient)),
	wire.Bind(new(assistant_internal_interface_provider.HubProvider), new(*assistant_infra_client.HubClient)),
	wire.Bind(new(assistant_internal_interface_provider.OpenAIProvider), new(*assistant_infra_client.OpenAIClient)),
	common_injection.NewHttpClient,
	common_utils.NewSyncUtil,
	assistant_infra_handler.NewAssistantHandler,
	assistant_internal_usecases.NewAssistantUsecase,
	assistant_infra_client.NewBdsproInternalClient,
	assistant_infra_rpc.NewBdsproRPCClient,
	assistant_infra_client.NewDeepseekClient,
	assistant_infra_client.NewGeminiClient,
	assistant_infra_client.NewHubClient,
	assistant_infra_rpc.NewHubRPCClient,
	assistant_initial.NewInitialApp,
	assistant_infra_client.NewOpenAIClient,
	assistant_config.NewRuntime,
)

func InitializeApp() (*assistant_initial.InitialApp, func(), error) {
	wire.Build(wireSet)
	return nil, nil, nil
}
