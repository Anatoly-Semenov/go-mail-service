package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/anatoly_dev/go-mail-service/internal/domain"
	"github.com/anatoly_dev/go-mail-service/internal/usecase"
	"github.com/anatoly_dev/go-mail-service/pkg/logger"
)

type Consumer struct {
	consumer sarama.Consumer
	usecase  *usecase.EmailUseCase
	logger   *logger.Logger
}

func NewConsumer(brokers []string, usecase *usecase.EmailUseCase, logger *logger.Logger) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &Consumer{
		consumer: consumer,
		usecase:  usecase,
		logger:   logger,
	}, nil
}

func (c *Consumer) Start(ctx context.Context, topic string) error {
	partitionConsumer, err := c.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return fmt.Errorf("failed to create partition consumer: %w", err)
	}

	go func() {
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				var emailMsg domain.EmailMessage
				if err := json.Unmarshal(msg.Value, &emailMsg); err != nil {
					c.logger.WithError(err).WithField("value", string(msg.Value)).Error("Failed to unmarshal message")
					continue
				}

				if err := c.usecase.ProcessEmailMessage(ctx, emailMsg); err != nil {
					c.logger.WithError(err).
						WithField("user_id", emailMsg.UserID).
						WithField("email_type", string(emailMsg.EmailType)).
						Error("Failed to process email message")
					continue
				}

				c.logger.WithField("user_id", emailMsg.UserID).
					WithField("email_type", string(emailMsg.EmailType)).
					Info("Successfully processed email message")

			case err := <-partitionConsumer.Errors():
				c.logger.WithError(err).Error("Error consuming message")

			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
