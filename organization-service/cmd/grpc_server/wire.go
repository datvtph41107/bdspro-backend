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
	repository.NewOrganizationRepository,
	repository.NewOrganizationMemberRepository,
)

var usecaseSet = wire.NewSet(
	usecase.NewOrganizationUsecase,
	usecase.NewOrganizationMemberUsecase,
)

var handlerSet = wire.NewSet(
	transformer.NewOrganizationTransformer,
	validator.NewOrganizationValidator,
	handler.NewOrganizationHandler,
	handler.NewOrganizationMemberHandler,
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
