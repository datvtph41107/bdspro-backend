package usecase

import (
	_errors "common/errors"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	seo_strategy "crm/internal/domain/seo_strategy"
	"crm/internal/dto"
	"crm/internal/repo"
)

type SeoObjectiveUsecase struct {
	repo repo.SeoObjectiveRepo
}

func NewSeoObjectiveUsecase(
	repository repo.SeoObjectiveRepo,
) *SeoObjectiveUsecase {
	return &SeoObjectiveUsecase{repo: repository}
}

func (u *SeoObjectiveUsecase) Create(
	ctx context.Context,
	req *dto.CreateSeoObjectiveRequest,
) (*seo_strategy.SeoObjectiveAggregate, error) {
	if req == nil {
		return nil, _errors.ReturnError(
			400,
			"request không được để trống",
		)
	}

	objective, err := seo_strategy.CreateSeoObjectiveDraft(
		req.Code,
		req.Name,
		time.Now(),
	)
	if err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}

	existing, err := u.repo.GetObjectiveByCode(ctx, objective.Code)
	if err == nil && existing != nil {
		return nil, _errors.ReturnError(
			409,
			"objective code đã tồn tại",
		)
	}
	if err != nil &&
		!errors.Is(err, repo.ErrSeoObjectiveNotFound) {
		return nil, err
	}

	if err := u.repo.CreateObjective(ctx, objective); err != nil {
		return nil, err
	}

	return u.Get(ctx, objective.ID)
}

func (u *SeoObjectiveUsecase) UpdateDefinition(
	ctx context.Context,
	req *dto.UpdateSeoObjectiveDefinitionRequest,
) (*seo_strategy.SeoObjectiveAggregate, error) {
	if req == nil || req.ID == 0 {
		return nil, _errors.ReturnError(
			400,
			"objective id không hợp lệ",
		)
	}

	if req.ExpectedVersion == 0 {
		return nil, _errors.ReturnError(
			400,
			"expectedVersion không hợp lệ",
		)
	}

	objective, err := u.repo.GetObjectiveByID(ctx, req.ID)
	if err != nil {
		return nil, mapObjectiveRepoError(err)
	}

	if objective.Status == seo_strategy.ObjectiveStatusArchived {
		return nil, _errors.ReturnError(
			409,
			"không thể sửa objective đã archived",
		)
	}

	objective.Name = strings.TrimSpace(req.Name)
	objective.AudienceSummary = strings.TrimSpace(req.AudienceSummary)
	objective.Situation = strings.TrimSpace(req.Situation)
	objective.ProblemStatement = strings.TrimSpace(req.ProblemStatement)
	objective.UserOutcome = strings.TrimSpace(req.UserOutcome)
	objective.ProductOutcome = strings.TrimSpace(req.ProductOutcome)
	objective.BusinessOutcome = strings.TrimSpace(req.BusinessOutcome)
	objective.NonGoals = compactObjectiveStrings(req.NonGoals)
	objective.Constraints = compactObjectiveStrings(req.Constraints)
	objective.Confidence = seo_strategy.ObjectiveConfidence(
		req.Confidence,
	)
	objective.OwnerCode = strings.TrimSpace(req.OwnerCode)
	objective.ReviewDueAt = req.ReviewDueAt

	if err := seo_strategy.ValidateObjectiveDefinition(objective); err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}

	if err := u.repo.UpdateObjectiveDefinition(
		ctx,
		objective,
		req.ExpectedVersion,
	); err != nil {
		return nil, mapObjectiveRepoError(err)
	}

	return u.Get(ctx, objective.ID)
}

func (u *SeoObjectiveUsecase) AddEvidence(
	ctx context.Context,
	req *dto.AddSeoObjectiveEvidenceRequest,
) (*seo_strategy.SeoObjectiveAggregate, error) {
	if req == nil || req.ObjectiveID == 0 {
		return nil, _errors.ReturnError(
			400,
			"objectiveId không hợp lệ",
		)
	}

	if _, err := u.repo.GetObjectiveByID(
		ctx,
		req.ObjectiveID,
	); err != nil {
		return nil, mapObjectiveRepoError(err)
	}

	metadata := req.Metadata
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}

	evidence := &seo_strategy.SeoObjectiveEvidence{
		ObjectiveID: req.ObjectiveID,
		Type: seo_strategy.EvidenceType(
			req.EvidenceType,
		),
		Stance: seo_strategy.EvidenceStance(
			req.Stance,
		),
		SourceSystem: strings.TrimSpace(
			req.SourceSystem,
		),
		SourceReference: strings.TrimSpace(
			req.SourceReference,
		),
		Finding: strings.TrimSpace(req.Finding),
		Confidence: seo_strategy.ObjectiveConfidence(
			req.Confidence,
		),
		ObservedAt: req.ObservedAt,
		ValidUntil: req.ValidUntil,
		Status:     seo_strategy.EvidenceStatusActive,
		Metadata:   metadata,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := seo_strategy.ValidateEvidence(evidence); err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}

	if err := u.repo.CreateEvidence(ctx, evidence); err != nil {
		return nil, err
	}

	return u.Get(ctx, req.ObjectiveID)
}

