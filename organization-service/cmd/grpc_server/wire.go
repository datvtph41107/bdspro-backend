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
	repository.NewGroupRepository,
	repository.NewGroupDealRepository,
	repository.NewGroupLogActivityRepository,
	repository.NewGroupNotificationRepository,
	repository.NewGroupSettingRepository,
	repository.NewGroupMemberRepository,
	repository.NewGroupDocumentRepository,
	repository.NewOrganizationLogActivityRepository,
)


var usecaseSet = wire.NewSet(
	usecase.NewOrganizationLogActivityUsecase,
	usecase.NewOrganizationLogWorker,
	usecase.NewOrganizationAuthUsecase,
	usecase.NewOrganizationUsecase,
	usecase.NewOrganizationMemberUsecase,
	usecase.NewOrganizationRoleUsecase,
	usecase.NewOrganizationPermissionUsecase,
	usecase.NewOrganizationBranchUsecase,
	usecase.NewGroupLogActivityUsecase,
	usecase.NewLogWorker,
	usecase.NewGroupUsecase,
	usecase.NewGroupDealUsecase,
	usecase.NewGroupNotificationUsecase,
	usecase.NewGroupSettingUsecase,
	usecase.NewGroupMemberUsecase,
	usecase.NewGroupDocumentUsecase,
)


var transformerSet = wire.NewSet(
	transformer.NewOrganizationTransformer,
	transformer.NewOrganizationMemberTransformer,
	transformer.NewOrganizationRoleTransformer,
	transformer.NewOrganizationPermissionTransformer,
	transformer.NewOrganizationBranchTransformer,
	transformer.NewGroupTransformer,
	transformer.NewGroupDealTransformer,
	transformer.NewGroupLogActivityTransformer,
	transformer.NewGroupNotificationTransformer,
	transformer.NewGroupSettingTransformer,
	transformer.NewGroupMemberTransformer,
	transformer.NewGroupDocumentTransformer,
	transformer.NewOrganizationLogActivityTransformer,
)
var validatorSet = wire.NewSet(
	validator.NewOrganizationValidator,
	validator.NewOrganizationMemberValidator,
	validator.NewOrganizationRoleValidator,
	validator.NewOrganizationPermissionValidator,
	validator.NewOrganizationBranchValidator,
	validator.NewGroupValidator,
	validator.NewGroupDealValidator,
	validator.NewGroupLogActivityValidator,
	validator.NewGroupNotificationValidator,
	validator.NewGroupSettingValidator,
	validator.NewGroupMemberValidator,
	validator.NewGroupDocumentValidator,
	validator.NewOrganizationLogActivityValidator,
)

var handlerSet = wire.NewSet(
	handler.NewOrganizationHandler,
	handler.NewOrganizationMemberHandler,
	handler.NewOrganizationRoleHandler,
	handler.NewOrganizationPermissionHandler,
	handler.NewOrganizationBranchHandler,
	handler.NewGroupHandler,
	handler.NewGroupDealHandler,
	handler.NewGroupLogActivityHandler,
	handler.NewGroupNotificationHandler,
	handler.NewGroupSettingHandler,
	handler.NewGroupMemberHandler,
	handler.NewGroupDocumentHandler,
	handler.NewOrganizationLogActivityHandler,
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
