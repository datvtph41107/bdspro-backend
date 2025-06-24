//go:build wireinject
// +build wireinject

package grpcserver

import (
	"github.com/google/wire"
	"github.com/tuandq2112/appointment-service/internal/infrastructure/repository"
	"github.com/tuandq2112/appointment-service/internal/interface/handler"
	"github.com/tuandq2112/appointment-service/internal/interface/transformer"
	"github.com/tuandq2112/appointment-service/internal/interface/validator"
	appointmentreminder "github.com/tuandq2112/appointment-service/internal/job/appointment_reminder"
	"github.com/tuandq2112/appointment-service/internal/server/grpc"
	"github.com/tuandq2112/appointment-service/internal/usecase"
	emailsender "github.com/tuandq2112/appointment-service/pkg/email_sender"
)

var repositorySet = wire.NewSet(
	repository.NewAppointmentRepository,
	repository.NewAppointmentReminderRepository,
	// repository.NewTransactionRepository,
)

var usecaseSet = wire.NewSet(
	usecase.NewAppointmentUsecase,
	usecase.NewAppointmentReminderUsecase,
	// usecase.NewTransactionUsecase,
)

var transformerSet = wire.NewSet(
	transformer.NewAppointmentTransformer,
	// transformer.NewTransactionTransformer,
)

var validatorSet = wire.NewSet(
	validator.NewAppointmentValidator,
	// validator.NewTransactionValidator,
)

var handlerSet = wire.NewSet(
	handler.NewAppointmentHandler,
	// handler.NewTransactionHandler,
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
