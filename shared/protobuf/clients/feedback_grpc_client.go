package clients

import (
	"context"
	"fmt"
	"log"
	feedbackpb "pb/types/feedback"
	"time"

	_middleware "common/middleware"

	"google.golang.org/grpc"
)

type FeedbackGrpcClient struct {
	feedbackpb.InternalServiceClient
}

func MakeFeedbackGrpcClient(env_runtime string) *FeedbackGrpcClient {
	feedbackClient := &FeedbackGrpcClient{}
	go func() {
		url := "localhost"
		name := "feedback"
		if env_runtime == "release" {
			url = name
		}

		conn, err := grpc.Dial(fmt.Sprintf("%s:8217", url),
			grpc.WithInsecure(),
			grpc.WithBackoffMaxDelay(5*time.Second),
			grpc.WithUnaryInterceptor(_middleware.UnaryClientInterceptor),
		)
		if err != nil {
			log.Fatalf("failed to dial: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		// Khởi chạy goroutine để theo dõi trạng thái kết nối
		go _middleware.MonitorConnection(ctx, name, conn)

		// Initialize all service clients
		feedbackClient.InternalServiceClient = feedbackpb.NewInternalServiceClient(conn)
	}()
	return feedbackClient
}

func (c *FeedbackGrpcClient) GetRateStats(ctx context.Context, targetId uint64, targetType uint8) error {
	// _, err := c.InternalServiceClient.GetRateStats(ctx, &feedbackpb.GetRateStatsRequest{
	// 	Name: "user",
	// })
	return nil
}
