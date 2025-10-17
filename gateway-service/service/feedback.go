package service

import (
	"context"
	"fmt"

	feedbackpb "pb/types/feedback"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func RegisterFeedbackService(ctx context.Context, mux *runtime.ServeMux, env string, opts []grpc.DialOption) {
	port := viper.GetString("service.port.feedback")
	url := "localhost"
	if env == "release" {
		url = "feedback_grpc"
	}
	endpoint := fmt.Sprintf("%s:%s", url, port)
	_ = feedbackpb.RegisterRateServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = feedbackpb.RegisterReportServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	// _ = feedbackpb.RegisterReportAdminServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = feedbackpb.RegisterReportReasonServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = feedbackpb.RegisterRegistryServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = feedbackpb.RegisterInternalServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
