package validator

import (
	_errors "common/errors"
	"strings"
	qh_domain "tqd/internal/domain/qh"
)

type QHPlanningValidator struct{}

func NewQHPlanningValidator() *QHPlanningValidator {
	return &QHPlanningValidator{}
}

func (v *QHPlanningValidator) ValidateProject(req *qh_domain.QHPlanningProject) error {
	if strings.TrimSpace(req.Name) == "" {
		return _errors.ReturnError(400, "name is required")
	}
	if strings.TrimSpace(req.Code) == "" {
		return _errors.ReturnError(400, "code is required")
	}
	if req.PlanningType == 0 {
		return _errors.ReturnError(400, "planningType is required")
	}
	if req.PlanningLevel == 0 {
		return _errors.ReturnError(400, "planningLevel is required")
	}
	if strings.TrimSpace(req.ValidityStatus) == "" {
		return _errors.ReturnError(400, "validityStatus is required")
	}
	return nil
}

func (v *QHPlanningValidator) ValidateEvent(req *qh_domain.QHPlanningEvent) error {
	if req == nil {
		return _errors.ReturnError(400, "event is required")
	}
	if req.PlanningProjectID == 0 {
		return _errors.ReturnError(400, "planningProjectId is required")
	}
	if strings.TrimSpace(req.EventName) == "" {
		return _errors.ReturnError(400, "eventName is required")
	}
	return nil
}

func (v *QHPlanningValidator) ValidateRelation(req *qh_domain.QHPlanningRelation) error {
	if req == nil {
		return _errors.ReturnError(400, "relation is required")
	}
	if req.FromID == 0 {
		return _errors.ReturnError(400, "fromId is required")
	}
	if req.ToID == 0 {
		return _errors.ReturnError(400, "toId is required")
	}
	if req.FromID == req.ToID {
		return _errors.ReturnError(400, "fromId và toId không được trùng nhau")
	}
	return nil
}

func (v *QHPlanningValidator) ValidateDocument(req *qh_domain.QHPlanningDocument) error {
	if req.PlanningProjectID == 0 {
		return _errors.ReturnError(400, "planningProjectId is required")
	}
	if strings.TrimSpace(req.Title) == "" {
		return _errors.ReturnError(400, "title is required")
	}
	if req.DocumentType == 0 {
		return _errors.ReturnError(400, "documentType is required")
	}
	return nil
}

func (v *QHPlanningValidator) ValidateAssignLayer(layerID uint64, docIDs []uint64) error {
	if layerID == 0 {
		return _errors.ReturnError(400, "layerId is required")
	}
	for _, docID := range docIDs {
		if docID > 0 {
			return nil
		}
	}
	return _errors.ReturnError(400, "docIds is required")
}

func (v *QHPlanningValidator) ValidateUnassignLayer(planningID, layerID uint64) error {
	if planningID == 0 {
		return _errors.ReturnError(400, "planningId is required")
	}
	if layerID == 0 {
		return _errors.ReturnError(400, "layerId is required")
	}
	return nil
}
