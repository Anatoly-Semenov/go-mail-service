package configs

import (
	"os"
	"strconv"

	"github.com/anatoly_dev/go-mail-service/internal/repository"
)

type Config struct {
	Kafka struct {
		Brokers []string
		Topic   string
	}
	SMTP repository.SMTPConfig
}

func NewConfig() *Config {
	cfg := &Config{
		Kafka: struct {
			Brokers []string
			Topic   string
		}{
			Brokers: []string{"localhost:9092"},
			Topic:   "mail.send-email",
		},
		SMTP: repository.SMTPConfig{
			Host:     "smtp.gmail.com",
			Port:     587,
			Username: "your-email@gmail.com",
			Password: "your-app-password",
			From:     "your-email@gmail.com",
		},
	}
	
	if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
		cfg.Kafka.Brokers = []string{brokers}
	}

	if host := os.Getenv("SMTP_HOST"); host != "" {
		cfg.SMTP.Host = host
	}

	if port := os.Getenv("SMTP_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.SMTP.Port = p
		}
	}

	if username := os.Getenv("SMTP_USERNAME"); username != "" {
		cfg.SMTP.Username = username
	}

	if password := os.Getenv("SMTP_PASSWORD"); password != "" {
		cfg.SMTP.Password = password
	}

	if from := os.Getenv("SMTP_FROM"); from != "" {
		cfg.SMTP.From = from
	}

	return cfg
}
