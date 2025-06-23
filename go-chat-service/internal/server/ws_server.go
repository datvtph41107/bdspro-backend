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
	"chat/infrastructure/database/postgres"
	"chat/internal/constants"
	middlewares "chat/internal/delivery/middleware"
	"chat/internal/delivery/wshandler"
	"chat/internal/repository"
	"chat/internal/usecases"

	"github.com/hyperledger/fabric/common/flogging"
)

var wsLogger = flogging.MustGetLogger("ws_server")

type WSServer struct {
	server *http.Server
}

func NewWSServer(port string) *WSServer {
	wsLogger.Info("Initializing WebSocket server...")

	configs.InitConfig()
	db, err := postgres.NewPostgresDB(configs.DNS)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	repository := repository.NewRepository(db)

	conversationUsecase := usecases.NewConversationUsecases(repository)
	wsHandler := wshandler.NewWebSocketHandler(conversationUsecase)

	mux := http.NewServeMux()

	handler := setupHandlers(mux, wsHandler)

	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        handler,
		IdleTimeout:    configs.IDLE_TIMEOUT,
		ReadTimeout:    configs.READ_TIMEOUT,
		WriteTimeout:   configs.WRITE_TIMEOUT,
		MaxHeaderBytes: 1 << 20,
	}
	go wsHandler.SubscribeToRedisChannel(constants.API_WS_CHANNEL)
	return &WSServer{
		server: server,
	}
}

func setupHandlers(mux *http.ServeMux, wsHandler *wshandler.WebSocketHandler) http.Handler {
	handler := middlewares.JWTMiddleware(mux)
	handler = middlewares.WSLoggingMiddleware(handler)

	mux.HandleFunc("/ws", wsHandler.HandleWebSocket)
	mux.HandleFunc("/ws/room", wsHandler.HandleRoomWebSocket)

	return handler
}

func (s *WSServer) Start() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		wsLogger.Infof("WebSocket server is running on %s", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			wsLogger.Fatalf("WebSocket server error: %v", err)
		}
	}()

	<-stop
	wsLogger.Info("Shutting down WebSocket server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		wsLogger.Errorf("WebSocket server shutdown error: %v", err)
	}
}
