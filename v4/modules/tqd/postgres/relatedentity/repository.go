package relatedentity

import (
	"context"
	"math"
	"strings"

	"userv4/modules/tqd/postgres/entitycandidate"
	discoverydomain "userv4/modules/tqd/entity/discovery"
	relateddomain "userv4/modules/tqd/entity/relatedentity"
	discoveryports "userv4/modules/tqd/services/discovery/ports"
	relatedports "userv4/modules/tqd/services/relatedentity/ports"

	"gorm.io/gorm"
)

type Repository struct {
	db           *gorm.DB
	entityReader relatedports.EntityReader
}

// NewEntityReader adapts the canonical Discovery read-model to the narrow port
// required by Related Entity. make wire discovers this constructor and keeps
// the module independent from Discovery search/identify operations.
func NewEntityReader(repository discoveryports.Repository) relatedports.EntityReader {
	return repository
}

func NewRepository(db *gorm.DB, entityReader relatedports.EntityReader) relatedports.Repository {
	return &Repository{
		db:           db,
		entityReader: entityReader,
	}
}

func (r *Repository) GetEntity(ctx context.Context, ref discoverydomain.EntityRef) (*discoverydomain.EntityCandidate, error) {
	return r.entityReader.GetEntity(ctx, ref)
}

type spatialRow = entitycandidate.Row

func rowsToCandidates(kind discoverydomain.EntityKind, rows []spatialRow) []relateddomain.Candidate {
	out := make([]relateddomain.Candidate, 0, len(rows))
	for _, row := range rows {
		relationship := relationshipKind(row.Relationship)
		out = append(out, relateddomain.Candidate{
			Entity: entitycandidate.Build(kind, row, entitycandidate.BuildOptions{
				MatchType: "related_entity",
			}),
			Relationship:   relationship,
			DistanceMeters: normalizedDistance(row.DistanceMeters),
			Score:          row.Score,
			ReasonCodes:    []string{"related_" + string(relationship)},
		})
	}
	return out
}

func relationshipKind(value string) relateddomain.RelationshipKind {
	normalized := relateddomain.RelationshipKind(strings.TrimSpace(value))
	if normalized.IsValid() {
		return normalized
	}
	return relateddomain.RelationshipUnknown
}

func normalizedDistance(value float64) float64 {
	if value < 0 || math.IsNaN(value) || math.IsInf(value, 0) {
		return math.MaxFloat64
	}
	return value
}
