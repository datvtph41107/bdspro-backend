package service

import (
	"common"
	"context"
	"fmt"
	chatpb "pb/types/chat"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func RegisterChatService(ctx context.Context, mux *runtime.ServeMux, env string, opts []grpc.DialOption) {
	port := viper.GetString("service.port.chat")
	url := common.GetPrefixProtobufUrl(env, "chat")

	endpoint := fmt.Sprintf("%s:%s", url, port)
	fmt.Println("url", url, "endpoint", endpoint)
	_ = chatpb.RegisterChatServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = chatpb.RegisterBackgroundImageServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
