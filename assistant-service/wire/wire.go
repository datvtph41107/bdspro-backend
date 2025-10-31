//go:build wireinject
// +build wireinject

package wire

import (
	assistant_infra_client "assistant/infra/client"
	assistant_infra_handler "assistant/infra/handler"
	"assistant/initial"
	assistant_initial "assistant/initial"
	assistant_internal_interface_provider "assistant/internal/interface/provider"
	assistant_internal_usecases "assistant/internal/usecases"
	_db "common/db"

	"github.com/google/wire"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	_db.NewTransactionRepo,
	wire.Bind(new(assistant_internal_interface_provider.DeepseekProvider), new(*assistant_infra_client.DeepseekClient)),
	wire.Bind(new(assistant_internal_interface_provider.OpenAIProvider), new(*assistant_infra_client.OpenAIClient)),
	assistant_infra_handler.NewAssistantHandler,
	assistant_internal_usecases.NewAssistantUsecase,
	assistant_infra_client.NewDeepseekClient,
	assistant_initial.NewInitialApp,
	assistant_infra_client.NewOpenAIClient,
)

func InitializeApp() (*initial.InitialApp, error) {
	wire.Build(wireSet)
	return nil, nil
}
