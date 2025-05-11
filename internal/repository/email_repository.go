package repository

import (
	"context"

	"github.com/anatoly_dev/go-mail-service/internal/domain"
)

type EmailRepository interface {
	SendEmail(ctx context.Context, userID string, template domain.EmailTemplate) error
	GetUserEmail(ctx context.Context, userID string) (string, error)
}
