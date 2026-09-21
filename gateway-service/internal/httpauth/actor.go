package httpauth

import (
	"common/identity"
	"common/jwtverify"
)

func ActorFromPrincipal(p *jwtverify.Principal) identity.Actor {
	if p == nil {
		return identity.Actor{}
	}
	return identity.Actor{
		AuthID:         p.AuthID,
		ProfileID:      p.ProfileID,
		OriginID:       p.OriginID,
		OrganizationID: nonZeroUint64(p.OrganizationID),
		SessionID:      p.SessionID,
		Role:           p.Role,
		TokenType:      p.Type,
	}
}

func cloneUint64(v *uint64) *uint64 {
	if v == nil {
		return nil
	}
	x := *v
	return &x
}

// nonZeroUint64 clones only positive IDs. Organization 0 must become nil so
// personal tokens pass identity.BindActor (pointer-to-zero is invalid).
func nonZeroUint64(v *uint64) *uint64 {
	if v == nil || *v == 0 {
		return nil
	}
	x := *v
	return &x
}
