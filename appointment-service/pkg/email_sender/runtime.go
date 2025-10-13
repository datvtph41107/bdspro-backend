package emailsender

import (
	"appointment/config"

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
		Host:     config.EMAIL_HOST,
		Port:     config.EMAIL_PORT,
		Username: config.EMAIL_USERNAME,
		Password: config.EMAIL_PASSWORD,
		From:     config.EMAIL_FROM,
		Name:     config.EMAIL_NAME,
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
