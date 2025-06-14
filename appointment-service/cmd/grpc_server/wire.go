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
	"github.com/tuandq2112/appointment-service/internal/job/appointment_reminder"
	emailsender "github.com/tuandq2112/appointment-service/pkg/email_sender"
)

var repositorySet = wire.NewSet(
	repository.NewAppointmentRepository,
	repository.NewAppointmentReminderRepository,
)

var usecaseSet = wire.NewSet(
	usecase.NewAppointmentUsecase,
	usecase.NewAppointmentReminderUsecase,
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

var emailSenderSet = wire.NewSet(
	emailsender.NewEmailSender,
)

var jobSet = wire.NewSet(
	appointmentreminder.NewAppointmentReminderJob,
)

func InitializeGRPCServer() *grpc.GRPCServer {
	wire.Build(
		repositorySet,
		usecaseSet,
		transformerSet,
		validatorSet,
		emailSenderSet,
		handlerSet,
		jobSet,
		grpc.NewGRPCServer,
	)
	return &grpc.GRPCServer{}
}
