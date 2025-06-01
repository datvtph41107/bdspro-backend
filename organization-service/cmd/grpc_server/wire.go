//go:build wireinject
// +build wireinject

package grpcserver

import (
	"github.com/google/wire"
	"github.com/tuandq2112/organization-service/internal/infrastructure/repository"
	"github.com/tuandq2112/organization-service/internal/interface/handler"
	"github.com/tuandq2112/organization-service/internal/interface/transformer"
	"github.com/tuandq2112/organization-service/internal/interface/validator"
	"github.com/tuandq2112/organization-service/internal/server/grpc"
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
	validator.NewOrganizationValidator,
	handler.NewOrganizationHandler,
)

func InitializeGRPCServer() *grpc.GRPCServer {
	wire.Build(
		repositorySet,
		usecaseSet,
		handlerSet,
		grpc.NewGRPCServer,
	)
	return &grpc.GRPCServer{}
}
