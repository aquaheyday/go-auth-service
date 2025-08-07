package mailer

import (
	"fmt"
	"net/smtp"
)

type SMTPMailer struct {
	host string
	port int
	user string
	pass string
}

func NewSMTPMailer(host string, port int, user, pass string) *SMTPMailer {
	return &SMTPMailer{host: host, port: port, user: user, pass: pass}
}

func (m *SMTPMailer) Send(to, subject, body string) error {
	auth := smtp.PlainAuth("", m.user, m.pass, m.host)
	msg := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))
	addr := fmt.Sprintf("%s:%d", m.host, m.port)
	return smtp.SendMail(addr, auth, m.user, []string{to}, msg)
}
