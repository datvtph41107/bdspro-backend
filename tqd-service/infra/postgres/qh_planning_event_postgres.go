package postgres

import (
	_db "common/db"
	"context"
	qh_domain "tqd/internal/domain/qh"
	qh_dto "tqd/internal/dto/qh"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
)

type QHPlanningEventRepo struct {
	db *_db.TransactionRepo
}

// NewQHPlanningEventRepo @bind: internal/interface/repo.IQHPlanningEventRepo
func NewQHPlanningEventRepo(db *_db.TransactionRepo) repo.IQHPlanningEventRepo {
	return &QHPlanningEventRepo{db: db}
}

func (r *QHPlanningEventRepo) Create(ctx context.Context, entity *qh_domain.QHPlanningEvent) error {
	normalizePlanningEvent(entity)
	return r.db.GetDB(ctx).Create(entity).Error
}

func (r *QHPlanningEventRepo) Update(ctx context.Context, id uint64, entity *qh_domain.QHPlanningEvent) error {
	normalizePlanningEvent(entity)
	updates := map[string]interface{}{
		"planning_project_id": entity.PlanningProjectID,
		"event_type":          entity.EventType,
		"event_name":          entity.EventName,
		"description":         entity.Description,
		"event_date":          entity.EventDate,
		"ver_no":              entity.VerNo,
		"legal_status":        entity.LegalStatus,
		"doc_ids":             entity.DocIDs,
		"metadata":            entity.Metadata,
	}
	return r.db.GetDB(ctx).Model(&qh_domain.QHPlanningEvent{}).Where("id = ?", id).Updates(updates).Error
}

func (r *QHPlanningEventRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.GetDB(ctx).Unscoped().Where("id = ?", id).Delete(&qh_domain.QHPlanningEvent{}).Error
}

