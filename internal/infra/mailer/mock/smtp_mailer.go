package mock

import (
	"log"
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
	log.Printf("[DummyMailer] to=%s subject=%s", to, subject)
	return nil
}
