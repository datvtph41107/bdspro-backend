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

type seoIntentUserNeedReader interface {
	GetByID(ctx context.Context, id uint64) (*seo_strategy.SeoUserNeed, error)
}

type seoIntentQueryEvidenceReader interface {
	GetByID(ctx context.Context, id uint64) (*seo_strategy.SeoQueryEvidence, error)
}

type SeoIntentHypothesisUsecase struct {
	repository          repo.SeoIntentHypothesisRepo
	userNeedReader      seoIntentUserNeedReader
	queryEvidenceReader seoIntentQueryEvidenceReader
	now                 func() time.Time
}

func NewSeoIntentHypothesisUsecase(
	repository repo.SeoIntentHypothesisRepo,
	userNeedRepository repo.SeoUserNeedRepo,
	queryEvidenceRepository repo.SeoQueryEvidenceRepo,
) *SeoIntentHypothesisUsecase {
	return &SeoIntentHypothesisUsecase{
		repository:          repository,
		userNeedReader:      userNeedRepository,
		queryEvidenceReader: queryEvidenceRepository,
		now:                 time.Now,
	}
}

func (u *SeoIntentHypothesisUsecase) currentTime() time.Time {
	if u.now == nil {
		return time.Now().UTC()
	}
	return u.now().UTC()
}

func (u *SeoIntentHypothesisUsecase) Create(
	ctx context.Context,
	req *dto.CreateSeoIntentHypothesisRequest,
) (*seo_strategy.SeoIntentHypothesisAggregate, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không được để trống")
	}
	if req.UserNeedID == 0 {
		return nil, _errors.ReturnError(400, "userNeedId không hợp lệ")
	}

	need, err := u.userNeedReader.GetByID(ctx, req.UserNeedID)
	if err != nil {
		return nil, mapIntentUserNeedError(err)
	}
	if need == nil {
		return nil, _errors.ReturnError(404, "seo user need không tồn tại")
	}
	if need.Status != seo_strategy.UserNeedStatusActive {
		return nil, _errors.ReturnError(
			409,
			"user need phải active trước khi tạo intent hypothesis",
		)
	}

	intent, err := seo_strategy.CreateSeoIntentHypothesisDraft(
		need.ObjectiveID,
		need.ID,
		req.Code,
		req.Name,
		u.currentTime(),
	)
	if err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}

	existing, err := u.repository.GetByCode(ctx, need.ID, intent.Code)
	if err == nil && existing != nil {
		return nil, _errors.ReturnError(
			409,
			"intent hypothesis code đã tồn tại trong user need",
		)
	}
	if err != nil && !errors.Is(err, repo.ErrSeoIntentHypothesisNotFound) {
		return nil, mapIntentRepoError(err)
	}

	if err := u.repository.Create(ctx, intent); err != nil {
		return nil, mapIntentRepoError(err)
	}
	return u.Get(ctx, intent.ID)
}

func (u *SeoIntentHypothesisUsecase) Update(
	ctx context.Context,
	req *dto.UpdateSeoIntentHypothesisRequest,
) (*seo_strategy.SeoIntentHypothesisAggregate, error) {
	if req == nil || req.ID == 0 || req.ExpectedVersion == 0 {
		return nil, _errors.ReturnError(400, "intent hypothesis update request không hợp lệ")
	}

	intent, err := u.repository.GetByID(ctx, req.ID)
	if err != nil {
		return nil, mapIntentRepoError(err)
	}
	if intent.Status == seo_strategy.IntentHypothesisStatusArchived {
		return nil, _errors.ReturnError(409, "không thể sửa intent hypothesis đã archived")
	}

	searchAction, err := seo_strategy.ParseIntentSearchAction(req.SearchAction)
	if err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}
	subjectType, err := seo_strategy.ParseIntentSubjectType(req.SubjectType)
	if err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}
	scopeType, err := seo_strategy.ParseIntentScopeType(req.ScopeType)
	if err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}
	confidence := seo_strategy.ObjectiveConfidence(strings.TrimSpace(req.Confidence))
	if confidence == "" {
		confidence = seo_strategy.ObjectiveConfidenceExploratory
	}
	if !seo_strategy.IsValidObjectiveConfidence(confidence) {
		return nil, _errors.ReturnError(400, "intent hypothesis confidence không hợp lệ")
	}

	intent.Name = strings.TrimSpace(req.Name)
	intent.GoalStatement = strings.TrimSpace(req.GoalStatement)
	intent.PrimaryQuestion = strings.TrimSpace(req.PrimaryQuestion)
	intent.ExpectedResult = strings.TrimSpace(req.ExpectedResult)
	intent.SearchAction = searchAction
	intent.SubjectType = subjectType
	intent.ScopeType = scopeType
	intent.Confidence = confidence
	intent.OwnerCode = strings.TrimSpace(req.OwnerCode)
	intent.ReviewDueAt = req.ReviewDueAt
	intent.UpdatedAt = u.currentTime()

	if err := seo_strategy.ValidateIntentHypothesisDefinition(intent); err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}
	if err := u.repository.UpdateDefinition(ctx, intent, req.ExpectedVersion); err != nil {
		return nil, mapIntentRepoError(err)
	}
	return u.Get(ctx, intent.ID)
}

