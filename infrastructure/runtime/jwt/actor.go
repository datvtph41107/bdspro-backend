package _jwt

import "common/identity"

/**
 * ActorFromPrincipal projects verified identity without entitlement state.
 */
func ActorFromPrincipal(principal *Principal) identity.Actor {
	if principal == nil {
		return identity.Actor{}
	}
	return identity.Actor{
		AuthID:         principal.AuthID,
		ProfileID:      principal.ProfileId,
		OriginID:       principal.OriginId,
		OrganizationID: cloneUint64Pointer(principal.OrganizationId),
		SessionID:      principal.Session,
		Role:           principal.Role,
		TokenType:      principal.Type,
	}
}

func cloneUint64Pointer(value *uint64) *uint64 {
	if value == nil {
		return nil
	}
	clone := *value
	return &clone
}
