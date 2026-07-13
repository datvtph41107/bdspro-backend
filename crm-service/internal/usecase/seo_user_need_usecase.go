package usecase

import (
	_errors "common/errors"
	"context"
	"errors"
	"strings"
	"time"

	seo_strategy "crm/internal/domain/seo_strategy"
	"crm/internal/dto"
	"crm/internal/repo"
)

type SeoUserNeedUsecase struct {
	userNeedRepo  repo.SeoUserNeedRepo
	audienceRepo  repo.SeoAudienceRepo
	objectiveRepo repo.SeoObjectiveRepo
}

func NewSeoUserNeedUsecase(userNeedRepo repo.SeoUserNeedRepo, audienceRepo repo.SeoAudienceRepo, objectiveRepo repo.SeoObjectiveRepo) *SeoUserNeedUsecase {
	return &SeoUserNeedUsecase{userNeedRepo: userNeedRepo, audienceRepo: audienceRepo, objectiveRepo: objectiveRepo}
}

func (u *SeoUserNeedUsecase) Create(ctx context.Context, req *dto.CreateSeoUserNeedRequest) (*seo_strategy.SeoUserNeedAggregate, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không được để trống")
	}
	objective, err := u.objectiveRepo.GetObjectiveByID(ctx, req.ObjectiveID)
	if err != nil || objective == nil {
		return nil, _errors.ReturnError(404, "objective không tồn tại")
	}
	if objective.Status == seo_strategy.ObjectiveStatusArchived {
		return nil, _errors.ReturnError(409, "không thể tạo user need cho objective archived")
	}
	audience, err := u.audienceRepo.GetByID(ctx, req.AudienceSegmentID)
	if err != nil || audience == nil {
		return nil, _errors.ReturnError(404, "audience không tồn tại")
	}
	if audience.Status == seo_strategy.AudienceStatusArchived {
		return nil, _errors.ReturnError(409, "không thể tạo user need cho audience archived")
	}
	need, err := seo_strategy.CreateSeoUserNeedDraft(req.ObjectiveID, req.AudienceSegmentID, req.Code, req.Name, time.Now())
	if err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}
	existing, err := u.userNeedRepo.GetByCode(ctx, need.ObjectiveID, need.Code)
	if err == nil && existing != nil {
		return nil, _errors.ReturnError(409, "user need code đã tồn tại trong objective")
	}
	if err != nil && !errors.Is(err, repo.ErrSeoUserNeedNotFound) {
		return nil, err
	}
	if err := u.userNeedRepo.Create(ctx, need); err != nil {
		return nil, err
	}
	return u.Get(ctx, need.ID)
}

func (u *SeoUserNeedUsecase) Update(ctx context.Context, req *dto.UpdateSeoUserNeedRequest) (*seo_strategy.SeoUserNeedAggregate, error) {
	if req == nil || req.ID == 0 || req.ExpectedVersion == 0 {
		return nil, _errors.ReturnError(400, "user need update request không hợp lệ")
	}
	need, err := u.userNeedRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, mapUserNeedRepoError(err)
	}
	if need.Status == seo_strategy.UserNeedStatusArchived {
		return nil, _errors.ReturnError(409, "không thể sửa user need đã archived")
	}
	need.Name = strings.TrimSpace(req.Name)
	need.DecisionStage = seo_strategy.UserNeedDecisionStage(req.DecisionStage)
	need.Situation = strings.TrimSpace(req.Situation)
	need.TriggerStatement = strings.TrimSpace(req.TriggerStatement)
	need.JobToBeDone = strings.TrimSpace(req.JobToBeDone)
	need.BarrierStatement = strings.TrimSpace(req.BarrierStatement)
	need.CurrentAlternative = strings.TrimSpace(req.CurrentAlternative)
	need.DesiredOutcome = strings.TrimSpace(req.DesiredOutcome)
	need.FailureImpact = strings.TrimSpace(req.FailureImpact)
	need.RiskLevel = seo_strategy.UserNeedRiskLevel(req.RiskLevel)
	need.Priority = seo_strategy.UserNeedPriority(req.Priority)
	need.Confidence = seo_strategy.ObjectiveConfidence(req.Confidence)
	need.OwnerCode = strings.TrimSpace(req.OwnerCode)
	need.ReviewDueAt = req.ReviewDueAt
	if err := seo_strategy.ValidateUserNeedDefinition(need); err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}
	if err := u.userNeedRepo.UpdateDefinition(ctx, need, req.ExpectedVersion); err != nil {
		return nil, mapUserNeedRepoError(err)
	}
	return u.Get(ctx, need.ID)
}

