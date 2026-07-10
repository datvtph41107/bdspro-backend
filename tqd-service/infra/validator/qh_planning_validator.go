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
