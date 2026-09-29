package appointmentreminder

import (
	"context"
	"log/slog"
	"time"

	"common/logging"
	emailsender "crm/infra/email_sender"
	"crm/internal/usecase"
)

type AppointmentReminderJob struct {
	appointmentReminderUsecase usecase.AppointmentReminderUsecase
	emailSender                *emailsender.EmailSender
}

func NewAppointmentReminderJob(appointmentReminderUsecase usecase.AppointmentReminderUsecase, emailSender *emailsender.EmailSender) *AppointmentReminderJob {
	return &AppointmentReminderJob{
		appointmentReminderUsecase: appointmentReminderUsecase,
		emailSender:                emailSender,
	}
}

func (j *AppointmentReminderJob) Run() {
	ctx := context.Background()
	reminders, err := j.appointmentReminderUsecase.GetReminders(ctx, 10)
	if err != nil {
		logging.WithComponent(ctx, "appointment-reminder-job").Error(
			"get appointment reminders failed",
			slog.Any("error", err),
		)
		return
	}
	for _, reminder := range reminders {
		if reminder.IsSent {
			continue
		}
		if reminder.RemindAt.Before(time.Now()) {
			j.emailSender.Send([]string{"doquoctuan311@gmail.com"}, "Reminder", "You have an appointment reminder", true)
		}
	}
}
