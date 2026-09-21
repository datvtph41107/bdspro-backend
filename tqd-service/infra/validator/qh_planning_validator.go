package validator

import (
	_errors "common/errors"
	"strings"
	"tqd/internal"
	qh_domain "tqd/internal/domain/qh"
)

type QHPlanningValidator struct{}

func NewQHPlanningValidator() *QHPlanningValidator {
	return &QHPlanningValidator{}
}

func (v *QHPlanningValidator) ValidateProject(req *qh_domain.QHPlanningProject) error {
	if strings.TrimSpace(req.Name) == "" {
		return _errors.ReturnError(service.PlanningNameRequired)
	}
	if strings.TrimSpace(req.Code) == "" {
		return _errors.ReturnError(service.PlanningCodeRequired)
	}
	if req.PlanningType == 0 {
		return _errors.ReturnError(service.PlanningTypeRequired)
	}
	if req.PlanningLevel == 0 {
		return _errors.ReturnError(service.PlanningLevelRequired)
	}
	if strings.TrimSpace(req.ValidityStatus) == "" {
		return _errors.ReturnError(service.PlanningValidityStatusRequired)
	}
	return nil
}

func (v *QHPlanningValidator) ValidateDocument(req *qh_domain.QHPlanningDocument) error {
	if req.PlanningProjectID == 0 {
		return _errors.ReturnError(service.PlanningProjectIDRequired)
	}
	if strings.TrimSpace(req.Title) == "" {
		return _errors.ReturnError(service.PlanningTitleRequired)
	}
	if req.DocumentType == 0 {
		return _errors.ReturnError(service.PlanningDocumentTypeRequired)
	}
	return nil
}
