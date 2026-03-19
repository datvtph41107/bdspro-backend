package server

import (
	_redis "common/redis"
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
	"chat/infrastructure/impl"
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
	server *http.Server
}

func NewHTTPServer(port string) *HTTPServer {
	logger.Info("Initializing HTTP server...")
	ctx := context.Background()

	config.LoadConfig()

	db, err := postgres.NewPostgresDB(config.AppProperties.Database.DSN)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	pRepository := repository.NewRepository(db)
	userClient := client.NewUserClient()
	bdsproClient := client.NewBdsproClient()
	converstationUsecase := usecases.NewConversationUsecases(pRepository, userClient)
	redisClient := _redis.NewRedisService()
	syncProvider := impl.NewSyncProvider(redisClient)
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
		bdsproClient,
		conversationMapper,
		messageMapper,
		backgroundMapper,
		syncProvider,
		pRepository,
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

	return &HTTPServer{server: server}
}

func (s *HTTPServer) Start() {
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
}
