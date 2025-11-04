package clients

import (
	"common"
	"context"
	"fmt"
	"log"
	hubpb "pb/types/hub"
	"time"

	_middleware "common/middleware"

	"google.golang.org/grpc"
)

type HubGrpcClient struct {
	Client         hubpb.LocationServiceClient
	InternalClient hubpb.HubInternalServiceClient
}

func MakeHubGrpcClient(env_runtime string) *HubGrpcClient {
	hubClient := &HubGrpcClient{Client: nil, InternalClient: nil}
	go func() {
		name := "hub"
		url := common.GetPrefixProtobufUrl(env_runtime, name)

		conn, err := grpc.Dial(fmt.Sprintf("%s:8280", url),
			grpc.WithInsecure(),
			grpc.WithBackoffMaxDelay(5*time.Second),
			grpc.WithUnaryInterceptor(_middleware.UnaryClientInterceptor),
		)
		if err != nil {
			log.Fatalf("failed to dial hub service: %v", err)
		}

		ctx, _ := context.WithCancel(context.Background())

		// Monitor connection status
		go _middleware.MonitorConnection(ctx, name, conn)

		hubClient.Client = hubpb.NewLocationServiceClient(conn)
		hubClient.InternalClient = hubpb.NewHubInternalServiceClient(conn)
	}()
	return hubClient
}