func (r *QHPlanningEventRepo) GetByID(ctx context.Context, id uint64) (*qh_domain.QHPlanningEvent, error) {
	var entity qh_domain.QHPlanningEvent
	err := r.db.GetDB(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	events := []qh_domain.QHPlanningEvent{entity}
	if err := r.attachDocs(ctx, events); err != nil {
		return nil, err
	}
	return &events[0], nil
}

func normalizePlanningEvent(entity *qh_domain.QHPlanningEvent) {
	if entity == nil {
		return
	}
	if len(entity.Metadata) == 0 {
		entity.Metadata = []byte("{}")
	}
}

func (r *QHPlanningEventRepo) GetList(ctx context.Context, req *qh_dto.ListPlanningEventsRequest) ([]qh_domain.QHPlanningEvent, int64, error) {
	var events []qh_domain.QHPlanningEvent
	var total int64

	buildQuery := func() *gorm.DB {
		q := r.db.GetDB(ctx).Model(&qh_domain.QHPlanningEvent{}).Where("deleted_at IS NULL")
		if req.PlanningProjectID != 0 {
			q = q.Where("planning_project_id = ?", req.PlanningProjectID)
		}
		if req.Ver {
			q = q.Where("ver_no IS NOT NULL AND ver_no <> ''")
		}
		return q
	}

	if err := buildQuery().Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := buildQuery().
		Select("qh_planning_events.*").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Order("qh_planning_events.event_date DESC NULLS LAST, qh_planning_events.created_at DESC").
		Find(&events).Error; err != nil {
		return nil, 0, err
	}

	if err := r.attachDocs(ctx, events); err != nil {
		return nil, 0, err
	}
	return events, total, nil
}

func (r *QHPlanningEventRepo) CountByProject(ctx context.Context, projectID uint64) (int64, error) {
	var total int64
	q := r.db.GetDB(ctx).Model(&qh_domain.QHPlanningEvent{}).Where("deleted_at IS NULL")
	if projectID != 0 {
		q = q.Where("planning_project_id = ?", projectID)
	}
	if err := q.Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *QHPlanningEventRepo) attachDocs(ctx context.Context, events []qh_domain.QHPlanningEvent) error {
	if len(events) == 0 {
		return nil
	}

	docIDSet := make(map[uint64]struct{})
	for _, event := range events {
		for _, id := range event.DocIDs {
			docIDSet[uint64(id)] = struct{}{}
		}
	}

	docIDs := make([]uint64, 0, len(docIDSet))
	for id := range docIDSet {
		docIDs = append(docIDs, id)
	}
	if len(docIDs) == 0 {
		return nil
	}

	var docs []qh_domain.QHPlanningDocument
	if err := r.db.GetDB(ctx).
		Where("id IN ? AND deleted_at IS NULL", docIDs).
		Find(&docs).Error; err != nil {
		return err
	}

	docMap := make(map[uint64]qh_domain.QHPlanningDocument, len(docs))
	for _, doc := range docs {
		docMap[doc.ID] = doc
	}

	for i := range events {
		if len(events[i].DocIDs) == 0 {
			continue
		}
		events[i].Docs = make([]qh_domain.QHPlanningDocument, 0, len(events[i].DocIDs))
		for _, id := range events[i].DocIDs {
			if doc, ok := docMap[uint64(id)]; ok {
				events[i].Docs = append(events[i].Docs, doc)
			}
		}
	}
	return nil
}

func (r *QHPlanningEventRepo) CreateRelation(ctx context.Context, entity *qh_domain.QHPlanningRelation) error {
	normalizePlanningRelation(entity)
	return r.db.GetDB(ctx).Create(entity).Error
}

func (r *QHPlanningEventRepo) UpdateRelation(ctx context.Context, id uint64, entity *qh_domain.QHPlanningRelation) error {
	normalizePlanningRelation(entity)
	return r.db.GetDB(ctx).Model(&qh_domain.QHPlanningRelation{}).Where("id = ?", id).Updates(map[string]interface{}{
		"from_id":                     entity.FromID,
		"to_id":                       entity.ToID,
		"is_peer":                     entity.IsPeer,
		"relation_type":               entity.RelationType,
		"planning_project_id":         entity.PlanningProjectID,
		"related_planning_project_id": entity.RelatedPlanningProjectID,
		"metadata":                    entity.Metadata,
	}).Error
}

func (r *QHPlanningEventRepo) DeleteRelation(ctx context.Context, id uint64) error {
	return r.db.GetDB(ctx).Unscoped().Where("id = ?", id).Delete(&qh_domain.QHPlanningRelation{}).Error
}

func (r *QHPlanningEventRepo) GetRelationByID(ctx context.Context, id uint64) (*qh_domain.QHPlanningRelationItem, error) {
	var relation qh_domain.QHPlanningRelation
	err := r.db.GetDB(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&relation).Error
	if err != nil {
		return nil, err
	}
	items, err := r.toRelationItems(ctx, []qh_domain.QHPlanningRelation{relation}, 0)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &items[0], nil
}

func normalizePlanningRelation(entity *qh_domain.QHPlanningRelation) {
	if entity == nil {
		return
	}
	if entity.PlanningProjectID == 0 {
		entity.PlanningProjectID = entity.FromID
	}
	if entity.RelatedPlanningProjectID == 0 {
		entity.RelatedPlanningProjectID = entity.ToID
	}
	if entity.RelationType == "" {
		if entity.IsPeer {
			entity.RelationType = "peer"
		} else {
			entity.RelationType = "related"
		}
	}
	if len(entity.Metadata) == 0 {
		entity.Metadata = []byte("{}")
	}
}

func (r *QHPlanningEventRepo) GetRelations(ctx context.Context, req *qh_dto.ListPlanningEventsRequest) ([]qh_domain.QHPlanningRelationItem, int64, error) {
	var relations []qh_domain.QHPlanningRelation
	var total int64

	buildQuery := func() *gorm.DB {
		q := r.db.GetDB(ctx).Model(&qh_domain.QHPlanningRelation{}).Where("deleted_at IS NULL")
		if req.PlanningProjectID != 0 {
			q = q.Where("from_id = ? OR to_id = ?", req.PlanningProjectID, req.PlanningProjectID)
		}
		return q
	}

	if err := buildQuery().Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := buildQuery().
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Order("created_at DESC").
		Find(&relations).Error; err != nil {
		return nil, 0, err
	}

	items, err := r.toRelationItems(ctx, relations, req.PlanningProjectID)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *QHPlanningEventRepo) CountRelationsByProject(ctx context.Context, projectID uint64) (int64, error) {
	var total int64
	q := r.db.GetDB(ctx).Model(&qh_domain.QHPlanningRelation{}).Where("deleted_at IS NULL")
	if projectID != 0 {
		q = q.Where("from_id = ? OR to_id = ?", projectID, projectID)
	}
	if err := q.Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *QHPlanningEventRepo) toRelationItems(
	ctx context.Context,
	relations []qh_domain.QHPlanningRelation,
	planningProjectID uint64,
) ([]qh_domain.QHPlanningRelationItem, error) {
	if len(relations) == 0 {
		return []qh_domain.QHPlanningRelationItem{}, nil
	}

	projectIDs := make([]uint64, 0, len(relations)*2)
	for _, rel := range relations {
		projectIDs = append(projectIDs, rel.FromID, rel.ToID)
	}

	var projects []qh_domain.QHPlanningProject
	if err := r.db.GetDB(ctx).
		Preload("Jurisdiction").
		Preload("Layers").
		Where("id IN ? AND deleted_at IS NULL", projectIDs).
		Find(&projects).Error; err != nil {
		return nil, err
	}

	projectMap := make(map[uint64]*qh_domain.QHPlanningProject, len(projects))
	for i := range projects {
		projectMap[projects[i].ID] = &projects[i]
	}

	items := make([]qh_domain.QHPlanningRelationItem, 0, len(relations))
	for i := range relations {
		rel := relations[i]
		item := qh_domain.QHPlanningRelationItem{Relation: &rel}
		if rel.IsPeer {
			peerID := rel.ToID
			if planningProjectID != 0 && rel.ToID == planningProjectID {
				peerID = rel.FromID
			}
			item.Peer = projectMap[peerID]
		} else {
			item.Parent = projectMap[rel.FromID]
			item.Child = projectMap[rel.ToID]
		}
		items = append(items, item)
	}
	return items, nil
}
