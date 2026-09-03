package server

import (
	chatrpc "chat/infrastructure/rpc"
	_redis "common/redis"
	_utils "common/utils"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"chat/config"
	"chat/infrastructure/client"
	"chat/infrastructure/database/postgres"
	"chat/infrastructure/handlers"
	"chat/infrastructure/mapper"
	middlewares "chat/infrastructure/middleware"
	"chat/internal/repository"
	"chat/internal/usecases"

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

	db, err := postgres.NewPostgresDB(config.AppProperties.Database.DSN)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
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
	converstationUsecase := usecases.NewConversationUsecases(pRepository, userClient)
	redisClient := _redis.NewRedisService()
	syncProvider := _utils.NewSyncUtil(redisClient)
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
	mux := runtime.NewServeMux()

	delivery.RegisterHTTPServer(ctx, mux)

	handler := middlewares.JWTMiddleware(mux)
	handler = middlewares.LoggingMiddleware(handler)
	handler = middlewares.PresenterMiddleware(handler)
	handler = middlewares.EnableCORS(handler)

	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        handler,
		IdleTimeout:    parseTimeout(config.AppProperties.Timeout.Idle),
		ReadTimeout:    parseTimeout(config.AppProperties.Timeout.Read),
		WriteTimeout:   parseTimeout(config.AppProperties.Timeout.Write),
		MaxHeaderBytes: 1 << 20,
	}

	return &HTTPServer{server: server, cleanup: cleanup}, nil
}

func (s *HTTPServer) Start() error {
	if s.cleanup != nil {
		defer s.cleanup()
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Infof("HTTP server is running on %s", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("HTTP server error: %v", err)
		}
	}()

	<-stop
	logger.Info("Shutting down HTTP server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		logger.Errorf("HTTP server shutdown error: %v", err)
	}
	return nil
}
