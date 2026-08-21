package grpc

import (
	_middleware "common/middleware"
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	paymentpb "pb/types/payment"

	"payment/appconfig"
	"payment/internal/interface/handler"
	"payment/internal/job"
	"payment/internal/server/grpc/interceptors"
	"payment/internal/usecase"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/oklog/run"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	paymentHandler           *handler.PaymentHandler
	bankHandler              *handler.BankHandler
	internalHandler          *handler.InternalHandler
	paymentDashboardStatsJob *job.PaymentDashboardStatsJob
	TransactionDashboardJob  *job.TransactionDashboardJob
	notificationWorker       *usecase.NotificationWorker
	logger                   *flogging.FabricLogger
}

func NewGRPCServer(
	paymentHandler *handler.PaymentHandler,
	bankHandler *handler.BankHandler,
	internalHandler *handler.InternalHandler,
	paymentDashboardStatsJob *job.PaymentDashboardStatsJob,
	TransactionDashboardJob *job.TransactionDashboardJob,
	notificationWorker *usecase.NotificationWorker,
) *GRPCServer {
	return &GRPCServer{
		paymentHandler:           paymentHandler,
		bankHandler:              bankHandler,
		internalHandler:          internalHandler,
		paymentDashboardStatsJob: paymentDashboardStatsJob,
		TransactionDashboardJob:  TransactionDashboardJob,
		notificationWorker:       notificationWorker,
		logger:                   flogging.MustGetLogger(appconfig.SERVICE_NAME),
	}
}

func (s *GRPCServer) Start() {
	g := &run.Group{}

	start, stop := buildGRPCServer(s.logger,
		s.paymentHandler,
		s.bankHandler,
		s.internalHandler,
	)

	g.Add(start, stop)

	signalStart, signalStop := buildSignalHandler(s.logger)
	g.Add(signalStart, signalStop)

	// The process root owns cancellation for every long-lived actor. No job may
	// self-spawn from a constructor or detach via its own Background context.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g.Add(func() error {
		s.paymentDashboardStatsJob.RunDailyAtMidnight(ctx)
		return nil
	}, func(error) { cancel() })
	g.Add(func() error {
		return s.TransactionDashboardJob.Run(ctx)
	}, func(error) { cancel() })
	g.Add(func() error {
		return s.notificationWorker.Run(ctx)
	}, func(error) { cancel() })

	if err := g.Run(); err != nil {
		s.logger.Errorf("Server exited: %v", err)
	}
}

func buildGRPCServer(logger *flogging.FabricLogger,
	paymentHandler *handler.PaymentHandler,
	bankHandler *handler.BankHandler,
	internalHandler *handler.InternalHandler,
) (func() error, func(error)) {
	unaryInts := []grpc.UnaryServerInterceptor{
		_middleware.ParseGrpcMetadataContextMiddleware,
		interceptors.UnaryLoggerInterceptor(logger),
		// interceptors.UnaryAuthInterceptor(logger, appconfig.Whitelist),
		interceptors.UnaryRecoveryInterceptor(logger),
	}
	streamInts := []grpc.StreamServerInterceptor{
		_middleware.ParseGrpcMetadataContextStreamMiddleware,
	}
	logger.Info(fmt.Sprintf("Starting GRPC server at port %s", appconfig.Port))
	_, start, stop, err := BuildGRPCServer(
		fmt.Sprintf("%s:%s", appconfig.Host, appconfig.Port),
		func(s *grpc.Server) {
			logger.Info("Registering PaymentServiceServer")
			paymentpb.RegisterPaymentServiceServer(s, paymentHandler)
			paymentpb.RegisterBankServiceServer(s, bankHandler)
			paymentpb.RegisterInternalServiceServer(s, internalHandler)
		},
		unaryInts,
		streamInts,
	)
	if err != nil {
		return nil, nil
	}

	return start, stop
}

func buildSignalHandler(logger *flogging.FabricLogger) (func() error, func(error)) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	return func() error {
			sig := <-sigs
			return fmt.Errorf("signal received: %v", sig)
		}, func(err error) {
			signal.Stop(sigs)
			close(sigs)
		}
}

func BuildGRPCServer(
	addr string,
	registerFunc func(*grpc.Server),
	unaryInterceptors []grpc.UnaryServerInterceptor,
	streamInterceptors []grpc.StreamServerInterceptor,
) (
	server *grpc.Server,
	start func() error,
	stop func(err error),
	err error,
) {
	var opts []grpc.ServerOption

	if len(unaryInterceptors) > 0 {
		opts = append(opts, grpc.ChainUnaryInterceptor(unaryInterceptors...))
	}

	if len(streamInterceptors) > 0 {
		opts = append(opts, grpc.ChainStreamInterceptor(streamInterceptors...))
	}

	server = grpc.NewServer(opts...)
	registerFunc(server)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to listen: %w", err)
	}

	start = func() error {
		return server.Serve(lis)
	}

	stop = func(err error) {
		done := make(chan struct{})
		go func() {
			server.GracefulStop()
			close(done)
		}()
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			server.Stop()
			<-done
		}
	}

	return server, start, stop, nil
}