func (u *SeoIntentHypothesisUsecase) Get(
	ctx context.Context,
	id uint64,
) (*seo_strategy.SeoIntentHypothesisAggregate, error) {
	if id == 0 {
		return nil, _errors.ReturnError(400, "intent hypothesis id không hợp lệ")
	}

	intent, err := u.repository.GetByID(ctx, id)
	if err != nil {
		return nil, mapIntentRepoError(err)
	}
	need, err := u.userNeedReader.GetByID(ctx, intent.UserNeedID)
	if err != nil {
		return nil, mapIntentUserNeedError(err)
	}
	if need == nil {
		return nil, _errors.ReturnError(404, "user need của intent hypothesis không tồn tại")
	}
	linkedEvidence, err := u.repository.ListLinkedEvidence(ctx, intent.ID)
	if err != nil {
		return nil, mapIntentRepoError(err)
	}

	return &seo_strategy.SeoIntentHypothesisAggregate{
		Intent:   intent,
		UserNeed: need,
		Evidence: linkedEvidence,
		Readiness: seo_strategy.EvaluateIntentHypothesisReadiness(
			need,
			intent,
			linkedEvidence,
			u.currentTime(),
		),
	}, nil
}

func (u *SeoIntentHypothesisUsecase) List(
	ctx context.Context,
	req *dto.SeoIntentHypothesisListRequest,
) ([]seo_strategy.SeoIntentHypothesis, int64, error) {
	if req == nil {
		req = &dto.SeoIntentHypothesisListRequest{}
	}

	filter := repo.SeoIntentHypothesisFilter{
		ObjectiveID: req.ObjectiveID,
		UserNeedID:  req.UserNeedID,
		Keyword:     strings.TrimSpace(req.Keyword),
		Page:        req.Page,
		Size:        req.Size,
	}

	if value := strings.TrimSpace(req.SearchAction); value != "" {
		parsed, err := seo_strategy.ParseIntentSearchAction(value)
		if err != nil {
			return nil, 0, _errors.ReturnError(400, err.Error())
		}
		filter.SearchAction = parsed
	}
	if value := strings.TrimSpace(req.SubjectType); value != "" {
		parsed, err := seo_strategy.ParseIntentSubjectType(value)
		if err != nil {
			return nil, 0, _errors.ReturnError(400, err.Error())
		}
		filter.SubjectType = parsed
	}
	if value := strings.TrimSpace(req.ScopeType); value != "" {
		parsed, err := seo_strategy.ParseIntentScopeType(value)
		if err != nil {
			return nil, 0, _errors.ReturnError(400, err.Error())
		}
		filter.ScopeType = parsed
	}
	if value := strings.TrimSpace(req.Status); value != "" {
		parsed, err := seo_strategy.ParseIntentHypothesisStatus(value)
		if err != nil {
			return nil, 0, _errors.ReturnError(400, err.Error())
		}
		filter.Status = parsed
	}
	if value := strings.TrimSpace(req.Confidence); value != "" {
		confidence := seo_strategy.ObjectiveConfidence(value)
		if !seo_strategy.IsValidObjectiveConfidence(confidence) {
			return nil, 0, _errors.ReturnError(400, "confidence không hợp lệ")
		}
		filter.Confidence = confidence
	}

	items, total, err := u.repository.List(ctx, filter)
	if err != nil {
		return nil, 0, mapIntentRepoError(err)
	}
	return items, total, nil
}

