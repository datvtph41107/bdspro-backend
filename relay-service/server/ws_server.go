package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	relayrpc "relay/rpc"
	"syscall"
	"time"

	"relay/client"
	"relay/config"
	"relay/constants"
	middlewares "relay/middleware"
	relayredis "relay/redis"
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
	server     *http.Server
	handler    *wshandler.WebSocketHandler
	redisActor func(context.Context) error
	cleanup    func()
}

func NewWSServer(port string) (*WSServer, error) {
	wsLogger.Info("Initializing WebSocket server...")

	userRPC, userCleanup, err := relayrpc.NewUserRPCClient()
	if err != nil {
		return nil, fmt.Errorf("configure relay user RPC: %w", err)
	}
	chatRPC, chatCleanup, err := relayrpc.NewChatRPCClient()
	if err != nil {
		userCleanup()
		return nil, fmt.Errorf("configure relay chat RPC: %w", err)
	}
	notificationRPC, notificationCleanup, err := relayrpc.NewNotificationRPCClient()
	if err != nil {
		chatCleanup()
		userCleanup()
		return nil, fmt.Errorf("configure relay notification RPC: %w", err)
	}
	chatClient := client.NewChatClient(chatRPC)
	userClient := client.NewUserClient(userRPC)
	notificationClient := client.NewNotificationClient(notificationRPC)
	redisClient, err := relayredis.NewRedisClient()
	if err != nil {
		notificationCleanup()
		chatCleanup()
		userCleanup()
		return nil, err
	}
	wsHandler := wshandler.NewWebSocketHandler(chatClient, userClient, notificationClient, redisClient)
	cleanup := func() {
		_ = wsHandler.Close()
		_ = redisClient.Close()
		notificationCleanup()
		chatCleanup()
		userCleanup()
	}

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
	return &WSServer{
		server:  server,
		handler: wsHandler,
		redisActor: func(ctx context.Context) error {
			return wsHandler.SubscribeToRedisChannel(ctx, constants.API_WS_CHANNEL)
		},
		cleanup: cleanup,
	}, nil
}

func setupHandlers(mux *http.ServeMux, wsHandler *wshandler.WebSocketHandler) http.Handler {
	handler := middlewares.JWTMiddleware(mux)
	handler = middlewares.WSLoggingMiddleware(handler)

	// mux.HandleFunc("/ws", wsHandler.HandleWebSocket)
	// mux.HandleFunc("/ws/room", wsHandler.HandleRoomWebSocket)
	mux.HandleFunc("/connect", wsHandler.HandleConnect)

	return handler
}

func (s *WSServer) Start() error {
	if s.cleanup != nil {
		defer s.cleanup()
	}
	signalCtx, stopSignals := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignals()
	ctx, cancel := context.WithCancel(signalCtx)
	defer cancel()
	errCh := make(chan error, 2)

	go func() {
		wsLogger.Infof("WebSocket server is running on %s", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("serve Relay WebSocket: %w", err)
		}
	}()
	redisDone := make(chan struct{})
	go func() {
		defer close(redisDone)
		if s.redisActor != nil {
			if err := s.redisActor(ctx); err != nil {
				errCh <- err
			}
		}
	}()

	var runtimeErr error
	select {
	case <-ctx.Done():
	case runtimeErr = <-errCh:
		cancel()
	}
	wsLogger.Info("Shutting down WebSocket server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := s.server.Shutdown(shutdownCtx); err != nil && runtimeErr == nil {
		runtimeErr = fmt.Errorf("shutdown Relay WebSocket: %w", err)
	}
	select {
	case <-redisDone:
	case <-shutdownCtx.Done():
		if runtimeErr == nil {
			runtimeErr = fmt.Errorf("stop Relay Redis actor: %w", shutdownCtx.Err())
		}
	}
	return runtimeErr
}
