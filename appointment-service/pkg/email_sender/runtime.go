package emailsender

import (
	"github.com/tuandq2112/appointment-service/appconfig"
	"gopkg.in/mail.v2"
)

type EmailSender struct {
	From     string
	Name     string 
	Host     string
	Port     int
	Username string
	Password string
}

func NewEmailSender() *EmailSender {
	return &EmailSender{
		Host:     appconfig.EMAIL_HOST,
		Port:     appconfig.EMAIL_PORT,
		Username: appconfig.EMAIL_USERNAME,
		Password: appconfig.EMAIL_PASSWORD,
		From:     appconfig.EMAIL_FROM,
		Name:     appconfig.EMAIL_NAME,
	}
}

func (s *EmailSender) Send(to []string, subject, body string, isHTML bool) error {
	m := mail.NewMessage()
	m.SetHeader("From", m.FormatAddress(s.From, s.Name))
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)

	if isHTML {
		m.SetBody("text/html", body)
	} else {
		m.SetBody("text/plain", body)
	}

	d := mail.NewDialer(s.Host, s.Port, s.Username, s.Password)

	return d.DialAndSend(m)
}