func (u *SeoUserNeedUsecase) Get(ctx context.Context, id uint64) (*seo_strategy.SeoUserNeedAggregate, error) {
	if id == 0 {
		return nil, _errors.ReturnError(400, "user need id không hợp lệ")
	}
	need, err := u.userNeedRepo.GetByID(ctx, id)
	if err != nil {
		return nil, mapUserNeedRepoError(err)
	}
	objective, err := u.objectiveRepo.GetObjectiveByID(ctx, need.ObjectiveID)
	if err != nil || objective == nil {
		return nil, _errors.ReturnError(404, "objective của user need không tồn tại")
	}
	audience, err := u.audienceRepo.GetByID(ctx, need.AudienceSegmentID)
	if err != nil || audience == nil {
		return nil, _errors.ReturnError(404, "audience của user need không tồn tại")
	}
	evidence, err := u.userNeedRepo.ListLinkedEvidence(ctx, need.ID)
	if err != nil {
		return nil, err
	}
	return &seo_strategy.SeoUserNeedAggregate{
		Need: need, Objective: objective, Audience: audience, Evidence: evidence,
		Readiness: seo_strategy.EvaluateUserNeedReadiness(objective, audience, need, evidence, time.Now()),
	}, nil
}

func (u *SeoUserNeedUsecase) List(ctx context.Context, req *dto.SeoUserNeedListRequest) ([]seo_strategy.SeoUserNeed, int64, error) {
	if req == nil {
		req = &dto.SeoUserNeedListRequest{}
	}
	return u.userNeedRepo.List(ctx, repo.SeoUserNeedFilter{
		ObjectiveID: req.ObjectiveID, AudienceID: req.AudienceID,
		Keyword: req.Keyword, Status: req.Status, Confidence: req.Confidence, Priority: req.Priority,
		Page: req.Page, Size: req.Size,
	})
}

func (u *SeoUserNeedUsecase) LinkEvidence(
	ctx context.Context,
	req *dto.LinkSeoUserNeedEvidenceRequest,
) (*seo_strategy.SeoUserNeedAggregate, error) {
	if req == nil ||
		req.UserNeedID == 0 ||
		req.EvidenceID == 0 {
		return nil, _errors.ReturnError(
			400,
			"userNeedId và evidenceId không hợp lệ",
		)
	}

	need, err := u.userNeedRepo.GetByID(
		ctx,
		req.UserNeedID,
	)
	if err != nil {
		return nil, mapSeoUserNeedRepoError(err)
	}

	evidence, err := u.objectiveRepo.GetEvidenceByID(
		ctx,
		req.EvidenceID,
	)
	if err != nil {
		if errors.Is(
			err,
			repo.ErrSeoObjectiveEvidenceNotFound,
		) {
			return nil, _errors.ReturnError(
				404,
				"seo objective evidence không tồn tại",
			)
		}

		return nil, err
	}

	if need.ObjectiveID != evidence.ObjectiveID {
		return nil, _errors.ReturnError(
			409,
			"evidence và user need không thuộc cùng objective",
		)
	}

	relationType :=
		seo_strategy.UserNeedEvidenceRelation(
			req.RelationType,
		)

	if !seo_strategy.IsValidUserNeedEvidenceRelation(
		relationType,
	) {
		return nil, _errors.ReturnError(
			400,
			"relationType không hợp lệ",
		)
	}

	link := &seo_strategy.SeoUserNeedEvidenceLink{
		ObjectiveID:  need.ObjectiveID,
		UserNeedID:   need.ID,
		EvidenceID:   evidence.ID,
		RelationType: relationType,
		RelevanceNote: strings.TrimSpace(
			req.RelevanceNote,
		),
		CreatedAt: time.Now(),
	}

	if err := u.userNeedRepo.LinkEvidence(
		ctx,
		link,
	); err != nil {
		return nil, mapSeoUserNeedRepoError(err)
	}

	return u.Get(ctx, need.ID)
}

