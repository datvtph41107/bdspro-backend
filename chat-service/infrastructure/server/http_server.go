package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"chat/configs"
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

var logger = flogging.MustGetLogger("http_server")

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(port string) *HTTPServer {
	logger.Info("Initializing HTTP server...")
	ctx := context.Background()

	configs.InitConfig()

	db, err := postgres.NewPostgresDB(configs.DNS)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	pRepository := repository.NewRepository(db)
	userClient := client.NewUserClient()
	converstationUsecase := usecases.NewConversationUsecases(pRepository, userClient)
	messageUsecase := usecases.NewMessageUsecases(pRepository)
	participantUsecase := usecases.NewParticipantUsecases(pRepository)
	readReceptUsecase := usecases.NewReadReceptUsecases(pRepository)
	messageReactionUsecase := usecases.NewMessageReactionUsecases(pRepository)
	backgroundImageUsecase := usecases.NewBackgroundImageUsecases(pRepository)
	conversationMapper := mapper.NewConversationMapper()

	delivery := handlers.NewDelivery(
		converstationUsecase,
		messageUsecase,
		readReceptUsecase,
		participantUsecase,
		messageReactionUsecase,
		backgroundImageUsecase,
		userClient,
		conversationMapper,
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
		IdleTimeout:    configs.IDLE_TIMEOUT,
		ReadTimeout:    configs.READ_TIMEOUT,
		WriteTimeout:   configs.WRITE_TIMEOUT,
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
