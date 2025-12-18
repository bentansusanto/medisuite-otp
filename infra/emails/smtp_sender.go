package emails

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"

	// "medisuite/pkg/logs"
	"net/smtp"
	"strings"

	"medisuite-api/config"
)

type SMTPSender struct{}

func NewSMTPSender() *SMTPSender {
	return &SMTPSender{}
}

func (s *SMTPSender) Send(email config.Email) error {
	body := "From: " + config.CONFIG_SMTP_USER + "\n" +
		"To: " + strings.Join(email.To, ",") + "\n" +
		"Cc: " + strings.Join(email.Cc, ",") + "\n" +
		"Subject: " + email.Subject + "\n\n" +
		email.Body

	auth := smtp.PlainAuth("", config.CONFIG_SMTP_USER, config.CONFIG_SMTP_PASS, config.CONFIG_SMTP_HOST)
	smtpAddr := fmt.Sprintf("%s:%s", config.CONFIG_SMTP_HOST, config.CONFIG_SMTP_PORT)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         config.CONFIG_SMTP_HOST,
	}

	conn, err := tls.Dial("tcp", smtpAddr, tlsConfig)
	if err != nil {
		slog.Error("error", "failed to connect to SMTP server", err)
		return errors.New("failed to connect to SMTP server")
	}

	client, err := smtp.NewClient(conn, config.CONFIG_SMTP_HOST)
	if err != nil {
		slog.Error("error", "failed to create SMTP client", err)
		return errors.New("failed to create SMTP client")
	}
	defer client.Close()

	if error := client.Auth(auth); error != nil {
		slog.Error("error", "failed to authenticate SMTP client", err)
		return errors.New("failed to authenticate SMTP client")
	}

	if error := client.Mail(config.CONFIG_SMTP_USER); error != nil {
		slog.Error("error", "failed to send SMTP mail", err)
		return errors.New("failed to send SMTP mail")
	}

	for _, addr := range append(email.To, email.Cc...) {
		if error := client.Rcpt(addr); error != nil {
			slog.Error("error", "failed to send SMTP mail", err)
			return errors.New("failed to send SMTP mail")
		}
	}

	w, error := client.Data()
	if error != nil {
		slog.Error("error", "failed to send SMTP mail", err)
		return errors.New("failed to send SMTP mail")
	}

	_, error = w.Write([]byte(body))
	if error != nil {
		slog.Error("error", "failed to send SMTP mail", err)
		return errors.New("failed to send SMTP mail")
	}

	if error := w.Close(); error != nil {
		slog.Error("error", "failed to send SMTP mail", err)
		return errors.New("failed to send SMTP mail")
	}

	return nil
}
