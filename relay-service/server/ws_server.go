package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"relay/client"
	"relay/config"
	"relay/constants"
	middlewares "relay/middleware"
	"relay/wshandler"

	"github.com/hyperledger/fabric/common/flogging"
)

// parseTimeout chuyển đổi string timeout thành time.Duration
func parseTimeout(timeoutStr string) time.Duration {
	duration, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return 30 * time.Second // Default 30s
	}
	return duration
}

var wsLogger = flogging.MustGetLogger("ws_server")

type WSServer struct {
	server *http.Server
}

func NewWSServer(port string) *WSServer {
	wsLogger.Info("Initializing WebSocket server...")

	chatClient := client.NewChatClient()
	userClient := client.NewUserClient()
	notificationClient := client.NewNotificationClient()
	wsHandler := wshandler.NewWebSocketHandler(chatClient, userClient, notificationClient)

	mux := http.NewServeMux()

	handler := setupHandlers(mux, wsHandler)

	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        handler,
		IdleTimeout:    parseTimeout(config.AppProperties.Timeout.Idle),
		ReadTimeout:    parseTimeout(config.AppProperties.Timeout.Read),
		WriteTimeout:   parseTimeout(config.AppProperties.Timeout.Write),
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

	// mux.HandleFunc("/ws", wsHandler.HandleWebSocket)
	// mux.HandleFunc("/ws/room", wsHandler.HandleRoomWebSocket)
	mux.HandleFunc("/connect", wsHandler.HandleConnect)

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
