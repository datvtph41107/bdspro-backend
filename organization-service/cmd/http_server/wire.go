//go:build wireinject
// +build wireinject

package httpserver

import (
	"github.com/google/wire"
	"github.com/tuandq2112/organization-service/internal/infrastructure/repository"
	"github.com/tuandq2112/organization-service/internal/interface/handler"
	"github.com/tuandq2112/organization-service/internal/interface/transformer"
	"github.com/tuandq2112/organization-service/internal/server/http"
	"github.com/tuandq2112/organization-service/internal/usecase"
)

var repositorySet = wire.NewSet(
	repository.NewRepository,
)

var usecaseSet = wire.NewSet(
	usecase.NewOrganizationUsecase,
)

var handlerSet = wire.NewSet(
	transformer.NewOrganizationTransformer,
	handler.NewOrganizationHandler,
)

func InitializeHttpServer() *http.HttpServer {
	wire.Build(
		repositorySet,
		usecaseSet,
		handlerSet,
		http.NewHttpServer,
	)
	return &http.HttpServer{}
}