func (u *SeoUserNeedUsecase) UnlinkEvidence(ctx context.Context, req *dto.UnlinkSeoUserNeedEvidenceRequest) (*seo_strategy.SeoUserNeedAggregate, error) {
	if req == nil || req.UserNeedID == 0 || req.EvidenceID == 0 {
		return nil, _errors.ReturnError(400, "evidence unlink request không hợp lệ")
	}
	if err := u.userNeedRepo.UnlinkEvidence(ctx, req.UserNeedID, req.EvidenceID); err != nil {
		return nil, err
	}
	return u.Get(ctx, req.UserNeedID)
}

func (u *SeoUserNeedUsecase) Activate(ctx context.Context, req *dto.SeoUserNeedCommandRequest) (*seo_strategy.SeoUserNeedAggregate, error) {
	if req == nil || req.ID == 0 || req.ExpectedVersion == 0 {
		return nil, _errors.ReturnError(400, "user need command không hợp lệ")
	}
	aggregate, err := u.Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if aggregate.Need.Status == seo_strategy.UserNeedStatusArchived {
		return nil, _errors.ReturnError(409, "không thể activate user need đã archived")
	}
	if !aggregate.Readiness.Ready {
		codes := make([]string, 0, len(aggregate.Readiness.Blockers))
		for _, item := range aggregate.Readiness.Blockers {
			codes = append(codes, item.Code)
		}
		return nil, _errors.ReturnError(422, "user need chưa sẵn sàng: "+strings.Join(codes, ", "))
	}
	if err := u.userNeedRepo.UpdateStatus(ctx, req.ID, seo_strategy.UserNeedStatusActive, req.ExpectedVersion); err != nil {
		return nil, mapUserNeedRepoError(err)
	}
	return u.Get(ctx, req.ID)
}

func (u *SeoUserNeedUsecase) Pause(ctx context.Context, req *dto.SeoUserNeedCommandRequest) (*seo_strategy.SeoUserNeedAggregate, error) {
	return u.changeStatus(ctx, req, seo_strategy.UserNeedStatusPaused)
}

func (u *SeoUserNeedUsecase) Archive(ctx context.Context, req *dto.SeoUserNeedCommandRequest) (*seo_strategy.SeoUserNeedAggregate, error) {
	return u.changeStatus(ctx, req, seo_strategy.UserNeedStatusArchived)
}

func (u *SeoUserNeedUsecase) changeStatus(ctx context.Context, req *dto.SeoUserNeedCommandRequest, status seo_strategy.UserNeedStatus) (*seo_strategy.SeoUserNeedAggregate, error) {
	if req == nil || req.ID == 0 || req.ExpectedVersion == 0 {
		return nil, _errors.ReturnError(400, "user need command không hợp lệ")
	}
	if err := u.userNeedRepo.UpdateStatus(ctx, req.ID, status, req.ExpectedVersion); err != nil {
		return nil, mapUserNeedRepoError(err)
	}
	return u.Get(ctx, req.ID)
}

func mapUserNeedRepoError(err error) error {
	switch {
	case errors.Is(err, repo.ErrSeoUserNeedNotFound):
		return _errors.ReturnError(404, "seo user need không tồn tại")
	case errors.Is(err, repo.ErrSeoUserNeedConflict):
		return _errors.ReturnError(409, "user need đã được cập nhật; hãy tải lại dữ liệu")
	default:
		return err
	}
}

func mapSeoUserNeedRepoError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, repo.ErrSeoUserNeedNotFound):
		return _errors.ReturnError(
			404,
			"seo user need không tồn tại",
		)

	case errors.Is(err, repo.ErrSeoUserNeedConflict):
		return _errors.ReturnError(
			409,
			"user need đã được cập nhật bởi thao tác khác; hãy tải lại dữ liệu",
		)

	case errors.Is(
		err,
		repo.ErrSeoUserNeedEvidenceAlreadyLinked,
	):
		return _errors.ReturnError(
			409,
			"evidence đã được liên kết với user need",
		)

	case errors.Is(
		err,
		repo.ErrSeoUserNeedEvidenceNotLinked,
	):
		return _errors.ReturnError(
			404,
			"evidence chưa được liên kết với user need",
		)

	default:
		return err
	}
}
