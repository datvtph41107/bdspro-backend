package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	common_db "common/db"
	"common/logging"
	"notification/config"
	"notification/infra/client"
	"notification/infra/firebase"
	postgres_eventing "notification/infra/postgres/eventing"
	notificationrpc "notification/infra/rpc"
	deliveryworker "notification/infra/worker/delivery"
	deliveryusecase "notification/internal/usecase/delivery"
)

func main() {
	closeLogger, err := logging.Configure("notification-service")
	if err != nil {
		fmt.Fprintf(os.Stderr, "configure Notification delivery-worker logging: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = closeLogger() }()

	if err := run(); err != nil {
		slog.Error(
			"notification delivery worker stopped",
			slog.String("component", "delivery-worker"),
			slog.Any("error", err),
		)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.LoadDeliveryConfig()
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
	userRPC, closeUser, err := notificationrpc.OpenUserRPCClient(cfg.UserGRPCTarget)
	if err != nil {
		return err
	}
	defer closeUser()
	firebaseConfig, err := config.NewFirebaseConfig()
	if err != nil {
		return err
	}
	if firebaseConfig == nil || firebaseConfig.FirebaseApp() == nil {
		return fmt.Errorf("notification delivery worker requires NOTIFICATION_FIREBASE_CREDENTIAL_FILE")
	}
	push := firebase.NewFirebaseService(firebaseConfig, database)
	tokens := client.NewAuthClient(userRPC)
	service := deliveryusecase.NewService(postgres_eventing.NewDeliveryStore(database), tokens, push, cfg.Lease)
	return deliveryworker.New(cfg.WorkerID, service, cfg.PollInterval).Run(ctx)
}
