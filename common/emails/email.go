package emails

import (
	"medisuite-api/config"
	"medisuite-api/infra/emails"
)

type EmailSender interface {
	Send(email config.Email) error
}

type Service struct {
	sender EmailSender
}

func NewService(sender EmailSender) *Service {
	return &Service{sender: sender}
}

// SendEmail sends an email using the configured email sender
func SendEmail(to []string, cc []string, subject string, body string) error {
	smtpSender := emails.NewSMTPSender()
	emailService := NewService(smtpSender)

	// Send the email using the configured email sender
	return emailService.sender.Send(config.Email{
		To:      to,
		Cc:      cc,
		Subject: subject,
		Body:    body,
	})
}
