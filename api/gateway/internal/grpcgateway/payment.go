package grpcgateway

import (
	"context"
	pb_payment "pb/types/payment"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func RegisterPaymentService(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return registerAll(ctx, mux, conn,
		pb_payment.RegisterPaymentServiceHandler,
		pb_payment.RegisterBankServiceHandler,
		pb_payment.RegisterPaymentProfileServiceHandler,
		pb_payment.RegisterAdminCommerceServiceHandler,
	)
}
