// payment-completed-consumer is an independently deployable integration-event
// worker. Its availability never gates Payment or the other consumer service.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	paymentconsumer "notification/internal/integrations/paymentcompleted"
)

func main() {
	if err := run(); err != nil {
		log.Printf("payment completed consumer stopped: %v", err)
		os.Exit(1)
	}
}

func run() error {
	dsn := os.Getenv("NOTIFICATION_DATABASE_DSN")
	rabbitURL := os.Getenv("QHPRO_RABBITMQ_URL")
	exchange := getenv("QHPRO_BUSINESS_EVENTS_EXCHANGE", "qhpro.business-events")
	if dsn == "" || rabbitURL == "" {
		return errors.New("NOTIFICATION_DATABASE_DSN and QHPRO_RABBITMQ_URL are required")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("open notification database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("access notification sql db: %w", err)
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
