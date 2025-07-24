package grpc

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb_marketing "pb/types/marketing"

	"marketing/env"
	"marketing/infrastructure/server/grpc/interceptors"
	"marketing/internal/handler"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/oklog/run"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	marketingHandler *handler.MarketingGrpcHandler
	logger           *flogging.FabricLogger
}

func NewGRPCServer(marketingHandler *handler.MarketingGrpcHandler) *GRPCServer {
	return &GRPCServer{
		marketingHandler: marketingHandler,
		logger:           flogging.MustGetLogger(env.SERVICE_NAME),
	}
}

func (s *GRPCServer) Start() {
	g := &run.Group{}

	start, stop := buildGRPCServer(s.logger, s.marketingHandler)
	g.Add(start, stop)

	signalStart, signalStop := buildSignalHandler(s.logger)
	g.Add(signalStart, signalStop)

	if err := g.Run(); err != nil {
		s.logger.Errorf("Server exited: %v", err)
	}
}

func buildGRPCServer(logger *flogging.FabricLogger, marketingHandler *handler.MarketingGrpcHandler) (func() error, func(error)) {
	unaryInts := []grpc.UnaryServerInterceptor{
		interceptors.UnaryLoggerInterceptor(logger),
		interceptors.UnaryRecoveryInterceptor(logger),
	}
	streamInts := []grpc.StreamServerInterceptor{}
	logger.Info(fmt.Sprintf("Starting GRPC server at port %s", env.GRPC_PORT))
	_, start, stop, err := BuildGRPCServer(
		fmt.Sprintf(":%s", env.GRPC_PORT),
		func(s *grpc.Server) {
			logger.Info("Registering MarketingServiceServer")
			pb_marketing.RegisterMarketingServiceServer(s, marketingHandler)
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
		<-sigs
		logger.Info("Received shutdown signal")
		return nil
	}, func(err error) {
		logger.Errorf("Signal handler error: %v", err)
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
	server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	)

	registerFunc(server)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to listen: %v", err)
	}

	start = func() error {
		return server.Serve(lis)
	}

	stop = func(err error) {
		server.GracefulStop()
	}

	return server, start, stop, nil
} 