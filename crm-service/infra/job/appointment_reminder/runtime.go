package appointmentreminder

import (
	"context"
	"log"
	"time"

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
	reminders, err := j.appointmentReminderUsecase.GetReminders(context.Background(), 10)
	if err != nil {
		log.Println("Error getting reminders:", err)
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
