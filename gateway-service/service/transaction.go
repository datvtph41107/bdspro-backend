package service

import (
	"context"
	"fmt"
	pb_transaction "pb/types/transaction"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func RegisterTransactionService(ctx context.Context, mux *runtime.ServeMux, env string, opts []grpc.DialOption) {
	port := viper.GetString("service.port.transaction")
	url := "localhost"
	if env == "release" {
		url = "transaction"
	}
	endpoint := fmt.Sprintf("%s:%s", url, port)
	_ = pb_transaction.RegisterTransactionServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
