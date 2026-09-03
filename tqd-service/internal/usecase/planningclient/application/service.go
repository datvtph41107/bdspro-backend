package application

import (
	"context"

	planningdomain "tqd/internal/domain/planningclient/model"
	publicdomain "tqd/internal/domain/publiccontent/model"
)

// PlanningRepository contains persistence operations used by client-facing
// planning workflows. Keeping the interface in the application package makes
// transport handlers independent from PostgreSQL/GORM.
type PlanningRepository interface {
	ListProjects(context.Context, planningdomain.ProjectListFilter) (*planningdomain.Page[planningdomain.Project], error)
	GetProject(context.Context, uint64) (*planningdomain.Project, error)
	ListDocuments(context.Context, uint64, int, int, string, string, string) (*planningdomain.Page[planningdomain.Document], error)
	GetDocument(context.Context, uint64) (*planningdomain.Document, error)
	ListEvents(context.Context, uint64, int, int) (*planningdomain.Page[planningdomain.Event], error)
	Follow(context.Context, uint64, uint64, string) (*planningdomain.Follow, error)
	GetFollow(context.Context, uint64, uint64) (*planningdomain.Follow, error)
	Unfollow(context.Context, uint64, uint64) error
	ListFollowed(context.Context, uint64, int, int, string) (*planningdomain.Page[planningdomain.FollowedProject], error)
	ListProjectLayers(context.Context, uint64) ([]planningdomain.ProjectLayer, error)
}

// ProjectionRepository owns source projections that CRM consumes through
// TQD gRPC. It is intentionally separate from the planning CRUD repository.
type ProjectionRepository interface {
	GetParcelQuickView(context.Context, uint64) (*publicdomain.ParcelQuickView, error)
	GetPlanningProjectProjection(context.Context, string) (*publicdomain.PlanningProjectProjection, error)
}

type Service struct {
	planning   PlanningRepository
	projection ProjectionRepository
}

func NewService(planning PlanningRepository, projection ProjectionRepository) *Service {
	return &Service{planning: planning, projection: projection}
}

func (s *Service) ListProjects(ctx context.Context, filter planningdomain.ProjectListFilter) (*planningdomain.Page[planningdomain.Project], error) {
	return s.planning.ListProjects(ctx, filter)
}
func (s *Service) GetProject(ctx context.Context, id uint64) (*planningdomain.Project, error) {
	return s.planning.GetProject(ctx, id)
}
func (s *Service) ListDocuments(ctx context.Context, projectID uint64, page, size int, keyword, documentType, status string) (*planningdomain.Page[planningdomain.Document], error) {
	return s.planning.ListDocuments(ctx, projectID, page, size, keyword, documentType, status)
}
func (s *Service) GetDocument(ctx context.Context, id uint64) (*planningdomain.Document, error) {
	return s.planning.GetDocument(ctx, id)
}
func (s *Service) ListEvents(ctx context.Context, projectID uint64, page, size int) (*planningdomain.Page[planningdomain.Event], error) {
	return s.planning.ListEvents(ctx, projectID, page, size)
}
func (s *Service) Follow(ctx context.Context, userID, projectID uint64, note string) (*planningdomain.Follow, error) {
	return s.planning.Follow(ctx, userID, projectID, note)
}
func (s *Service) GetFollow(ctx context.Context, userID, projectID uint64) (*planningdomain.Follow, error) {
	return s.planning.GetFollow(ctx, userID, projectID)
}
func (s *Service) Unfollow(ctx context.Context, userID, projectID uint64) error {
	return s.planning.Unfollow(ctx, userID, projectID)
}
func (s *Service) ListFollowed(ctx context.Context, userID uint64, page, size int, search string) (*planningdomain.Page[planningdomain.FollowedProject], error) {
	return s.planning.ListFollowed(ctx, userID, page, size, search)
}
func (s *Service) ListProjectLayers(ctx context.Context, projectID uint64) ([]planningdomain.ProjectLayer, error) {
	return s.planning.ListProjectLayers(ctx, projectID)
}
func (s *Service) GetParcelQuickView(ctx context.Context, parcelID uint64) (*publicdomain.ParcelQuickView, error) {
	return s.projection.GetParcelQuickView(ctx, parcelID)
}
func (s *Service) GetPlanningProjectProjection(ctx context.Context, identity string) (*publicdomain.PlanningProjectProjection, error) {
	return s.projection.GetPlanningProjectProjection(ctx, identity)
}
