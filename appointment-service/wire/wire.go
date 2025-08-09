//go:build wireinject
// +build wireinject
	
package wire

import (
	"github.com/google/wire"
	_db "common/db"
	"appointment/config"
	appointment_internal_infrastructure_client "appointment/internal/infrastructure/client"
	appointment_internal_infrastructure_database "appointment/internal/infrastructure/database"
	appointment_internal_infrastructure_repository "appointment/internal/infrastructure/repository"
	appointment_internal_infrastructure_repository_postgres "appointment/internal/infrastructure/repository/postgres"
	appointment_internal_interface_handler "appointment/internal/interface/handler"
	appointment_internal_interface_transformer "appointment/internal/interface/transformer"
	appointment_internal_interface_validator "appointment/internal/interface/validator"
	appointment_internal_job_appointment_reminder "appointment/internal/job/appointment_reminder"
	appointment_internal_server_grpc "appointment/internal/server/grpc"
	appointment_internal_usecase "appointment/internal/usecase"
	appointment_pkg_email_sender "appointment/pkg/email_sender"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	appointment_internal_interface_handler.NewAppointmentHandler,
	appointment_internal_infrastructure_repository_postgres.NewAppointmentPostgresRepository,
	appointment_internal_job_appointment_reminder.NewAppointmentReminderJob,
	appointment_internal_infrastructure_repository_postgres.NewAppointmentReminderPostgresRepository,
	appointment_internal_infrastructure_repository.NewAppointmentReminderRepository,
	appointment_internal_usecase.NewAppointmentReminderUsecase,
	appointment_internal_infrastructure_repository.NewAppointmentRepository,
	appointment_internal_interface_transformer.NewAppointmentTransformer,
	appointment_internal_usecase.NewAppointmentUsecase,
	appointment_internal_interface_validator.NewAppointmentValidator,
	appointment_internal_infrastructure_client.NewBdsproClient,
	appointment_pkg_email_sender.NewEmailSender,
	appointment_internal_server_grpc.NewGRPCServer,
	appointment_internal_infrastructure_client.NewOrganizationClient,
	appointment_internal_infrastructure_database.NewPostgresDB,
	appointment_internal_infrastructure_client.NewUserClient,
)
func InitializeApp() (*config.ConfigApp, error) {
	wire.Build(wireSet)
	return nil, nil
}