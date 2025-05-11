package usecase

import (
	"context"
	"fmt"

	"github.com/anatoly_dev/go-mail-service/internal/domain"
	"github.com/anatoly_dev/go-mail-service/internal/repository"
	"github.com/anatoly_dev/go-mail-service/pkg/logger"
)

type EmailUseCase struct {
	repo   repository.EmailRepository
	logger *logger.Logger
}

func NewEmailUseCase(repo repository.EmailRepository, logger *logger.Logger) *EmailUseCase {
	return &EmailUseCase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *EmailUseCase) ProcessEmailMessage(ctx context.Context, msg domain.EmailMessage) error {
	userEmail, err := uc.repo.GetUserEmail(ctx, msg.UserID)
	if err != nil {
		uc.logger.WithField("user_id", msg.UserID).WithError(err).Error("Failed to get user email")
		return fmt.Errorf("failed to get user email: %w", err)
	}

	template, err := uc.getEmailTemplate(msg.EmailType)
	if err != nil {
		uc.logger.WithField("email_type", msg.EmailType).WithError(err).Error("Failed to get email template")
		return fmt.Errorf("failed to get email template: %w", err)
	}

	if err := uc.repo.SendEmail(ctx, userEmail, template); err != nil {
		uc.logger.WithField("user_email", userEmail).WithError(err).Error("Failed to send email")
		return fmt.Errorf("failed to send email: %w", err)
	}

	uc.logger.WithField("user_email", userEmail).WithField("email_type", msg.EmailType).Info("Email processed successfully")
	return nil
}

func (uc *EmailUseCase) getEmailTemplate(emailType domain.EmailType) (domain.EmailTemplate, error) {
	switch emailType {
	case domain.RegistrationEmail:
		return domain.EmailTemplate{
			Subject: "Добро пожаловать!",
			Body:    `<h1>Добро пожаловать!</h1><p>Спасибо за регистрацию в нашем сервисе.</p>`,
		}, nil
	case domain.PaymentReminderEmail:
		return domain.EmailTemplate{
			Subject: "Напоминание об оплате",
			Body:    `<h1>Напоминание об оплате</h1><p>Ваша подписка истекает через 5 дней.</p>`,
		}, nil
	case domain.PaymentSuccessEmail:
		return domain.EmailTemplate{
			Subject: "Оплата успешно произведена",
			Body:    `<h1>Оплата успешно произведена</h1><p>Спасибо за продление подписки!</p>`,
		}, nil
	case domain.PasswordChangeEmail:
		return domain.EmailTemplate{
			Subject: "Смена пароля",
			Body:    `<h1>Смена пароля</h1><p>Ваш пароль был успешно изменен.</p>`,
		}, nil
	default:
		return domain.EmailTemplate{}, fmt.Errorf("unknown email type: %s", emailType)
	}
}
