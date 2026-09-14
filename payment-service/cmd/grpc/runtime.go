package grpcserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	common_db "common/db"
	_middleware "common/middleware"
	qhprorpc "common/rpc"
	"common/rpcenv"
	"payment/config"
	rabbit "payment/infra/broker/rabbitmq"
	authclient "payment/infra/client/auth"
	notificationclient "payment/infra/client/notification"
	commercegrpc "payment/infra/handler/grpc"
	bankhandler "payment/infra/handler/grpc/bank"
	commercialprofilegrpc "payment/infra/handler/grpc/commercialprofile"
	wallethandler "payment/infra/handler/grpc/wallet"
	paymentinfra "payment/infra/postgres"
	bankpg "payment/infra/postgres/bank"
	paymentpg "payment/infra/postgres/payment"
	walletpg "payment/infra/postgres/wallet"
	"payment/infra/provider/sepay"
	"payment/infra/reference"
	fulfillmentworker "payment/infra/worker/fulfillment"
	"payment/internal/job"
	"payment/internal/server/grpc/interceptors"
	adminusecase "payment/internal/usecase/admin"
	"payment/internal/usecase/attempt"
	bankuc "payment/internal/usecase/bank"
	commercialprofileusecase "payment/internal/usecase/commercialprofile"
	"payment/internal/usecase/fulfillment"
	"payment/internal/usecase/order"
	"payment/internal/usecase/outbox"
	"payment/internal/usecase/settlement"
	walletuc "payment/internal/usecase/wallet"
	paymentworker "payment/worker"
	paymentpb "pb/types/payment"

	"github.com/hyperledger/fabric/common/flogging"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

var GrpcServerCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Start the Payment gRPC process",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(cmd.Context())
	},
}

// run is the only Payment gRPC process composition root. It owns typed config,
// one PostgreSQL pool, outbound connections, long-lived actors, gRPC lifecycle
// and close order. Business packages receive capabilities; they do not open
// process resources themselves.
func run(parent context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load Payment config: %w", err)
	}

	db, closeDB, err := paymentinfra.Open(cfg.Database)
	if err != nil {
		return err
	}
	defer closeDB()
	schemaPolicy, err := common_db.LoadSchemaPolicy("payment")
	if err != nil {
		return fmt.Errorf("load Payment schema policy: %w", err)
	}
	logger := flogging.MustGetLogger("payment-service")
	logger.Infof("Payment database schema mode=%s source=%s", schemaPolicy.Mode, schemaPolicy.Source)
	if err := common_db.ApplySchemaPolicy(db, schemaPolicy, paymentinfra.AutoMigrate); err != nil {
		return fmt.Errorf("apply Payment schema policy: %w", err)
	}

	// Legacy Wallet/Bank remains a live public capability during the commercial
	// cutover, but it shares the exact same process-owned DB lifecycle.
	paymentMethodRepo := walletpg.NewPaymentMethodRepository(db)
	walletRepo := walletpg.NewWalletsRepository(db)
	walletTxRepo := walletpg.NewWalletTransactionRepository(db)
	walletAuditRepo := walletpg.NewWalletAuditLogRepository(db)
	withdrawalRepo := walletpg.NewWithdrawalRequestsRepository(db)
	dashboardRepo := walletpg.NewDashboardMetricPostgresRepository(db)

	notificationRPC, closeNotification, err := notificationclient.ProvideNotificationRPCClientAt(cfg.Notification.Address)
	if err != nil {
		return fmt.Errorf("open Notification client: %w", err)
	}
	defer closeNotification()
	notificationClient := notificationclient.NewNotificationClient(notificationRPC)
	notificationWorker := walletuc.NewNotificationWorker(notificationClient)

	paymentUsecase := walletuc.NewPaymentUsecase(
		paymentMethodRepo,
		walletRepo,
		walletTxRepo,
		walletAuditRepo,
		withdrawalRepo,
		notificationWorker,
		dashboardRepo,
	)
	paymentDashboardJob := job.NewPaymentDashboardStatsJob(paymentUsecase, db)
	transactionDashboardJob := job.NewTransactionDashboardJob(paymentUsecase)

	bankRepo := bankpg.NewBankPostgresRepository(db)
	bankUsecase := bankuc.NewBankUsecase(bankRepo)
	bankHandler := bankhandler.NewBankHandler(bankUsecase)
	internalHandler := bankhandler.NewInternalHandler(bankUsecase)

	// Canonical durable commercial state machine uses the same process DB.
	store := paymentpg.NewStore(db)
	if strings.TrimSpace(cfg.Rabbit.URL) == "" {
		return errors.New("PAYMENT_RABBIT_URL is required")
	}
	outboxFactory := func() (*outbox.Service, func(), error) {
		rabbitConnection, err := amqp.Dial(cfg.Rabbit.URL)
		if err != nil {
			return nil, nil, fmt.Errorf("connect Payment RabbitMQ: %w", err)
		}
		rabbitPublisher, err := rabbit.NewPublisher(rabbitConnection, cfg.Rabbit.Exchange)
		if err != nil {
			_ = rabbitConnection.Close()
			return nil, nil, fmt.Errorf("create Payment event publisher: %w", err)
		}
		closeResources := func() {
			_ = rabbitPublisher.Close()
			_ = rabbitConnection.Close()
		}
		return outbox.NewService(store, rabbitPublisher, time.Now, cfg.Rabbit.PublisherRetry), closeResources, nil
	}
	outboxSupervisor := paymentworker.NewOutboxSupervisor(
		outboxFactory,
		paymentworker.ProcessID("payment-outbox"),
		cfg.Rabbit.PublisherLease,
		cfg.Rabbit.PublisherPoll,
		cfg.Rabbit.PublisherRetry,
	)
	now := time.Now
	orderService := order.NewService(store, reference.New(), now, cfg.Commerce.OrderTTL)
	provider := sepay.New("sepay")
	attemptService := attempt.NewService(store, store, []attempt.Provider{provider}, now)
	settlementService := settlement.NewService(store, store, now)
	userClient, closeUser, err := newUserClient(cfg.User.Address, cfg.User.Timeout)
	if err != nil {
		return err
	}
	defer closeUser()
	authClient, closeAuth, err := authclient.NewRPCClientAt(cfg.Auth.Address)
	if err != nil {
		return fmt.Errorf("open Auth client: %w", err)
	}
	defer closeAuth()
	fulfillmentService := fulfillment.NewService(store, userClient, now, cfg.Commerce.FulfillmentRetry)
	fulfillmentWorker := fulfillmentworker.New(
		fulfillmentService,
		store,
		workerID("payment-fulfillment"),
		cfg.Commerce.FulfillmentLease,
		cfg.Commerce.FulfillmentPoll,
		now,
	)

	paymentHandler := wallethandler.NewPaymentHandler(paymentUsecase, paymentDashboardJob, settlementService, cfg.Provider.SepayAPIKey)
	commerceHandler := commercegrpc.NewCommerceHandler(orderService, attemptService)
	commercialProfileHandler := commercialprofilegrpc.New(commercialprofileusecase.NewService(store))
	adminCommerceHandler := commercegrpc.NewAdminCommerceHandler(adminusecase.NewService(store, settlementService, now), authClient)

	grpcServer, listener, err := buildServer(
		cfg.Server.Address,
		logger,
		paymentHandler,
		bankHandler,
		internalHandler,
		commerceHandler,
		commercialProfileHandler,
		adminCommerceHandler,
	)
	if err != nil {
		return err
	}
	defer listener.Close()

	ctx, stop := signal.NotifyContext(parent, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	group, actorCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		if err := grpcServer.Serve(listener); err != nil && actorCtx.Err() == nil {
			return fmt.Errorf("serve Payment gRPC: %w", err)
		}
		return nil
	})
	group.Go(func() error { return paymentDashboardJob.Run(actorCtx) })
	group.Go(func() error { return transactionDashboardJob.Run(actorCtx) })
	group.Go(func() error { return notificationWorker.Run(actorCtx) })
	group.Go(func() error {
		fulfillmentWorker.Run(actorCtx)
		return nil
	})
	group.Go(func() error { return outboxSupervisor.Run(actorCtx) })

	<-actorCtx.Done()
	gracefulStop(grpcServer, cfg.Server.ShutdownTimeout)
	return group.Wait()
}

