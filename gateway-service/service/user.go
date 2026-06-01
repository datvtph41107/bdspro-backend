package service

import (
	"common"
	"context"
	"fmt"

	userpb "pb/types/user"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func RegisterUserService(ctx context.Context, mux *runtime.ServeMux, env string, opts []grpc.DialOption) {
	port := viper.GetString("service.port.user")
	url := common.GetPrefixProtobufUrl(env, "user")

	endpoint := fmt.Sprintf("%s:%s", url, port)
	_ = userpb.RegisterProfileServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	// _ = userpb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = userpb.RegisterReportServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = userpb.RegisterReportReasonServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = userpb.RegisterReportProofServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = userpb.RegisterBookmarkUserServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = userpb.RegisterAdminProfileServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = userpb.RegisterAdminUserProfileServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	// _ = userpb.RegisterSystemConfigServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = userpb.RegisterMembershipServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = userpb.RegisterPlanServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = userpb.RegisterMainAreaServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = userpb.RegisterPurposeUseServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = userpb.RegisterProfessionServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = userpb.RegisterTagServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = userpb.RegisterKYCServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = userpb.RegisterPriceTableServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
