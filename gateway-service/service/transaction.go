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
	_ = pb_transaction.RegisterDealContractServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_transaction.RegisterDealCostServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_transaction.RegisterCostTypeServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	// _ = pb_transaction.RegisterDealTransactionServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
