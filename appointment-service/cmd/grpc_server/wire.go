//go:build wireinject
// +build wireinject

package grpcserver

import (
	"github.com/google/wire"
	"github.com/tuandq2112/appointment-service/internal/infrastructure/repository"
	"github.com/tuandq2112/appointment-service/internal/interface/handler"
	"github.com/tuandq2112/appointment-service/internal/interface/transformer"
	"github.com/tuandq2112/appointment-service/internal/interface/validator"
	"github.com/tuandq2112/appointment-service/internal/server/grpc"
	"github.com/tuandq2112/appointment-service/internal/usecase"
)

var repositorySet = wire.NewSet(
	repository.NewAppointmentRepository,
)

var usecaseSet = wire.NewSet(
	usecase.NewAppointmentUsecase,
)

var transformerSet = wire.NewSet(
	transformer.NewAppointmentTransformer,
)

var validatorSet = wire.NewSet(
	validator.NewAppointmentValidator,
)

var handlerSet = wire.NewSet(
	handler.NewAppointmentHandler,
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
