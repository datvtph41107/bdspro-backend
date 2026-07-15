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

type seoQueryEvidenceUserNeedReader interface {
	GetByID(
		ctx context.Context,
		id uint64,
	) (*seo_strategy.SeoUserNeed, error)
}

type SeoQueryEvidenceUsecase struct {
	repository     repo.SeoQueryEvidenceRepo
	userNeedReader seoQueryEvidenceUserNeedReader
	now            func() time.Time
}

func NewSeoQueryEvidenceUsecase(
	repository repo.SeoQueryEvidenceRepo,
	userNeedRepository repo.SeoUserNeedRepo,
) *SeoQueryEvidenceUsecase {
	return &SeoQueryEvidenceUsecase{
		repository:     repository,
		userNeedReader: userNeedRepository,
		now:            time.Now,
	}
}

func (u *SeoQueryEvidenceUsecase) currentTime() time.Time {
	if u.now == nil {
		return time.Now().UTC()
	}
	return u.now().UTC()
}

func (u *SeoQueryEvidenceUsecase) Upsert(
	ctx context.Context,
	req *dto.UpsertSeoQueryEvidenceRequest,
) (*seo_strategy.SeoQueryEvidence, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không được để trống")
	}
	if req.UserNeedID == 0 {
		return nil, _errors.ReturnError(400, "userNeedId không hợp lệ")
	}

	need, err := u.userNeedReader.GetByID(ctx, req.UserNeedID)
	if err != nil {
		if errors.Is(err, repo.ErrSeoUserNeedNotFound) {
			return nil, _errors.ReturnError(404, "seo user need không tồn tại")
		}
		return nil, err
	}
	if need == nil {
		return nil, _errors.ReturnError(404, "seo user need không tồn tại")
	}
	if need.Status != seo_strategy.UserNeedStatusActive {
		return nil, _errors.ReturnError(
			409,
			"user need phải ở trạng thái active trước khi ghi nhận query evidence",
		)
	}

	sourceKind, err := seo_strategy.ParseQueryEvidenceSourceKind(req.SourceKind)
	if err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}
	observationType, err := seo_strategy.ParseQueryObservationType(req.ObservationType)
	if err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}
	deviceCategory, err := seo_strategy.ParseQueryDeviceCategory(req.DeviceCategory)
	if err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}

	confidence := seo_strategy.ObjectiveConfidence(strings.TrimSpace(req.Confidence))
	if confidence == "" {
		confidence = seo_strategy.ObjectiveConfidenceExploratory
	}
	if !seo_strategy.IsValidObjectiveConfidence(confidence) {
		return nil, _errors.ReturnError(400, "query evidence confidence không hợp lệ")
	}

	now := u.currentTime()
	evidence, err := seo_strategy.CreateSeoQueryEvidence(
		need.ObjectiveID,
		need.ID,
		req.RawQuery,
		req.SourceSystem,
		sourceKind,
		observationType,
		req.SourceRecordKey,
		now,
	)
	if err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}

	if locale := strings.TrimSpace(req.Locale); locale != "" {
		evidence.Locale = locale
	}
	if countryCode := strings.TrimSpace(req.CountryCode); countryCode != "" {
		evidence.CountryCode = strings.ToUpper(countryCode)
	}
	evidence.DeviceCategory = deviceCategory
	evidence.SourceReference = strings.TrimSpace(req.SourceReference)
	evidence.ObservedFrom = req.ObservedFrom
	evidence.ObservedTo = req.ObservedTo
	evidence.ValidUntil = req.ValidUntil
	evidence.Impressions = req.Impressions
	evidence.Clicks = req.Clicks
	evidence.CTR = req.CTR
	evidence.AvgPosition = req.AvgPosition
	evidence.SearchVolume = req.SearchVolume
	evidence.Confidence = confidence

	metadata := req.Metadata
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	evidence.Metadata = metadata

	// Locale có thể đã thay đổi sau factory nên hash phải tính lại.
	evidence.NormalizedQuery = seo_strategy.NormalizeSearchQuery(evidence.RawQuery)
	evidence.QueryHash = seo_strategy.BuildQueryEvidenceHash(
		evidence.NormalizedQuery,
		evidence.Locale,
	)
	evidence.UpdatedAt = now

	if err := seo_strategy.ValidateQueryEvidence(evidence); err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}
	if err := u.repository.Upsert(ctx, evidence); err != nil {
		return nil, mapQueryEvidenceRepoError(err)
	}
	return evidence, nil
}

func (u *SeoQueryEvidenceUsecase) Get(
	ctx context.Context,
	id uint64,
) (*seo_strategy.SeoQueryEvidence, error) {
	if id == 0 {
		return nil, _errors.ReturnError(400, "query evidence id không hợp lệ")
	}
	item, err := u.repository.GetByID(ctx, id)
	if err != nil {
		return nil, mapQueryEvidenceRepoError(err)
	}
	return item, nil
}

