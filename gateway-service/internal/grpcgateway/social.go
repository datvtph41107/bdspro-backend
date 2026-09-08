package grpcgateway

import (
	"context"
	pb_social "pb/types/social"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func RegisterSocialService(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return registerAll(ctx, mux, conn,
		pb_social.RegisterNewsFeedServiceHandler,
		pb_social.RegisterCommentServiceHandler,
		pb_social.RegisterReportServiceHandler,
		pb_social.RegisterLikeServiceHandler,
	)
}
