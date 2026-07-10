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
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Order("event_date DESC NULLS LAST, created_at DESC").
		Find(&events).Error; err != nil {
		return nil, 0, err
	}

	if len(events) == 0 {
		return events, total, nil
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
		return events, total, nil
	}

	var docs []qh_domain.QHPlanningDocument
	if err := r.db.GetDB(ctx).
		Where("id IN ? AND deleted_at IS NULL", docIDs).
		Find(&docs).Error; err != nil {
		return nil, 0, err
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

	return events, total, nil
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

	if len(relations) == 0 {
		return []qh_domain.QHPlanningRelationItem{}, total, nil
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
		return nil, 0, err
	}

	projectMap := make(map[uint64]*qh_domain.QHPlanningProject, len(projects))
	for i := range projects {
		projectMap[projects[i].ID] = &projects[i]
	}

	items := make([]qh_domain.QHPlanningRelationItem, 0, len(relations))
	for _, rel := range relations {
		if rel.IsPeer {
			peerID := rel.ToID
			if req.PlanningProjectID != 0 && rel.ToID == req.PlanningProjectID {
				peerID = rel.FromID
			}
			items = append(items, qh_domain.QHPlanningRelationItem{
				Peer: projectMap[peerID],
			})
			continue
		}

		items = append(items, qh_domain.QHPlanningRelationItem{
			Parent: projectMap[rel.FromID],
			Child:  projectMap[rel.ToID],
		})
	}

	return items, total, nil
}