func (u *SeoIntentHypothesisUsecase) LinkEvidence(
	ctx context.Context,
	req *dto.LinkSeoIntentEvidenceRequest,
) (*seo_strategy.SeoIntentHypothesisAggregate, error) {
	if req == nil || req.IntentHypothesisID == 0 || req.QueryEvidenceID == 0 {
		return nil, _errors.ReturnError(400, "intent evidence link request không hợp lệ")
	}

	intent, err := u.repository.GetByID(ctx, req.IntentHypothesisID)
	if err != nil {
		return nil, mapIntentRepoError(err)
	}
	if intent.Status == seo_strategy.IntentHypothesisStatusArchived {
		return nil, _errors.ReturnError(409, "không thể link evidence vào intent hypothesis đã archived")
	}

	evidence, err := u.queryEvidenceReader.GetByID(ctx, req.QueryEvidenceID)
	if err != nil {
		if errors.Is(err, repo.ErrSeoQueryEvidenceNotFound) {
			return nil, _errors.ReturnError(404, "seo query evidence không tồn tại")
		}
		return nil, err
	}
	if evidence == nil {
		return nil, _errors.ReturnError(404, "seo query evidence không tồn tại")
	}
	if evidence.ObjectiveID != intent.ObjectiveID || evidence.UserNeedID != intent.UserNeedID {
		return nil, _errors.ReturnError(
			409,
			"intent hypothesis và query evidence không cùng user need/objective",
		)
	}
	if !seo_strategy.IsQueryEvidenceCurrent(*evidence, u.currentTime()) {
		return nil, _errors.ReturnError(
			409,
			"chỉ Query Evidence active và còn hiệu lực mới được link vào intent hypothesis",
		)
	}

	relationType, err := seo_strategy.ParseIntentEvidenceRelation(req.RelationType)
	if err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}
	now := u.currentTime()
	link := &seo_strategy.SeoIntentEvidenceLink{
		ObjectiveID:        intent.ObjectiveID,
		UserNeedID:         intent.UserNeedID,
		IntentHypothesisID: intent.ID,
		QueryEvidenceID:    evidence.ID,
		RelationType:       relationType,
		RelevanceNote:      strings.TrimSpace(req.RelevanceNote),
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if err := seo_strategy.ValidateIntentEvidenceLink(link); err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}
	if err := u.repository.LinkEvidence(ctx, link); err != nil {
		return nil, mapIntentRepoError(err)
	}
	return u.Get(ctx, intent.ID)
}

func (u *SeoIntentHypothesisUsecase) UnlinkEvidence(
	ctx context.Context,
	req *dto.UnlinkSeoIntentEvidenceRequest,
) (*seo_strategy.SeoIntentHypothesisAggregate, error) {
	if req == nil || req.IntentHypothesisID == 0 || req.QueryEvidenceID == 0 {
		return nil, _errors.ReturnError(400, "intent evidence unlink request không hợp lệ")
	}
	intent, err := u.repository.GetByID(ctx, req.IntentHypothesisID)
	if err != nil {
		return nil, mapIntentRepoError(err)
	}
	if intent.Status == seo_strategy.IntentHypothesisStatusArchived {
		return nil, _errors.ReturnError(409, "không thể unlink evidence khỏi intent hypothesis đã archived")
	}
	if err := u.repository.UnlinkEvidence(ctx, intent.ID, req.QueryEvidenceID); err != nil {
		return nil, mapIntentRepoError(err)
	}
	return u.Get(ctx, intent.ID)
}

