package server

import (
	chatrpc "chat/infrastructure/rpc"
	"common/logging"
	_middleware "common/middleware"
	_redis "common/redis"
	_utils "common/utils"
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"chat/config"
	"chat/infrastructure/client"
	"chat/infrastructure/database/postgres"
	"chat/infrastructure/handlers"
	"chat/infrastructure/mapper"
	"chat/internal/repository"
	"chat/internal/usecases"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	server  *grpc.Server
	port    string
	cleanup func()
}

func NewGRPCServer(port string) (*GRPCServer, error) {
	logger := logging.WithComponent(context.Background(), "grpc_server")
	logger.Info("Initializing gRPC server...")

	config.LoadConfig()

	db, err := postgres.NewPostgresDB(config.AppProperties.Database.DSN)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to database: %v", err))
		os.Exit(1)
	}
	pRepository := repository.NewRepository(db)
	userRPC, userCleanup, err := chatrpc.NewUserRPCClient()
	if err != nil {
		return nil, fmt.Errorf("configure chat-v1 user RPC: %w", err)
	}
	bdsproRPC, bdsproCleanup, err := chatrpc.NewBdsproRPCClient()
	if err != nil {
		userCleanup()
		return nil, fmt.Errorf("configure chat-v1 bdspro RPC: %w", err)
	}
	crmRPC, crmCleanup, err := chatrpc.NewCrmRPCClient()
	if err != nil {
		bdsproCleanup()
		userCleanup()
		return nil, fmt.Errorf("configure chat-v1 crm RPC: %w", err)
	}
	cleanup := func() { crmCleanup(); bdsproCleanup(); userCleanup() }
	userClient := client.NewUserClient(userRPC)
	crmClient := client.NewCrmClient(crmRPC)
	bdsproClient := client.NewBdsproClient(bdsproRPC)
	redisClient := _redis.NewRedisService()
	syncProvider := _utils.NewSyncUtil(redisClient)
	// conversationMapper := mapper.NewConversationMapper()

	converstationUsecase := usecases.NewConversationUsecases(pRepository, userClient)
	messageUsecase := usecases.NewMessageUsecases(pRepository, converstationUsecase, syncProvider)
	participantUsecase := usecases.NewParticipantUsecases(pRepository)
	readReceptUsecase := usecases.NewReadReceptUsecases(pRepository)
	messageReactionUsecase := usecases.NewMessageReactionUsecases(pRepository)
	backgroundImageUsecase := usecases.NewBackgroundImageUsecases(pRepository)
	backgroundMapper := mapper.NewBackgroundMapper()
	messageMapper := mapper.NewMessageMapper()
	conversationMapper := mapper.NewConversationMapper(backgroundMapper, messageMapper)

	delivery := handlers.NewDelivery(
		converstationUsecase,
		messageUsecase,
		readReceptUsecase,
		participantUsecase,
		messageReactionUsecase,
		backgroundImageUsecase,
		userClient,
		crmClient,
		bdsproClient,
		conversationMapper,
		messageMapper,
		backgroundMapper,
		syncProvider,
		pRepository,
		usecases.NewSystemMessageUsecase(pRepository, userClient),
	)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(_middleware.ParseGrpcMetadataContextMiddleware),
		grpc.StreamInterceptor(_middleware.ParseGrpcMetadataContextStreamMiddleware),
	)
	delivery.RegisterGRPCServer(server)

	return &GRPCServer{server: server, port: port, cleanup: cleanup}, nil
}

func (s *GRPCServer) Start() error {
	logger := logging.WithComponent(context.Background(), "grpc_server")
	if s.cleanup != nil {
		defer s.cleanup()
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to listen: %v", err))
		os.Exit(1)
	}

	go func() {
		logger.Info(fmt.Sprintf("gRPC server is running on port %s", s.port))
		if err := s.server.Serve(lis); err != nil {
			logger.Error(fmt.Sprintf("Failed to serve: %v", err))
			os.Exit(1)
		}
	}()

	<-stop
	logger.Info("Shutting down gRPC server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	done := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(done)
	}()

	select {
	case <-ctx.Done():
		logger.Warn("Shutdown timed out, forcing stop")
		s.server.Stop()
	case <-done:
		logger.Info("Server stopped gracefully")
	}
	return nil
}
