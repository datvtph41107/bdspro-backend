// commerce-event-publisher is a process-owned worker for publishing durable
// Payment integration events. It intentionally runs separately from the gRPC
// server so broker health cannot become request-path availability.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	commercepg "payment/internal/commerce/adapters/postgres"
	commercerabbit "payment/internal/commerce/adapters/rabbitmq"
	"payment/internal/commerce/application"

	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := run(); err != nil {
		log.Printf("commerce event publisher stopped: %v", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Database and broker connections are process-owned resources. Constructors
	// below never spawn hidden goroutines and every opened resource has a Close.
	db, err := gorm.Open(postgres.Open(cfg.databaseDSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("open payment database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("access payment sql db: %w", err)
	}
	defer sqlDB.Close()

	conn, err := amqp.Dial(cfg.rabbitURL)
	if err != nil {
		return fmt.Errorf("connect rabbitmq: %w", err)
	}
	defer conn.Close()

	publisher, err := commercerabbit.NewPublisher(conn, cfg.exchange)
	if err != nil {
		return err
	}
	defer publisher.Close()

	store := commercepg.NewStore(db)
	svc := application.NewOutboxPublisherService(store, publisher, time.Now, cfg.retryDelay)
	ticker := time.NewTicker(cfg.pollInterval)
	defer ticker.Stop()

	// Drain one eligible responsibility at a time. PostgreSQL SKIP LOCKED and
	// claim_version fencing make multiple process replicas safe; throughput is
	// scaled by replicas rather than by hiding concurrency inside this loop.
	for {
		processed, err := svc.RunOne(ctx, cfg.workerID, cfg.lease)
		if err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("publish attempt failed; durable outbox remains owner: %v", err)
		}
		if ctx.Err() != nil {
			return nil
		}
		if processed {
			continue
		}
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
}

type runtimeConfig struct {
	databaseDSN  string
	rabbitURL    string
	exchange     string
	workerID     string
	pollInterval time.Duration
	lease        time.Duration
	retryDelay   time.Duration
}

func loadConfig() (runtimeConfig, error) {
	cfg := runtimeConfig{
		databaseDSN:  os.Getenv("PAYMENT_DATABASE_DSN"),
		rabbitURL:    os.Getenv("QHPRO_RABBITMQ_URL"),
		exchange:     getenv("QHPRO_BUSINESS_EVENTS_EXCHANGE", "qhpro.business-events"),
		workerID:     getenv("QHPRO_PAYMENT_EVENT_WORKER_ID", hostnameWorkerID()),
		pollInterval: durationEnv("QHPRO_EVENT_POLL_INTERVAL", time.Second),
		lease:        durationEnv("QHPRO_EVENT_CLAIM_LEASE", 30*time.Second),
		retryDelay:   durationEnv("QHPRO_EVENT_RETRY_DELAY", 5*time.Second),
	}
	if cfg.databaseDSN == "" || cfg.rabbitURL == "" {
		return runtimeConfig{}, errors.New("PAYMENT_DATABASE_DSN and QHPRO_RABBITMQ_URL are required")
	}
	if cfg.pollInterval <= 0 || cfg.lease <= 0 || cfg.retryDelay < 0 {
		return runtimeConfig{}, errors.New("invalid commerce event worker durations")
	}
	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
func hostnameWorkerID() string {
	h, _ := os.Hostname()
	if h == "" {
		h = "unknown"
	}
	return "payment-events-" + h + "-" + strconv.Itoa(os.Getpid())
}
func durationEnv(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return -1
	}
	return d
}