func (u *SeoIntentHypothesisUsecase) Activate(
	ctx context.Context,
	req *dto.SeoIntentHypothesisCommandRequest,
) (*seo_strategy.SeoIntentHypothesisAggregate, error) {
	if req == nil || req.ID == 0 || req.ExpectedVersion == 0 {
		return nil, _errors.ReturnError(400, "intent hypothesis command không hợp lệ")
	}
	aggregate, err := u.Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if aggregate.Intent.LockVersion != req.ExpectedVersion {
		return nil, _errors.ReturnError(409, "intent hypothesis đã được cập nhật; hãy tải lại dữ liệu")
	}
	if err := seo_strategy.ValidateIntentHypothesisStatusTransition(
		aggregate.Intent.Status,
		seo_strategy.IntentHypothesisStatusActive,
	); err != nil {
		return nil, _errors.ReturnError(409, err.Error())
	}
	if !aggregate.Readiness.Ready {
		codes := make([]string, 0, len(aggregate.Readiness.Blockers))
		for _, item := range aggregate.Readiness.Blockers {
			codes = append(codes, item.Code)
		}
		return nil, _errors.ReturnError(
			422,
			"intent hypothesis chưa sẵn sàng: "+strings.Join(codes, ", "),
		)
	}
	if err := u.repository.UpdateStatus(
		ctx,
		aggregate.Intent.ID,
		seo_strategy.IntentHypothesisStatusActive,
		req.ExpectedVersion,
	); err != nil {
		return nil, mapIntentRepoError(err)
	}
	return u.Get(ctx, aggregate.Intent.ID)
}

func (u *SeoIntentHypothesisUsecase) Pause(
	ctx context.Context,
	req *dto.SeoIntentHypothesisCommandRequest,
) (*seo_strategy.SeoIntentHypothesisAggregate, error) {
	return u.changeStatus(ctx, req, seo_strategy.IntentHypothesisStatusPaused)
}

func (u *SeoIntentHypothesisUsecase) Archive(
	ctx context.Context,
	req *dto.SeoIntentHypothesisCommandRequest,
) (*seo_strategy.SeoIntentHypothesisAggregate, error) {
	return u.changeStatus(ctx, req, seo_strategy.IntentHypothesisStatusArchived)
}

func (u *SeoIntentHypothesisUsecase) changeStatus(
	ctx context.Context,
	req *dto.SeoIntentHypothesisCommandRequest,
	target seo_strategy.IntentHypothesisStatus,
) (*seo_strategy.SeoIntentHypothesisAggregate, error) {
	if req == nil || req.ID == 0 || req.ExpectedVersion == 0 {
		return nil, _errors.ReturnError(400, "intent hypothesis command không hợp lệ")
	}
	current, err := u.repository.GetByID(ctx, req.ID)
	if err != nil {
		return nil, mapIntentRepoError(err)
	}
	if current.LockVersion != req.ExpectedVersion {
		return nil, _errors.ReturnError(409, "intent hypothesis đã được cập nhật; hãy tải lại dữ liệu")
	}
	if err := seo_strategy.ValidateIntentHypothesisStatusTransition(current.Status, target); err != nil {
		return nil, _errors.ReturnError(409, err.Error())
	}
	if err := u.repository.UpdateStatus(ctx, current.ID, target, req.ExpectedVersion); err != nil {
		return nil, mapIntentRepoError(err)
	}
	return u.Get(ctx, current.ID)
}

func mapIntentRepoError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, repo.ErrSeoIntentHypothesisNotFound):
		return _errors.ReturnError(404, "seo intent hypothesis không tồn tại")
	case errors.Is(err, repo.ErrSeoIntentHypothesisConflict):
		return _errors.ReturnError(409, "intent hypothesis đã được cập nhật; hãy tải lại dữ liệu")
	case errors.Is(err, repo.ErrSeoIntentHypothesisInvalid):
		return _errors.ReturnError(400, err.Error())
	case errors.Is(err, repo.ErrSeoIntentEvidenceAlreadyLinked):
		return _errors.ReturnError(409, "query evidence đã được liên kết với intent hypothesis")
	case errors.Is(err, repo.ErrSeoIntentEvidenceNotLinked):
		return _errors.ReturnError(404, "query evidence chưa được liên kết với intent hypothesis")
	case errors.Is(err, repo.ErrSeoIntentEvidenceMismatch):
		return _errors.ReturnError(409, "intent hypothesis và query evidence không cùng user need/objective")
	default:
		return err
	}
}

func mapIntentUserNeedError(err error) error {
	if errors.Is(err, repo.ErrSeoUserNeedNotFound) {
		return _errors.ReturnError(404, "seo user need không tồn tại")
	}
	return err
}
