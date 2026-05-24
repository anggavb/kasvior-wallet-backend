package pkg

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"os"
)

type SMTPMailer struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewSMTPMailerFromEnv() *SMTPMailer {
	return &SMTPMailer{
		host:     os.Getenv("SMTP_HOST"),
		port:     os.Getenv("SMTP_PORT"),
		username: os.Getenv("SMTP_USERNAME"),
		password: os.Getenv("SMTP_PASSWORD"),
		from:     os.Getenv("SMTP_FROM"),
	}
}

func (m *SMTPMailer) SendMail(ctx context.Context, to, subject, body string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if m.host == "" || m.port == "" || m.from == "" {
		return errors.New("smtp host, port, and from are required")
	}

	addr := net.JoinHostPort(m.host, m.port)
	var auth smtp.Auth
	if m.username != "" || m.password != "" {
		auth = smtp.PlainAuth("", m.username, m.password, m.host)
	}

	var message []byte
	message = fmt.Appendf(
		message,
		"From: %s\nTo: %s\nSubject: %s\nMIME-Version: 1.0\nContent-Type: text/plain; charset=UTF-8\n\n%s",
		m.from,
		to,
		subject,
		body,
	)

	return smtp.SendMail(addr, auth, m.from, []string{to}, message)
}