func (u *SeoObjectiveUsecase) UpsertMetric(
	ctx context.Context,
	req *dto.UpsertSeoObjectiveMetricRequest,
) (*seo_strategy.SeoObjectiveAggregate, error) {
	if req == nil || req.ObjectiveID == 0 {
		return nil, _errors.ReturnError(
			400,
			"objectiveId không hợp lệ",
		)
	}

	if _, err := u.repo.GetObjectiveByID(
		ctx,
		req.ObjectiveID,
	); err != nil {
		return nil, mapObjectiveRepoError(err)
	}

	status := seo_strategy.MetricStatus(req.Status)
	if status == "" {
		status = seo_strategy.MetricStatusActive
	}

	unit := strings.TrimSpace(req.Unit)
	if unit == "" {
		unit = "count"
	}

	metric := &seo_strategy.SeoObjectiveMetric{
		ObjectiveID: req.ObjectiveID,
		Code:        req.MetricCode,
		Name:        req.MetricName,
		Layer: seo_strategy.MetricLayer(
			req.MetricLayer,
		),
		Description: strings.TrimSpace(req.Description),
		SourceSystem: strings.TrimSpace(
			req.SourceSystem,
		),
		CalculationMethod: strings.TrimSpace(
			req.CalculationMethod,
		),
		Unit: unit,
		TargetDirection: seo_strategy.TargetDirection(
			req.TargetDirection,
		),
		BaselineValue: req.BaselineValue,
		TargetValue:   req.TargetValue,
		IsPrimary:     req.IsPrimary,
		InstrumentationStatus: seo_strategy.InstrumentationStatus(
			req.InstrumentationStatus,
		),
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := seo_strategy.ValidateMetric(metric); err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}

	if err := u.repo.UpsertMetric(ctx, metric); err != nil {
		return nil, err
	}

	return u.Get(ctx, req.ObjectiveID)
}

func (u *SeoObjectiveUsecase) Get(
	ctx context.Context,
	id uint64,
) (*seo_strategy.SeoObjectiveAggregate, error) {
	if id == 0 {
		return nil, _errors.ReturnError(
			400,
			"objective id không hợp lệ",
		)
	}

	objective, err := u.repo.GetObjectiveByID(ctx, id)
	if err != nil {
		return nil, mapObjectiveRepoError(err)
	}

	evidence, err := u.repo.ListEvidence(ctx, id)
	if err != nil {
		return nil, err
	}

	metrics, err := u.repo.ListMetrics(ctx, id)
	if err != nil {
		return nil, err
	}

	readiness := seo_strategy.EvaluateObjectiveReadiness(
		objective,
		evidence,
		metrics,
		time.Now(),
	)

	return &seo_strategy.SeoObjectiveAggregate{
		Objective: objective,
		Evidence:  evidence,
		Metrics:   metrics,
		Readiness: readiness,
	}, nil
}

func (u *SeoObjectiveUsecase) List(
	ctx context.Context,
	req *dto.SeoObjectiveListRequest,
) ([]seo_strategy.SeoObjective, int64, error) {
	if req == nil {
		req = &dto.SeoObjectiveListRequest{}
	}

	return u.repo.ListObjectives(
		ctx,
		repo.SeoObjectiveFilter{
			Keyword:    req.Keyword,
			Status:     req.Status,
			Confidence: req.Confidence,
			Page:       req.Page,
			Size:       req.Size,
		},
	)
}

func (u *SeoObjectiveUsecase) Activate(
	ctx context.Context,
	req *dto.SeoObjectiveCommandRequest,
) (*seo_strategy.SeoObjectiveAggregate, error) {
	aggregate, err := u.Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if aggregate.Objective.Status ==
		seo_strategy.ObjectiveStatusArchived {
		return nil, _errors.ReturnError(
			409,
			"không thể activate objective đã archived",
		)
	}

	if !aggregate.Readiness.Ready {
		codes := make([]string, 0, len(aggregate.Readiness.Blockers))
		for _, item := range aggregate.Readiness.Blockers {
			codes = append(codes, item.Code)
		}

		return nil, _errors.ReturnError(
			422,
			"objective chưa sẵn sàng: "+strings.Join(codes, ", "),
		)
	}

	if err := u.repo.UpdateObjectiveStatus(
		ctx,
		req.ID,
		seo_strategy.ObjectiveStatusActive,
		req.ExpectedVersion,
	); err != nil {
		return nil, mapObjectiveRepoError(err)
	}

	return u.Get(ctx, req.ID)
}

func (u *SeoObjectiveUsecase) Pause(
	ctx context.Context,
	req *dto.SeoObjectiveCommandRequest,
) (*seo_strategy.SeoObjectiveAggregate, error) {
	if err := u.repo.UpdateObjectiveStatus(
		ctx,
		req.ID,
		seo_strategy.ObjectiveStatusPaused,
		req.ExpectedVersion,
	); err != nil {
		return nil, mapObjectiveRepoError(err)
	}

	return u.Get(ctx, req.ID)
}

func (u *SeoObjectiveUsecase) Archive(
	ctx context.Context,
	req *dto.SeoObjectiveCommandRequest,
) (*seo_strategy.SeoObjectiveAggregate, error) {
	if err := u.repo.UpdateObjectiveStatus(
		ctx,
		req.ID,
		seo_strategy.ObjectiveStatusArchived,
		req.ExpectedVersion,
	); err != nil {
		return nil, mapObjectiveRepoError(err)
	}

	return u.Get(ctx, req.ID)
}

func mapObjectiveRepoError(err error) error {
	switch {
	case errors.Is(err, repo.ErrSeoObjectiveNotFound):
		return _errors.ReturnError(
			404,
			"seo objective không tồn tại",
		)

	case errors.Is(err, repo.ErrSeoObjectiveConflict):
		return _errors.ReturnError(
			409,
			"objective đã được cập nhật; hãy tải lại dữ liệu",
		)

	default:
		return err
	}
}

func compactObjectiveStrings(values []string) []string {
	result := make([]string, 0, len(values))

	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			result = append(result, value)
		}
	}

	return result
}
