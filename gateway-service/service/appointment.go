package service

import (
	"context"
	"fmt"
	pb_appointment "pb/types/appointment"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func RegisterAppointmentService(ctx context.Context, mux *runtime.ServeMux, env string, opts []grpc.DialOption) {

	port := viper.GetString("service.port.appointment")
	url := "localhost"
	if env == "release" {
		url = "appointment"
	}
	endpoint := fmt.Sprintf("%s:%s", url, port)
	_ = pb_appointment.RegisterAppointmentServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