func (u *SeoQueryEvidenceUsecase) List(
	ctx context.Context,
	req *dto.SeoQueryEvidenceListRequest,
) ([]seo_strategy.SeoQueryEvidence, int64, error) {
	if req == nil {
		req = &dto.SeoQueryEvidenceListRequest{}
	}

	filter := repo.SeoQueryEvidenceFilter{
		ObjectiveID: req.ObjectiveID,
		UserNeedID:  req.UserNeedID,
		Keyword:     strings.TrimSpace(req.Keyword),
		SourceSystem: strings.TrimSpace(
			req.SourceSystem,
		),
		Page: req.Page,
		Size: req.Size,
	}

	if value := strings.TrimSpace(req.SourceKind); value != "" {
		parsed, err := seo_strategy.ParseQueryEvidenceSourceKind(value)
		if err != nil {
			return nil, 0, _errors.ReturnError(400, err.Error())
		}
		filter.SourceKind = parsed
	}
	if value := strings.TrimSpace(req.ObservationType); value != "" {
		parsed, err := seo_strategy.ParseQueryObservationType(value)
		if err != nil {
			return nil, 0, _errors.ReturnError(400, err.Error())
		}
		filter.ObservationType = parsed
	}
	if value := strings.TrimSpace(req.DeviceCategory); value != "" {
		parsed, err := seo_strategy.ParseQueryDeviceCategory(value)
		if err != nil {
			return nil, 0, _errors.ReturnError(400, err.Error())
		}
		filter.DeviceCategory = parsed
	}
	if value := strings.TrimSpace(req.Confidence); value != "" {
		confidence := seo_strategy.ObjectiveConfidence(value)
		if !seo_strategy.IsValidObjectiveConfidence(confidence) {
			return nil, 0, _errors.ReturnError(400, "confidence không hợp lệ")
		}
		filter.Confidence = confidence
	}
	if value := strings.TrimSpace(req.Status); value != "" {
		status, err := seo_strategy.ParseQueryEvidenceStatus(value)
		if err != nil {
			return nil, 0, _errors.ReturnError(400, err.Error())
		}
		filter.Status = status
	}
	if req.CurrentOnly {
		now := u.currentTime()
		filter.CurrentAt = &now
	}

	items, total, err := u.repository.List(ctx, filter)
	if err != nil {
		return nil, 0, mapQueryEvidenceRepoError(err)
	}
	return items, total, nil
}

func (u *SeoQueryEvidenceUsecase) Supersede(
	ctx context.Context,
	req *dto.SeoQueryEvidenceStatusCommandRequest,
) (*seo_strategy.SeoQueryEvidence, error) {
	return u.changeStatus(ctx, req, seo_strategy.QueryEvidenceStatusSuperseded)
}

func (u *SeoQueryEvidenceUsecase) Reject(
	ctx context.Context,
	req *dto.SeoQueryEvidenceStatusCommandRequest,
) (*seo_strategy.SeoQueryEvidence, error) {
	return u.changeStatus(ctx, req, seo_strategy.QueryEvidenceStatusRejected)
}

func (u *SeoQueryEvidenceUsecase) changeStatus(
	ctx context.Context,
	req *dto.SeoQueryEvidenceStatusCommandRequest,
	target seo_strategy.QueryEvidenceStatus,
) (*seo_strategy.SeoQueryEvidence, error) {
	if req == nil || req.ID == 0 {
		return nil, _errors.ReturnError(400, "query evidence command không hợp lệ")
	}
	if strings.TrimSpace(req.ExpectedStatus) == "" {
		return nil, _errors.ReturnError(400, "expectedStatus không được để trống")
	}

	expectedStatus, err := seo_strategy.ParseQueryEvidenceStatus(req.ExpectedStatus)
	if err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}
	current, err := u.repository.GetByID(ctx, req.ID)
	if err != nil {
		return nil, mapQueryEvidenceRepoError(err)
	}
	if current.Status != expectedStatus {
		return nil, _errors.ReturnError(
			409,
			"query evidence đã thay đổi trạng thái; hãy tải lại dữ liệu",
		)
	}
	if err := seo_strategy.ValidateQueryEvidenceStatusTransition(current.Status, target); err != nil {
		return nil, _errors.ReturnError(409, err.Error())
	}
	if err := u.repository.UpdateStatus(ctx, current.ID, expectedStatus, target); err != nil {
		return nil, mapQueryEvidenceRepoError(err)
	}
	return u.Get(ctx, current.ID)
}

func mapQueryEvidenceRepoError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, repo.ErrSeoQueryEvidenceNotFound):
		return _errors.ReturnError(404, "seo query evidence không tồn tại")
	case errors.Is(err, repo.ErrSeoQueryEvidenceConflict):
		return _errors.ReturnError(
			409,
			"query evidence xung đột với dữ liệu hiện tại; hãy tải lại dữ liệu",
		)
	case errors.Is(err, repo.ErrSeoQueryEvidenceInvalid):
		return _errors.ReturnError(400, err.Error())
	case errors.Is(err, repo.ErrSeoQueryEvidenceUserNeedMismatch):
		return _errors.ReturnError(409, "query evidence và user need không cùng objective")
	default:
		return err
	}
}
