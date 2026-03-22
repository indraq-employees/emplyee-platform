package services

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

type MailerService struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewMailerService(host, port, username, password, from string) *MailerService {
	return &MailerService{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (m *MailerService) IsConfigured() bool {
	return m != nil && m.host != "" && m.port != "" && m.username != "" && m.password != "" && m.from != ""
}

func (m *MailerService) SendHTML(to, subject, htmlBody string) error {
	if !m.IsConfigured() {
		return fmt.Errorf("smtp is not configured")
	}

	message := strings.Join([]string{
		fmt.Sprintf("From: %s", m.from),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=UTF-8",
		"",
		htmlBody,
	}, "\r\n")

	addr := fmt.Sprintf("%s:%s", m.host, m.port)
	auth := smtp.PlainAuth("", m.username, m.password, m.host)

	if m.port == "465" {
		conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: m.host})
		if err != nil {
			return err
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, m.host)
		if err != nil {
			return err
		}
		defer client.Close()

		if err = client.Auth(auth); err != nil {
			return err
		}
		if err = client.Mail(m.from); err != nil {
			return err
		}
		if err = client.Rcpt(to); err != nil {
			return err
		}

		writer, err := client.Data()
		if err != nil {
			return err
		}
		if _, err = writer.Write([]byte(message)); err != nil {
			_ = writer.Close()
			return err
		}
		if err = writer.Close(); err != nil {
			return err
		}
		return client.Quit()
	}

	return smtp.SendMail(addr, auth, m.from, []string{to}, []byte(message))
}