package _constants

type CtxKey string

const (
	ProfileIDKey      CtxKey = "profileId"
	OrganizationIDKey CtxKey = "organizationId"
	AuthIDKey         CtxKey = "authId"
	SessionKey        CtxKey = "session"
	RoleKey           CtxKey = "role"
	TypeKey           CtxKey = "type"
	PlanIDKey         CtxKey = "planId"
	PlanFromKey       CtxKey = "planFrom"
)
