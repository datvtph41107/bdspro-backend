package server

import (
	_middleware "common/middleware"
	_redis "common/redis"
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
	"chat/infrastructure/impl"
	"chat/infrastructure/mapper"
	"chat/internal/repository"
	"chat/internal/usecases"

	"github.com/hyperledger/fabric/common/flogging"
	"google.golang.org/grpc"
)

var grpcLogger = flogging.MustGetLogger("grpc_server")

type GRPCServer struct {
	server *grpc.Server
	port   string
}

func NewGRPCServer(port string) *GRPCServer {
	grpcLogger.Info("Initializing gRPC server...")

	config.LoadConfig()

	db, err := postgres.NewPostgresDB(config.AppProperties.Database.DSN)
	if err != nil {
		grpcLogger.Fatalf("Failed to connect to database: %v", err)
	}
	pRepository := repository.NewRepository(db)
	userClient := client.NewUserClient()
	bdsproClient := client.NewBdsproClient()
	redisClient := _redis.NewRedisService()
	syncProvider := impl.NewSyncProvider(redisClient)
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
		bdsproClient,
		conversationMapper,
		messageMapper,
		backgroundMapper,
		syncProvider,
		pRepository,
	)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(_middleware.ParseGrpcMetadataContextMiddleware),
	)
	delivery.RegisterGRPCServer(server)

	return &GRPCServer{
		server: server,
		port:   port,
	}
}

func (s *GRPCServer) Start() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		grpcLogger.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		grpcLogger.Infof("gRPC server is running on port %s", s.port)
		if err := s.server.Serve(lis); err != nil {
			grpcLogger.Fatalf("Failed to serve: %v", err)
		}
	}()

	<-stop
	grpcLogger.Info("Shutting down gRPC server...")

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
		grpcLogger.Warning("Shutdown timed out, forcing stop")
		s.server.Stop()
	case <-done:
		grpcLogger.Info("Server stopped gracefully")
	}
}
