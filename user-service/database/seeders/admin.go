package seeders

type adminSeedSpec struct {
	Username             string
	PasswordHash         string
	FullName             string
	Email                string
	Phone                string
	RoleKey              string
	RoleName             string
	RoleDescription      string
	AdminRoleDescription string
	RoleType             string
	RoleCode             int
	PermissionKeys       []string
	GrantAllPermissions  bool
}

var acceptanceAdminPermissionKeys = []string{
	"CATALOG_PLAN_VIEW",
	"CATALOG_PLAN_MANAGE",
	"CATALOG_PLAN_PUBLISH",
	"COMMERCIAL_SUBSCRIPTION_VIEW",
	"COMMERCIAL_USAGE_VIEW",
	"COMMERCIAL_USAGE_RECONCILE",
	"PAYMENT_ORDER_VIEW",
	"PAYMENT_FULFILLMENT_VIEW",
	"PAYMENT_FULFILLMENT_REDRIVE",
	"USER_ADMIN_VIEW",
	"USER_ADMIN_MANAGE",
	"IAM_ROLE_VIEW",
	"IAM_ROLE_MANAGE",
	"ORGANIZATION_VIEW",
	"ORGANIZATION_MANAGE",
}
