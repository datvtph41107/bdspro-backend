package service

import (
	"common"
	"context"
	"fmt"
	"log"

	crmpb "pb/types/crm"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// DialCRMGRPC creates the single internal transport from Gateway to CRM.
// CRM no longer exposes public HTTP; all public/admin HTTP routes are owned by
// Gateway and translated to this gRPC connection.
func DialCRMGRPC(env string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	port := viper.GetString("service.port.crm")
	url := common.GetPrefixProtobufUrl(env, "crm")
	return grpc.NewClient(fmt.Sprintf("%s:%s", url, port), opts...)
}

func RegisterCrmService(ctx context.Context, mux *runtime.ServeMux, env string, opts []grpc.DialOption) *grpc.ClientConn {
	port := viper.GetString("service.port.crm")
	url := common.GetPrefixProtobufUrl(env, "crm")
	endpoint := fmt.Sprintf("%s:%s", url, port)

	conn, err := DialCRMGRPC(env, opts)
	if err != nil {
		log.Printf("[CRM] Không tạo được client tới %s; đăng ký grpc-gateway theo endpoint: %v", endpoint, err)
		registerCRMFromEndpoint(ctx, mux, endpoint, opts)
		return nil
	}

	log.Printf("[CRM] Đã tạo client tới %s (kết nối khi có request)", endpoint)
	registerCRMWithConn(ctx, mux, conn)
	return conn
}

func registerCRMWithConn(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) {
	_ = crmpb.RegisterPipelineServiceHandler(ctx, mux, conn)
	_ = crmpb.RegisterContactServiceHandler(ctx, mux, conn)
	_ = crmpb.RegisterLeadServiceHandler(ctx, mux, conn)
	_ = crmpb.RegisterStageServiceHandler(ctx, mux, conn)
	_ = crmpb.RegisterRuleServiceHandler(ctx, mux, conn)
	_ = crmpb.RegisterBlockServiceHandler(ctx, mux, conn)
	_ = crmpb.RegisterFollowServiceHandler(ctx, mux, conn)
	_ = crmpb.RegisterFriendGroupServiceHandler(ctx, mux, conn)
	_ = crmpb.RegisterFriendServiceHandler(ctx, mux, conn)
	_ = crmpb.RegisterInvitationInstallServiceHandler(ctx, mux, conn)
	_ = crmpb.RegisterSharingAccessServiceHandler(ctx, mux, conn)
	_ = crmpb.RegisterRateServiceHandler(ctx, mux, conn)
	_ = crmpb.RegisterReportServiceHandler(ctx, mux, conn)
	_ = crmpb.RegisterReportAdminServiceHandler(ctx, mux, conn)
	_ = crmpb.RegisterReportReasonServiceHandler(ctx, mux, conn)
	_ = crmpb.RegisterSupportTicketServiceHandler(ctx, mux, conn)
	_ = crmpb.RegisterRegistryServiceHandler(ctx, mux, conn)
	_ = crmpb.RegisterEnumServiceHandler(ctx, mux, conn)
	_ = crmpb.RegisterAppointmentServiceHandler(ctx, mux, conn)
	_ = crmpb.RegisterSeoDomainServiceHandler(ctx, mux, conn)
	_ = crmpb.RegisterSeoStrategyServiceHandler(ctx, mux, conn)
	_ = crmpb.RegisterCRMPublicContentServiceHandler(ctx, mux, conn)
	_ = crmpb.RegisterSeoMeasurementRuntimeServiceHandler(ctx, mux, conn)
}

func registerCRMFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) {
	_ = crmpb.RegisterPipelineServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterContactServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterLeadServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterStageServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterRuleServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterBlockServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterFollowServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterFriendGroupServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterFriendServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterInvitationInstallServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterSharingAccessServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterRateServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterReportServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterReportAdminServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterReportReasonServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterSupportTicketServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterRegistryServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterEnumServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterAppointmentServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterSeoDomainServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterSeoStrategyServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterCRMPublicContentServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	_ = crmpb.RegisterSeoMeasurementRuntimeServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
