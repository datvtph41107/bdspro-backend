package repository

import (
	"github.com/tuandq2112/appointment-service/appconfig"
	"github.com/tuandq2112/appointment-service/internal/domain/repository"
	"github.com/tuandq2112/appointment-service/internal/infrastructure/database"
	"github.com/tuandq2112/appointment-service/internal/infrastructure/repository/postgres"
)

func NewAppointmentRepository() repository.AppointmentRepository {
	db := database.NewPostgresDB(appconfig.DB_DNS)
	db.AutoMigrate(&postgres.AppointmentModel{})
	return postgres.NewAppointmentPostgresRepository(db)
}
func NewAppointmentReminderRepository() repository.AppointmentReminderRepository {
	db := database.NewPostgresDB(appconfig.DB_DNS)
	db.AutoMigrate(&postgres.AppointmentReminderModel{}, &postgres.AppointmentParticipantModel{})
	return postgres.NewAppointmentReminderPostgresRepository(db)
}