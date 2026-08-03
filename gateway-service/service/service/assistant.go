package service

import (
	"common"
	"context"
	"fmt"

	assistantpb "pb/types/assistant"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func RegisterAssistantService(ctx context.Context, mux *runtime.ServeMux, env string, opts []grpc.DialOption) {
	port := viper.GetString("service.port.assistant")
	if port == "" {
		port = "8218"
	}

	url := common.GetPrefixProtobufUrl(env, "assistant")

	endpoint := fmt.Sprintf("%s:%s", url, port)

	// Register Assistant Service
	_ = assistantpb.RegisterAssistantServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
