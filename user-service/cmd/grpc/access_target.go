package main

/**
 * trustedIdentityMethod trả true cho RPC cần canonical identity đã được ký.
 *
 * Legacy quota interceptor phải bỏ qua các method này: chúng không phải user business action.
 */
func trustedIdentityMethod(fullMethod string) bool {
	switch fullMethod {
	case "/userpb.InternalAccessService/GetAccess",
		"/userpb.CommercialProfileService/ListAvailablePlans",
		"/userpb.CommercialProfileService/GetMyCommercialProfile",
		"/userpb.CheckoutService/Checkout",
		"/userpb.CheckoutService/CreatePaymentAttempt",
		"/userpb.InternalSubscriptionService/ApplySettlement",
		"/userpb.AdminCatalogService/ListPlanVersions",
		"/userpb.AdminCatalogService/GetPlanVersion",
		"/userpb.AdminCatalogService/PublishPlanVersion",
		"/userpb.AdminCommercialService/ListSubscriptions",
		"/userpb.AdminCommercialService/GetUserProjection",
		"/userpb.AdminCommercialService/ListUserProjections":
		return true
	default:
		return false
	}
}