func buildServer(
	address string,
	logger *flogging.FabricLogger,
	paymentHandler paymentpb.PaymentServiceServer,
	bankHandler paymentpb.BankServiceServer,
	internalHandler paymentpb.InternalServiceServer,
	commerceHandler paymentpb.InternalCommerceServiceServer,
	commercialProfileHandler paymentpb.PaymentProfileServiceServer,
	adminCommerceHandler paymentpb.AdminCommerceServiceServer,
) (*grpc.Server, net.Listener, error) {
	transport := qhprorpc.ServerTransport(rpcenv.LoadTransportConfig())
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			transport.Unary,
			_middleware.ParseGrpcMetadataContextMiddleware,
			interceptors.UnaryLoggerInterceptor(logger),
			interceptors.UnaryRecoveryInterceptor(logger),
		),
		grpc.ChainStreamInterceptor(
			transport.Stream,
			_middleware.ParseGrpcMetadataContextStreamMiddleware,
		),
	)
	paymentpb.RegisterPaymentServiceServer(server, paymentHandler)
	paymentpb.RegisterBankServiceServer(server, bankHandler)
	paymentpb.RegisterInternalServiceServer(server, internalHandler)
	paymentpb.RegisterInternalCommerceServiceServer(server, commerceHandler)
	paymentpb.RegisterPaymentProfileServiceServer(server, commercialProfileHandler)
	paymentpb.RegisterAdminCommerceServiceServer(server, adminCommerceHandler)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, nil, fmt.Errorf("listen Payment gRPC %s: %w", address, err)
	}
	return server, listener, nil
}

func gracefulStop(server *grpc.Server, timeout time.Duration) {
	if server == nil {
		return
	}
	done := make(chan struct{})
	go func() {
		server.GracefulStop()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(timeout):
		server.Stop()
		<-done
	}
}

func workerID(prefix string) string {
	host, _ := os.Hostname()
	if host == "" {
		host = "unknown"
	}
	return prefix + "-" + host + "-" + strconv.Itoa(os.Getpid())
}
