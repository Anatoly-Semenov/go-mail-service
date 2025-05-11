package repository

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/anatoly_dev/go-mail-service/internal/domain"
	"github.com/anatoly_dev/go-mail-service/pkg/logger"
)

type emailRepository struct {
	logger *logger.Logger
	config *SMTPConfig
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func NewEmailRepository(logger *logger.Logger, config *SMTPConfig) EmailRepository {
	return &emailRepository{
		logger: logger,
		config: config,
	}
}

func (r *emailRepository) SendEmail(ctx context.Context, userEmail string, template domain.EmailTemplate) error {
	auth := smtp.PlainAuth("", r.config.Username, r.config.Password, r.config.Host)

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"+
		"%s",
		r.config.From, userEmail, template.Subject, template.Body)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", r.config.Host, r.config.Port),
		auth,
		r.config.From,
		[]string{userEmail},
		[]byte(msg),
	)

	if err != nil {
		r.logger.WithField("to", userEmail).WithError(err).Error("Failed to send email")
		return fmt.Errorf("failed to send email: %w", err)
	}

	r.logger.WithField("to", userEmail).WithField("subject", template.Subject).Info("Email sent successfully")
	return nil
}

func (r *emailRepository) GetUserEmail(ctx context.Context, userID string) (string, error) {
	return fmt.Sprintf("user_%s@example.com", userID), nil
}
