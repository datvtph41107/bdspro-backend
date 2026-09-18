package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os/signal"
	"syscall"

	common_db "common/db"
	_middleware "common/middleware"
	process "common/process"
	sharedredis "common/redis"
	"notification/config"
	notificationdb "notification/db"
	brokerrabbit "notification/infra/broker/rabbitmq"
	"notification/infra/client"
	"notification/infra/firebase"
	postgres_eventing "notification/infra/postgres/eventing"
	notificationrpc "notification/infra/rpc"
	deliveryworker "notification/infra/worker/delivery"
	deliveryusecase "notification/internal/usecase/delivery"
	paymentcompleted "notification/internal/usecase/eventing/paymentcompleted"
	"notification/wire"
	notificationpb "pb/types/notification"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var GrpcCmd = &cobra.Command{
	Use: "grpc", Short: "Khởi chạy gRPC server",
	RunE: func(cmd *cobra.Command, args []string) error { return startGRPCServer() },
}

func startGRPCServer() error {
	cfg, err := config.LoadRuntimeConfig()
	if err != nil {
		return fmt.Errorf("load notification config: %w", err)
	}
	database, err := common_db.Open(common_db.DatabaseConfig{DSN: cfg.DatabaseDSN})
	if err != nil {
		return err
	}
	sqlDB, err := database.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()
	schemaPolicy, err := common_db.LoadSchemaPolicy("notification")
	if err != nil {
		return fmt.Errorf("load Notification schema policy: %w", err)
	}
	slog.Info(
		"notification database schema policy loaded",
		slog.Any("database.schema.mode", schemaPolicy.Mode),
		slog.Any("database.schema.source", schemaPolicy.Source),
	)
	if err := common_db.ApplySchemaPolicy(database, schemaPolicy, notificationdb.AutoMigrate); err != nil {
		return fmt.Errorf("apply Notification schema policy: %w", err)
	}
	redisService, err := sharedredis.Open(sharedredis.Config{Address: cfg.RedisAddress, Password: cfg.RedisPassword, DB: cfg.RedisDB})
	if err != nil {
		return err
	}
	defer redisService.Close()
	userRPC, closeUser, err := notificationrpc.OpenUserRPCClient(cfg.UserGRPCTarget)
	if err != nil {
		return err
	}
	defer closeUser()
	firebaseConfig, err := config.NewFirebaseConfig()
	if err != nil {
		return fmt.Errorf("load Firebase configuration: %w", err)
	}

	app, cleanup, err := wire.InitializeApp(database, redisService, userRPC, firebaseConfig)
	if err != nil {
		return fmt.Errorf("initialize notification app: %w", err)
	}
	defer cleanup()

	var paymentConsumer *brokerrabbit.PaymentCompletedSupervisor
	if cfg.RabbitURL != "" {
		paymentService := paymentcompleted.NewService(postgres_eventing.NewPaymentCompletedStore(database))
		paymentConsumer, err = brokerrabbit.NewPaymentCompletedSupervisor(
			cfg.RabbitURL,
			"notification-service",
			cfg.BusinessEventsExchange,
			paymentService,
		)
		if err != nil {
			return fmt.Errorf("create payment completed supervisor: %w", err)
		}
	} else {
		slog.Info(
			"notification payment consumer disabled",
			slog.String("component", "payment-consumer"),
			slog.String("reason", "rabbitmq-not-configured"),
		)
	}

	var pushDeliveryWorker *deliveryworker.Worker
	if firebaseConfig != nil && firebaseConfig.FirebaseApp() != nil {
		push := firebase.NewFirebaseService(firebaseConfig, database)
		tokens := client.NewAuthClient(userRPC)
		deliveryService := deliveryusecase.NewService(
			postgres_eventing.NewDeliveryStore(database),
			tokens,
			push,
			cfg.DeliveryLease,
		)
		pushDeliveryWorker = deliveryworker.New(cfg.DeliveryWorkerID, deliveryService, cfg.DeliveryPollInterval)
	} else {
		slog.Info(
			"notification delivery worker disabled",
			slog.String("component", "delivery-worker"),
			slog.String("reason", "firebase-not-configured"),
		)
	}

	port := fmt.Sprintf(":%d", cfg.ServerPort)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("listen notification gRPC on %s: %w", port, err)
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(_middleware.ParseGrpcMetadataContextMiddleware),
		grpc.StreamInterceptor(_middleware.ParseGrpcMetadataContextStreamMiddleware),
	)
	notificationpb.RegisterGatewayNotificationServiceServer(s, app.NotificationHandler)
	notificationpb.RegisterInternalNotificationServiceServer(s, app.InternalHandler)
	notificationpb.RegisterHistoryServiceServer(s, app.HistoryHandler)
	notificationpb.RegisterDealHistoryServiceServer(s, app.DealHistoryHandler)
	notificationpb.RegisterPropertyHistoryServiceServer(s, app.PropertyHistoryHandler)

	signalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	ctx, cancel := context.WithCancel(signalCtx)
	defer cancel()
	actorErrors := make(chan error, 2)
	actors := make([]process.Actor, 0, 2)
	if paymentConsumer != nil {
		actors = append(actors, process.ActorFunc(func(ctx context.Context) {
			if err := paymentConsumer.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
				select {
				case actorErrors <- fmt.Errorf("payment event consumer: %w", err):
				default:
				}
				cancel()
			}
		}))
	}
	if pushDeliveryWorker != nil {
		actors = append(actors, process.ActorFunc(func(ctx context.Context) {
			if err := pushDeliveryWorker.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
				select {
				case actorErrors <- fmt.Errorf("push delivery worker: %w", err):
				default:
				}
				cancel()
			}
		}))
	}
	slog.Info(
		"notification gRPC listening",
		slog.String("server.address", port),
	)
	if err := process.Run(ctx, cancel, s, lis, actors, process.Config{GracefulStopTimeout: cfg.ShutdownTimeout}); err != nil {
		return fmt.Errorf("notification gRPC lifecycle: %w", err)
	}
	select {
	case err := <-actorErrors:
		return err
	default:
	}
	return nil
}
