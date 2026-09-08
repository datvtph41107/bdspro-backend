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
		OrganizationID: cloneUint64(p.OrganizationID),
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
