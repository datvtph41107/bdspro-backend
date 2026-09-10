package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	common_db "common/db"
	"notification/config"
	brokerrabbit "notification/infra/broker/rabbitmq"
	postgres_eventing "notification/infra/postgres/eventing"
	paymentcompleted "notification/internal/usecase/eventing/paymentcompleted"
)

func main() {
	if err := run(); err != nil {
		log.Printf("payment completed consumer stopped: %v", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.LoadEventingConfig()
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	database, err := common_db.Open(common_db.DatabaseConfig{DSN: cfg.DatabaseDSN})
	if err != nil {
		return err
	}
	sqlDB, err := database.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()
	service := paymentcompleted.NewService(postgres_eventing.NewPaymentCompletedStore(database))
	consumer, err := brokerrabbit.NewPaymentCompletedSupervisor(
		cfg.RabbitURL,
		"notification-payment-completed",
		cfg.Exchange,
		service,
	)
	if err != nil {
		return err
	}
	if err := consumer.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}
