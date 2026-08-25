package cmd_grpc

import (
	"errors"
	"fmt"
	"user/infra/handler/cataloggrpc"
	"user/infra/postgres"
	"user/infra/postgres/iampermission"
	useraccess "user/internal/entitlement/evaluate"
	"user/internal/entitlement/grpc"
	accesspostgres "user/internal/entitlement/postgres"
	"user/internal/modules/catalog/admin"
	"user/internal/modules/catalog/publication"

	_db "common/db"
	userpb "pb/types/user"

	"google.golang.org/grpc"
)

/**
 * registerTargetServices lắp các capability theo trusted identity vào gRPC server.
 *
 * Hàm dùng cùng PostgreSQL connection của user-service và không mở HTTP route.
 */
func registerTargetServices(server *grpc.Server) error {
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
	registerAdminCatalogService(server)
	return nil
}

/** registerAdminCatalogService installs the versioned catalog control plane. */
func registerAdminCatalogService(server *grpc.Server) {
	repository := postgres.NewCatalogRepo(_db.NewTransactionRepo(_db.DB))
	publisher := publication.NewService(
		_db.NewTransactionRepo(_db.DB),
		repository,
	)
	authorizer := iampermission.NewChecker(_db.DB)
	handler := cataloggrpc.NewHandler(admin.NewService(repository, publisher), authorizer)
	userpb.RegisterAdminCatalogServiceServer(server, handler)
}

/**
 * trustedIdentityMethod trả true cho RPC cần canonical identity đã được ký.
 *
 * Legacy quota interceptor phải bỏ qua các method này: chúng không phải user business action.
 */
func trustedIdentityMethod(fullMethod string) bool {
	switch fullMethod {
	case "/userpb.InternalAccessService/GetAccess",
		"/userpb.AdminCatalogService/ListPlanVersions",
		"/userpb.AdminCatalogService/GetPlanVersion",
		"/userpb.AdminCatalogService/PublishPlanVersion":
		return true
	default:
		return false
	}
}

func targetRegistrationError(err error) error {
	return fmt.Errorf("register trusted target services: %w", err)
}
