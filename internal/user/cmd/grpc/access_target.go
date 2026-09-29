package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
	paymentclient "user/infra/client/payment"

	"user/infra/handler/grpc/entitlement"
	"user/infra/postgres"
	checkoutpostgres "user/infra/postgres/checkout"
	accesspostgres "user/infra/postgres/entitlement"
	"user/infra/postgres/iampermission"
	organizationpostgres "user/infra/postgres/organization"
	subscriptionpostgres "user/infra/postgres/subscription"
	"user/internal/usecase/commercialprofile"
	useraccess "user/internal/usecase/entitlement"
	organizationusecase "user/internal/usecase/organization"
	"user/internal/usecase/plan/admin"
	"user/internal/usecase/plan/publish"
	subscriptionusecase "user/internal/usecase/subscription"
	checkoutapp "user/internal/usecase/subscription/checkout"
	settlementapp "user/internal/usecase/subscription/settlement"

	_db "common/db"
	"gorm.io/gorm"
	authpb "pb/types/auth"
	userpb "pb/types/user"
	usergrpc "user/infra/handler/grpc"

	"google.golang.org/grpc"
)

/**
 * registerTargetServices lắp các capability theo trusted identity vào gRPC server.
 *
 * Hàm dùng cùng PostgreSQL connection của user-service và không mở HTTP route.
 */
func registerTargetServices(server *grpc.Server, database *gorm.DB) (func(), error) {
	if server == nil {
		return nil, errors.New("grpc server is nil")
	}
	if database == nil {
		return nil, errors.New("user database is not initialized")
	}
	subscriptions := accesspostgres.NewSubscriptionStore(database)
	plans := accesspostgres.NewPlanAccessStore(database)
	service := useraccess.NewService(subscriptions, plans)
	handler := accessgrpc.NewHandler(service)

	userpb.RegisterInternalAccessServiceServer(server, handler)
	registerAdminPlanService(server, database)
	registerAdminSubscriptionService(server, database)
	registerCommercialProfileService(server, database)
	registerAdminOrganizationService(server, database)
	registerInternalOrganizationMembershipService(server, database)

	tx := _db.NewTransactionRepo(database)
	checkoutStore := checkoutpostgres.NewCheckoutStore(tx)
	paymentCommands, cleanup, err := paymentclient.NewCommands()
	if err != nil {
		return nil, fmt.Errorf("create Payment commercial channel: %w", err)
	}
	checkoutService := checkoutapp.NewService(checkoutStore, checkoutStore, checkoutStore, paymentCommands, paymentCommands, checkoutStore, time.Now)
	organizationAuthorizer := organizationpostgres.NewCheckoutAuthorizer(database)
	userpb.RegisterCheckoutServiceServer(server, usergrpc.NewSubscriptionCheckoutHandler(checkoutService, organizationAuthorizer))

	settlementStore := subscriptionpostgres.NewSettlementStore(tx)
	userpb.RegisterInternalSubscriptionServiceServer(server, usergrpc.NewSubscriptionSettlementHandler(settlementapp.NewService(settlementStore)))
	return cleanup, nil
}

func registerAdminSubscriptionService(server *grpc.Server, database *gorm.DB) {
	store := subscriptionpostgres.NewAdminProjectionStore(database)
	authorizer := iampermission.NewChecker(database)
	handler := usergrpc.NewSubscriptionAdminHandler(subscriptionusecase.NewAdminService(store), authorizer)
	userpb.RegisterAdminCommercialServiceServer(server, handler)
}

func registerCommercialProfileService(server *grpc.Server, database *gorm.DB) {
	tx := _db.NewTransactionRepo(database)
	planRepository := postgres.NewPlanVersionRepo(tx)
	planService := admin.NewService(planRepository, nil)
	subscriptionStore := subscriptionpostgres.NewAdminProjectionStore(database)
	subscriptionService := subscriptionusecase.NewAdminService(subscriptionStore)
	handler := usergrpc.NewCommercialProfileHandler(
		commercialprofile.NewService(planService, subscriptionService, time.Now),
	)
	userpb.RegisterCommercialProfileServiceServer(server, handler)
}

// registerAdminOrganizationService installs the User-owned organization
// directory control plane on the existing process database and lifecycle.
func registerAdminOrganizationService(server *grpc.Server, database *gorm.DB) {
	repository := organizationpostgres.NewRepository(database)
	authorizer := iampermission.NewChecker(database)
	handler := usergrpc.NewAdminOrganizationHandler(organizationusecase.NewService(repository), authorizer)
	userpb.RegisterAdminOrganizationServiceServer(server, handler)
}

func registerInternalOrganizationMembershipService(server *grpc.Server, database *gorm.DB) {
	repository := organizationpostgres.NewRepository(database)
	handler := usergrpc.NewInternalOrganizationMembershipHandler(organizationusecase.NewService(repository))
	userpb.RegisterInternalOrganizationMembershipServiceServer(server, handler)
}

/** registerAdminPlanService installs the versioned Plan control plane. */
func registerAdminPlanService(server *grpc.Server, database *gorm.DB) {
	tx := _db.NewTransactionRepo(database)
	repository := postgres.NewPlanVersionRepo(tx)
	publisher := publish.NewService(
		tx,
		repository,
	)
	authorizer := iampermission.NewChecker(database)
	handler := usergrpc.NewPlanVersionAdminHandler(admin.NewService(repository, publisher), authorizer)
	userpb.RegisterAdminCatalogServiceServer(server, handler)
}

/**
 * trustedIdentityMethod trả true cho RPC cần canonical identity đã được ký.
 *
 * Legacy quota interceptor phải bỏ qua các method này: chúng không phải user business action.
 */
func trustedIdentityMethod(fullMethod string) bool {
	if strings.HasPrefix(fullMethod, "/userpb.AdminUserProfileService/") ||
		strings.HasPrefix(fullMethod, "/userpb.AdminProfileService/") ||
		strings.HasPrefix(fullMethod, "/userpb.AdminOrganizationService/") ||
		strings.HasPrefix(fullMethod, "/userpb.InternalOrganizationMembershipService/") {
		return true
	}
	switch fullMethod {
	case authpb.AuthService_ListAdmins_FullMethodName,
		"/userpb.InternalAccessService/GetAccess",
		"/userpb.CommercialProfileService/ListAvailablePlans",
		"/userpb.CommercialProfileService/GetMyCommercialProfile",
		"/userpb.CheckoutService/Checkout",
		"/userpb.CheckoutService/CreatePaymentAttempt",
		"/userpb.InternalSubscriptionService/ApplySettlement",
		"/userpb.AdminCatalogService/CreatePlanVersionDraft",
		"/userpb.AdminCatalogService/ListPlanVersions",
		"/userpb.AdminCatalogService/GetPlanVersion",
		"/userpb.AdminCatalogService/UpdatePlanVersionDraft",
		"/userpb.AdminCatalogService/ValidatePlanVersionDraft",
		"/userpb.AdminCatalogService/DeletePlanVersionDraft",
		"/userpb.AdminCatalogService/PublishPlanVersion",
		"/userpb.AdminCatalogService/RetirePlanVersion",
		"/userpb.AdminCommercialService/ListSubscriptions",
		"/userpb.AdminCommercialService/GetUserProjection",
		"/userpb.AdminCommercialService/ListUserProjections":
		return true
	default:
		return false
	}
}

func targetRegistrationError(err error) error {
	return fmt.Errorf("register trusted target services: %w", err)
}
