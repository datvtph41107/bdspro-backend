// internal/domain/enums/actor_role.go
package enums

type ActorRole int32

const (
	ActorRoleOwner        ActorRole = 10
	ActorRoleAdmin        ActorRole = 20
	ActorRoleSystem       ActorRole = 30
	ActorRoleCollaborator ActorRole = 40
)

func (e ActorRole) String() string {
	switch e {
	case ActorRoleOwner:
		return "OWNER"
	case ActorRoleAdmin:
		return "ADMIN"
	case ActorRoleSystem:
		return "SYSTEM"
	case ActorRoleCollaborator:
		return "COLLABORATOR"
	default:
		return "UNKNOWN"
	}
}

func (e ActorRole) HasFullAccess() bool {
	return e == ActorRoleAdmin || e == ActorRoleSystem
}

func (e ActorRole) CanBypassRestrictions() bool {
	return e == ActorRoleAdmin
}
