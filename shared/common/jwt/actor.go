package _jwt

import _identitycontext "common/identitycontext"

/**
 * ActorFromPrincipal projects verified identity without entitlement state.
 */
func ActorFromPrincipal(principal *Principal) _identitycontext.Actor {
	if principal == nil {
		return _identitycontext.Actor{}
	}
	return _identitycontext.Actor{
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
