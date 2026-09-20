// payment-completed-consumer is an independently deployable integration-event
// worker. Its availability never gates Payment or the other consumer service.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"common/logging"
	paymentconsumer "crm/internal/integrations/paymentcompleted"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	os.Exit(runProcess())
}

func runProcess() int {
	closeLogger, err := logging.Configure("crm-payment-completed-consumer")
	if err != nil {
		fmt.Fprintf(os.Stderr, "configure CRM payment consumer logging: %v\n", err)
		return 1
	}
	defer func() { _ = closeLogger() }()

	if err := run(); err != nil {
		slog.Error("payment completed consumer stopped", slog.Any("error", err))
		return 1
	}
	return 0
}

func run() error {
	dsn := os.Getenv("CRM_DATABASE_DSN")
	rabbitURL := os.Getenv("QHPRO_RABBITMQ_URL")
	exchange := getenv("QHPRO_BUSINESS_EVENTS_EXCHANGE", "qhpro.business-events")
	if dsn == "" || rabbitURL == "" {
		return errors.New("CRM_DATABASE_DSN and QHPRO_RABBITMQ_URL are required")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("open crm database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("access crm sql db: %w", err)
	}
	defer sqlDB.Close()

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return fmt.Errorf("connect rabbitmq: %w", err)
	}
	defer conn.Close()

	consumer, err := paymentconsumer.NewConsumer(conn, exchange, paymentconsumer.NewStore(db))
	if err != nil {
		return err
	}
	defer consumer.Close()

	err = consumer.Run(ctx)
	if errors.Is(err, context.Canceled) {
		return nil
	}
	return err
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
