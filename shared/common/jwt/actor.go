package _jwt

import _identitycontext "common/identitycontext"

/**
 * ActorFromPrincipal projects verified identity without entitlement state.
 * OrganizationID 0 is treated as absent (nil) so BindActor accepts personal tokens.
 */
func ActorFromPrincipal(principal *Principal) _identitycontext.Actor {
	if principal == nil {
		return _identitycontext.Actor{}
	}
	return _identitycontext.Actor{
		AuthID:         principal.AuthID,
		ProfileID:      principal.ProfileId,
		OriginID:       principal.OriginId,
		OrganizationID: nonZeroUint64Pointer(principal.OrganizationId),
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

/** nonZeroUint64Pointer clones value only when it is a positive ID. */
func nonZeroUint64Pointer(value *uint64) *uint64 {
	if value == nil || *value == 0 {
		return nil
	}
	clone := *value
	return &clone
}
