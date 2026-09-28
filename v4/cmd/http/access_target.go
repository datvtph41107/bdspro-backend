package main

import (
	"strings"

	authpb "pb/types/auth"
)

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
		"/userpb.CommercialProfileService/CancelMySubscription",
		"/userpb.CheckoutService/Checkout",
		"/userpb.CheckoutService/ListPaymentMethods",
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
