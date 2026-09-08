package server

import (
	"chat/config"
	"chat/infra/db"
	"chat/infra/handler"
	"chat/infra/postgres"
	chatrpc "chat/infra/rpc"
	"chat/internal/usecases"
	_db "common/db"
	_redis "common/redis"
	_utils "common/utils"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	middlewares "chat/infra/middleware"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hyperledger/fabric/common/flogging"
)

func parseTimeout(timeoutStr string) time.Duration {
	duration, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return 30 * time.Second // Default 30s
	}
	return duration
}

var logger = flogging.MustGetLogger("http_server")

type HTTPServer struct {
	server  *http.Server
	cleanup func()
}

func NewHTTPServer(port string) (*HTTPServer, error) {
	logger.Info("Initializing HTTP server...")
	ctx := context.Background()

	config.LoadConfig()

	// Initialize database using shared/common/db
	dbInstance, err := _db.NewDB()
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
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

	// Setup HTTP mux for gRPC gateway
	mux := runtime.NewServeMux()
	// Register HTTP handlers
	err = handler.RegisterHTTPServer(ctx, mux, chatHandler, backgroundHandler)
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("register chat HTTP handlers: %w", err)
	}

	// Apply middlewares
	httpHandler := middlewares.JWTMiddleware(mux)
	httpHandler = middlewares.LoggingMiddleware(httpHandler)
	httpHandler = middlewares.PresenterMiddleware(httpHandler)
	httpHandler = middlewares.EnableCORS(httpHandler)

	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        httpHandler,
		IdleTimeout:    parseTimeout(config.AppProperties.Timeout.Idle),
		ReadTimeout:    parseTimeout(config.AppProperties.Timeout.Read),
		WriteTimeout:   parseTimeout(config.AppProperties.Timeout.Write),
		MaxHeaderBytes: 1 << 20,
	}

	return &HTTPServer{
		server:  server,
		cleanup: cleanup,
	}, nil
}

func (s *HTTPServer) Start() error {
	if s.cleanup != nil {
		defer s.cleanup()
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	errCh := make(chan error, 1)

	go func() {
		logger.Infof("HTTP server is running on %s", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-stop:
		logger.Info("Shutting down HTTP server...")
	case err := <-errCh:
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	logger.Info("HTTP server stopped gracefully")
	return nil
}
