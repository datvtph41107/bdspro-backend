package base_enum

type ContextKey string

const (
	ProfileIDKey      ContextKey = "profileId"
	OrganizationIDKey ContextKey = "organizationId"
	AuthIDKey         ContextKey = "authId"
	SessionKey        ContextKey = "session"
	RoleKey           ContextKey = "role"
	TypeKey           ContextKey = "type"
	PlanIDKey         ContextKey = "planId"
	PlanFromKey       ContextKey = "planFrom"
)
