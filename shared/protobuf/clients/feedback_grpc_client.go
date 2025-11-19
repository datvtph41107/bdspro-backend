package clients

import (
	"common"
	"context"
	"fmt"
	"log"
	feedbackpb "pb/types/feedback"
	sharepb "pb/types/shared"
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
		name := "feedback"
		url := common.GetPrefixProtobufUrl(env_runtime, name)

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

func (c *FeedbackGrpcClient) GetRateStats(ctx context.Context, req *sharepb.IdRequest) (*sharepb.RateStatsV3Proto, error) {
	resp, err := c.InternalServiceClient.GetRateStats(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
