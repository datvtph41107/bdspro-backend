package clients

import (
	"common"
	"context"
	"fmt"
	"log"
	crmpb "pb/types/crm"
	"time"

	_middleware "common/middleware"

	"google.golang.org/grpc"
)

type FeedbackGrpcClient struct {
	crmpb.CrmInternalServiceClient
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
		feedbackClient.CrmInternalServiceClient = crmpb.NewCrmInternalServiceClient(conn)
	}()
	return feedbackClient
}
