//go:build wireinject
// +build wireinject

package wire

import (
	assistant_infra_client "assistant/infra/client"
	assistant_infra_handler "assistant/infra/handler"
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
	wire.Bind(new(assistant_internal_interface_provider.HubProvider), new(*assistant_infra_client.HubClient)),
	common_injection.NewAuthClient,
	common_injection.NewBDSProClient,
	common_injection.NewCRMClient,
	common_injection.NewChatClient,
	common_injection.NewHubClient,
	common_injection.NewNotification,
	common_injection.NewOrganizationClient,
	common_injection.NewPaymentClient,
	common_injection.NewTQClient,
	common_injection.NewTransactionClient,
	common_injection.NewUserClient,
	common_injection.NewHttpClient,
	common_utils.NewSyncUtil,
	assistant_infra_handler.NewAssistantHandler,
	assistant_internal_usecases.NewAssistantUsecase,
	assistant_infra_client.NewBdsproInternalClient,
	assistant_infra_client.NewDeepseekClient,
	assistant_infra_client.NewFileClient,
	assistant_infra_client.NewGeminiClient,
	assistant_infra_client.NewHubClient,
	assistant_initial.NewInitialApp,
	assistant_infra_client.NewOpenAIClient,
)

func InitializeApp() (*assistant_initial.InitialApp, error) {
	wire.Build(wireSet)
	return nil, nil
}
