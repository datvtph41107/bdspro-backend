package _enum

type ContextKey string

const (
	ProfileIDKey      ContextKey = "profileId"
	OriginIDKey       ContextKey = "originId"
	OrganizationIDKey ContextKey = "organizationId"
	AuthIDKey         ContextKey = "authId"
	SessionKey        ContextKey = "session"
	RoleKey           ContextKey = "role"
	TypeKey           ContextKey = "type"
	PlanIDKey         ContextKey = "planId"
	PlanFromKey       ContextKey = "planFrom"
	ClientIPKey       ContextKey = "client-ip"
	UserAgentKey      ContextKey = "user-agent"
	APIKeyKey         ContextKey = "api-key"
	DeviceIDKey       ContextKey = "device-id"
	RoleIDsKey        ContextKey = "roleIds"
)
