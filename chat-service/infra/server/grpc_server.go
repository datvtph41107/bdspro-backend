package server

import (
	"chat/config"
	"chat/infra/db"
	"chat/infra/handler"
	"chat/infra/postgres"
	chatrpc "chat/infra/rpc"
	"chat/internal/usecases"
	_db "common/db"
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

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	server  *grpc.Server
	port    string
	cleanup func()
	// wsHandler  *wshandler.WebSocketHandler
}

func NewGRPCServer(port string) (*GRPCServer, error) {
	logger := logging.WithComponent(context.Background(), "grpc_server")
	logger.Info("Initializing gRPC server...")

	config.LoadConfig()
	viper.Set("redis.pass", config.AppProperties.Redis.Password)

	// Initialize database using shared/common/db
	dbInstance, err := _db.NewDB()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to database: %v", err))
		os.Exit(1)
	}

	// Migrate domain
	db.MigrateDomain()

	// Create TransactionRepo
	transactionRepo := _db.NewTransactionRepo(dbInstance)

	// Initialize repositories
	conversationRepo := postgres.NewConversationRepo(transactionRepo)
	messageRepo := postgres.NewMessageRepo(transactionRepo)
	participantRepo := postgres.NewParticipantRepo(transactionRepo)
	readPositionRepo := postgres.NewReadPositionRepo(transactionRepo)
	messageReactionRepo := postgres.NewMessageReactionRepo(transactionRepo)
	backgroundRepo := postgres.NewBackgroundRepo(transactionRepo)

	// Process-owned outbound RPC dependencies.
	userRPC, userCleanup, err := chatrpc.NewUserRPCClient()
	if err != nil {
		return nil, fmt.Errorf("configure chat user RPC: %w", err)
	}

	bdsproRPC, bdsproCleanup, err := chatrpc.NewBdsproRPCClient()
	if err != nil {
		userCleanup()
		return nil, fmt.Errorf("configure chat bdspro RPC: %w", err)
	}

	rpcCleanup := func() {
		bdsproCleanup()
		userCleanup()
	}

	userClient := handler.NewUserClient(userRPC)
	bdsproClient := handler.NewBdsproClient(bdsproRPC)

	// Initialize usecases
	conversationUsecase := usecases.NewConversationUsecases(conversationRepo, participantRepo, userClient)

	messageUsecase := usecases.NewMessageUsecases(messageRepo, conversationRepo, participantRepo, conversationUsecase, transactionRepo)
	participantUsecase := usecases.NewParticipantUsecases(conversationRepo, participantRepo)
	readMarkUsecase := usecases.NewReadMarkUsecases(conversationRepo, participantRepo, readPositionRepo, messageRepo)
	messageReactionUsecase := usecases.NewMessageReactionUsecases(conversationRepo, participantRepo, messageRepo, messageReactionRepo)
	backgroundImageUsecase := usecases.NewBackgroundImageUsecases(backgroundRepo)

	// Initialize mappers
	backgroundMapper := handler.NewBackgroundMapper()
	conversationMapper := handler.NewConversationMapper(backgroundMapper)

	redisService := _redis.NewRedisService()
	syncProvider := _utils.NewSyncUtil(redisService)

	// Initialize handlers
	chatHandler := handler.NewChatHandler(
		conversationUsecase,
		messageUsecase,
		readMarkUsecase,
		participantUsecase,
		messageReactionUsecase,
		userClient,
		bdsproClient,
		conversationMapper,
		backgroundMapper,
		syncProvider,
	)
	chatHandler.Start()
	cleanup := func() {
		_ = chatHandler.Close()
		_ = redisService.Close()
		rpcCleanup()
	}

	backgroundHandler := handler.NewBackgroundHandler(backgroundImageUsecase)

	// Initialize gRPC server
	server := grpc.NewServer(
		grpc.UnaryInterceptor(_middleware.ParseGrpcMetadataContextMiddleware),
		grpc.StreamInterceptor(_middleware.ParseGrpcMetadataContextStreamMiddleware),
	)

	// Register handlers
	handler.RegisterGRPCServer(server, chatHandler, backgroundHandler)

	return &GRPCServer{
		server:  server,
		port:    port,
		cleanup: cleanup,
	}, nil
}

func (s *GRPCServer) Start() error {
	logger := logging.WithComponent(context.Background(), "grpc_server")
	if s.cleanup != nil {
		defer s.cleanup()
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	errCh := make(chan error, 1)

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

	select {
	case <-stop:
		logger.Info("Shutting down servers...")
	case err := <-errCh:
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Shutdown gRPC server
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
		logger.Info("Servers stopped gracefully")
	}

	return nil
}
