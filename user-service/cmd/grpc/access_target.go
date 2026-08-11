package cmd_grpc

import (
	"errors"
	"fmt"
	"user/infra/handler/accessgrpc"
	accesspostgres "user/infra/postgres/access"
	useraccess "user/internal/access"

	_db "common/db"
	userpb "pb/types/user"

	"google.golang.org/grpc"
)

/**
 * registerAccessService lắp access flow mới vào gRPC server hiện tại.
 *
 * Hàm dùng cùng PostgreSQL connection của user-service và không mở HTTP route.
 */
func registerAccessService(server *grpc.Server) error {
	if server == nil {
		return errors.New("grpc server is nil")
	}
	if _db.DB == nil {
		return errors.New("user database is not initialized")
	}

	subscriptions := accesspostgres.NewSubscriptionStore(_db.DB)
	plans := accesspostgres.NewPlanAccessStore(_db.DB)
	service := useraccess.NewService(subscriptions, plans)
	handler := accessgrpc.NewHandler(service)

	userpb.RegisterInternalAccessServiceServer(server, handler)
	return nil
}

/**
 * internalAccessMethod trả true cho gRPC method chỉ dùng để đọc access.
 *
 * Legacy quota interceptor phải bỏ qua method này vì đây không phải user business action.
 */
func internalAccessMethod(fullMethod string) bool {
	return fullMethod == "/userpb.InternalAccessService/GetAccess"
}

func accessRegistrationError(err error) error {
	return fmt.Errorf("register internal access service: %w", err)
}
