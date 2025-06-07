package service

import (
	"context"
	"fmt"

	pb_payment "pb/types/payment"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func RegisterPaymentService(ctx context.Context, mux *runtime.ServeMux, url string, opts []grpc.DialOption) {
	port := viper.GetString("service.port.payment")
	endpoint := fmt.Sprintf("%s:%s", url, port)
	_ = pb_payment.RegisterPaymentServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
