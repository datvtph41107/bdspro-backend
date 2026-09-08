package emailsender

import (
	"crm/config"

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
		Host:     config.AppProperties.Email.Host,
		Port:     config.AppProperties.Email.Port,
		Username: config.AppProperties.Email.Username,
		Password: config.AppProperties.Email.Password,
		From:     config.AppProperties.Email.From,
		Name:     config.AppProperties.Email.Name,
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