package grpc

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/oklog/run"
	"github.com/tuandq2112/payment-service/appconfig"
	"github.com/tuandq2112/payment-service/internal/interface/handler"
	"github.com/tuandq2112/payment-service/internal/server/grpc/interceptors"
	paymentpb "pb/types/payment"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	paymentHandler *handler.PaymentHandler

	logger                         *flogging.FabricLogger
}

func NewGRPCServer(
	paymentHandler *handler.PaymentHandler,


) *GRPCServer {
	return &GRPCServer{
		paymentHandler: paymentHandler,

		logger: flogging.MustGetLogger(appconfig.SERVICE_NAME),
	}
}

func (s *GRPCServer) Start() {
	g := &run.Group{}

	start, stop := buildGRPCServer(s.logger,
		s.paymentHandler,
	)

	g.Add(start, stop)

	signalStart, signalStop := buildSignalHandler(s.logger)
	g.Add(signalStart, signalStop)

	if err := g.Run(); err != nil {
		s.logger.Errorf("Server exited: %v", err)
	}
}

func buildGRPCServer(logger *flogging.FabricLogger,
	paymentHandler *handler.PaymentHandler,
) (func() error, func(error)) {
	unaryInts := []grpc.UnaryServerInterceptor{
		interceptors.UnaryLoggerInterceptor(logger),
		interceptors.UnaryAuthInterceptor(logger, appconfig.Whitelist),
		interceptors.UnaryRecoveryInterceptor(logger),
	}
	streamInts := []grpc.StreamServerInterceptor{}
	logger.Info(fmt.Sprintf("Starting GRPC server at port %s", appconfig.Port))
	_, start, stop, err := BuildGRPCServer(
		fmt.Sprintf("%s:%s", appconfig.Host, appconfig.Port),
		func(s *grpc.Server) {
			logger.Info("Registering PaymentServiceServer")
			paymentpb.RegisterPaymentServiceServer(s, paymentHandler)
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
		server.GracefulStop()
	}

	return server, start, stop, nil
}
