package application

import (
	"context"

	"tqd/internal/enums"
	planningdomain "tqd/internal/modules/planningclient/domain"
	publicdomain "tqd/internal/modules/publiccontent/domain"
)

// PlanningRepository contains persistence operations used by client-facing
// planning workflows. Keeping the interface in the application package makes
// transport handlers independent from PostgreSQL/GORM.
type PlanningRepository interface {
	ListProjects(context.Context, planningdomain.ProjectListFilter) (*planningdomain.Page[planningdomain.Project], error)
	GetProject(context.Context, uint64) (*planningdomain.Project, error)
	ListDocuments(context.Context, uint64, int, int, string, string, string) (*planningdomain.Page[planningdomain.Document], error)
	GetDocument(context.Context, uint64) (*planningdomain.Document, error)
	ListDocumentsByParcelRelation(context.Context, uint64, int, int, string) (*planningdomain.Page[planningdomain.Document], error)
	ListEvents(context.Context, uint64, int, int) (*planningdomain.Page[planningdomain.Event], error)
	Follow(context.Context, uint64, uint64, string) (*planningdomain.Follow, error)
	GetFollow(context.Context, uint64, uint64) (*planningdomain.Follow, error)
	Unfollow(context.Context, uint64, uint64) error
	ListFollowed(context.Context, uint64, int, int, string) (*planningdomain.Page[planningdomain.FollowedProject], error)
	ListProjectsByParcelRelation(context.Context, uint64, int, int) (*planningdomain.Page[planningdomain.Project], error)
	ListProjectsByLayerRelation(context.Context, []uint64, int, int) (*planningdomain.Page[planningdomain.Project], error)
	ListProjectLayers(context.Context, uint64) ([]planningdomain.ProjectLayer, error)
	ListProvinces(context.Context, planningdomain.ProvinceListFilter) ([]planningdomain.ProvinceItem, error)
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
func (s *Service) ListDocumentsByParcelRelation(ctx context.Context, parcelID uint64, page, size int, docType string) (*planningdomain.Page[planningdomain.Document], error) {
	return s.planning.ListDocumentsByParcelRelation(ctx, parcelID, page, size, docType)
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
func (s *Service) ListProjectsByParcelRelation(ctx context.Context, parcelID uint64, page, size int) (*planningdomain.Page[planningdomain.Project], error) {
	return s.planning.ListProjectsByParcelRelation(ctx, parcelID, page, size)
}
func (s *Service) ListProjectsByLayerRelation(ctx context.Context, layerIDs []uint64, page, size int) (*planningdomain.Page[planningdomain.Project], error) {
	return s.planning.ListProjectsByLayerRelation(ctx, layerIDs, page, size)
}
func (s *Service) ListProjectLayers(ctx context.Context, projectID uint64) ([]planningdomain.ProjectLayer, error) {
	return s.planning.ListProjectLayers(ctx, projectID)
}
func (s *Service) ListProvinces(ctx context.Context, filter planningdomain.ProvinceListFilter) ([]planningdomain.ProvinceItem, error) {
	items, err := s.planning.ListProvinces(ctx, filter)
	if err != nil {
		return nil, err
	}
	for i := range items {
		items[i].UnitTypeName = enums.LocationUnitType(items[i].UnitType).String()
		items[i].DataSourceName = enums.LocationDataSource(items[i].DataSource).String()
		if items[i].PlanDocID != nil && *items[i].PlanDocID == 0 {
			items[i].PlanDocID = nil
		}
	}
	return items, nil
}
func (s *Service) GetParcelQuickView(ctx context.Context, parcelID uint64) (*publicdomain.ParcelQuickView, error) {
	return s.projection.GetParcelQuickView(ctx, parcelID)
}
func (s *Service) GetPlanningProjectProjection(ctx context.Context, identity string) (*publicdomain.PlanningProjectProjection, error) {
	return s.projection.GetPlanningProjectProjection(ctx, identity)
}
