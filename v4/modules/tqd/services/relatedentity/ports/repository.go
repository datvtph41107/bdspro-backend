package ports

import (
	"context"

	discoverydomain "mdn4/modules/tqd/entity/discovery"
	relateddomain "mdn4/modules/tqd/entity/relatedentity"
)

// EntityReader is intentionally narrow. Related Entity only needs canonical
// entity lookup; it must not depend on Discovery identify/search capabilities.
type EntityReader interface {
	GetEntity(ctx context.Context, ref discoverydomain.EntityRef) (*discoverydomain.EntityCandidate, error)
}

type Repository interface {
	EntityReader
	FindPlanningEntities(ctx context.Context, query relateddomain.Query) ([]relateddomain.Candidate, error)
	FindAdministrativeUnits(ctx context.Context, query relateddomain.Query) ([]relateddomain.Candidate, error)
	FindParcels(ctx context.Context, query relateddomain.Query) ([]relateddomain.Candidate, error)
	FindPOIs(ctx context.Context, query relateddomain.Query) ([]relateddomain.Candidate, error)
}
