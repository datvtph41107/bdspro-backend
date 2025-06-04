//go:build wireinject
// +build wireinject

package grpcserver

import (
	"github.com/google/wire"
	"github.com/tuandq2112/organization-service/internal/infra/repository"
	"github.com/tuandq2112/organization-service/internal/interface/handler"
	"github.com/tuandq2112/organization-service/internal/interface/transformer"
	"github.com/tuandq2112/organization-service/internal/interface/validator"
	"github.com/tuandq2112/organization-service/internal/server/grpc"
	"github.com/tuandq2112/organization-service/internal/usecase"
)

var repositorySet = wire.NewSet(
	repository.NewOrganizationRepository,
	repository.NewOrganizationMemberRepository,
	repository.NewOrganizationRoleRepository,
	repository.NewOrganizationPermissionRepository,
	repository.NewOrganizationBranchRepository,
	repository.NewOrganizationBranchMemberRepository,
)

var usecaseSet = wire.NewSet(
	usecase.NewOrganizationAuthUsecase,
	usecase.NewOrganizationUsecase,
	usecase.NewOrganizationMemberUsecase,
	usecase.NewOrganizationRoleUsecase,
	usecase.NewOrganizationPermissionUsecase,
	usecase.NewOrganizationBranchUsecase,
)
var transformerSet = wire.NewSet(
	transformer.NewOrganizationTransformer,
	transformer.NewOrganizationMemberTransformer,
	transformer.NewOrganizationRoleTransformer,
	transformer.NewOrganizationPermissionTransformer,
	transformer.NewOrganizationBranchTransformer,
)
var validatorSet = wire.NewSet(
	validator.NewOrganizationValidator,
	validator.NewOrganizationMemberValidator,
	validator.NewOrganizationRoleValidator,
	validator.NewOrganizationPermissionValidator,
	validator.NewOrganizationBranchValidator,
)

var handlerSet = wire.NewSet(
	handler.NewOrganizationHandler,
	handler.NewOrganizationMemberHandler,
	handler.NewOrganizationRoleHandler,
	handler.NewOrganizationPermissionHandler,
	handler.NewOrganizationBranchHandler,
)

func InitializeGRPCServer() *grpc.GRPCServer {
	wire.Build(
		repositorySet,
		usecaseSet,
		transformerSet,
		validatorSet,
		handlerSet,
		grpc.NewGRPCServer,
	)
	return &grpc.GRPCServer{}
}
