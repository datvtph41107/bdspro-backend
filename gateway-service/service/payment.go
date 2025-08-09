package service

import (
	"context"
	"fmt"
	"log"
	pb_payment "pb/types/payment"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func RegisterPaymentService(ctx context.Context, mux *runtime.ServeMux, env string, opts []grpc.DialOption) {
	port := viper.GetString("service.port.payment")
	url := "localhost"
	if env == "release" {
		url = "payment"
	}
	endpoint := fmt.Sprintf("%s:%s", url, port)
	err := pb_payment.RegisterPaymentServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = pb_payment.RegisterBankServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		log.Fatalf("Failed to register payment service: %v", err)
	}
}
