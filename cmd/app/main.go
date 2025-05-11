package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/anatoly_dev/go-mail-service/configs"
	"github.com/anatoly_dev/go-mail-service/internal/delivery/kafka"
	"github.com/anatoly_dev/go-mail-service/internal/repository"
	"github.com/anatoly_dev/go-mail-service/internal/usecase"
	"github.com/anatoly_dev/go-mail-service/pkg/logger"
)

type Application struct {
	logger       *logger.Logger
	config       *configs.Config
	consumer     *kafka.Consumer
	emailRepo    repository.EmailRepository
	emailUseCase *usecase.EmailUseCase
}

func NewApplication() *Application {
	return &Application{}
}

func (a *Application) initializeLogger() error {
	if err := logger.InitGlobalLogger("info", false); err != nil {
		return err
	}
	a.logger = logger.GetLogger()
	a.logger.Info("Logger initialized successfully")
	return nil
}

func (a *Application) loadConfig() {
	a.config = configs.NewConfig()
	a.logger.Info("Configuration loaded successfully")
}

func (a *Application) initializeComponents() error {
	a.logger.Info("Initializing application components...")

	a.emailRepo = repository.NewEmailRepository(a.logger, &a.config.SMTP)
	a.logger.Info("Email repository initialized")

	a.emailUseCase = usecase.NewEmailUseCase(a.emailRepo, a.logger)
	a.logger.Info("Email usecase initialized")

	var err error
	a.consumer, err = kafka.NewConsumer(a.config.Kafka.Brokers, a.emailUseCase, a.logger)
	if err != nil {
		return err
	}
	a.logger.Info("Kafka consumer initialized")

	return nil
}

func (a *Application) startConsumer(ctx context.Context) error {
	a.logger.Info("Starting Kafka consumer...")
	if err := a.consumer.Start(ctx, a.config.Kafka.Topic); err != nil {
		return err
	}
	a.logger.Info("Kafka consumer started successfully")
	return nil
}

func (a *Application) waitForShutdown() {
	a.logger.Info("Waiting for shutdown signal...")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	a.logger.Info("Received shutdown signal")
}

func (a *Application) cleanup() {
	a.logger.Info("Performing cleanup...")
	if a.consumer != nil {
		a.consumer.Close()
	}
	a.logger.Sync()
	a.logger.Info("Cleanup completed")
}

func (a *Application) Run() error {
	if err := a.initializeLogger(); err != nil {
		return err
	}
	defer a.cleanup()

	a.loadConfig()

	if err := a.initializeComponents(); err != nil {
		a.logger.WithError(err).Fatal("Failed to initialize components")
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := a.startConsumer(ctx); err != nil {
		a.logger.WithError(err).Fatal("Failed to start consumer")
		return err
	}

	a.logger.Info("Service started successfully")
	a.waitForShutdown()

	return nil
}

func main() {
	app := NewApplication()
	if err := app.Run(); err != nil {
		os.Exit(1)
	}
}
