package usecase

import (
	_errors "common/errors"
	"context"
	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/interface/provider"
	"crm/internal/repo"
	"encoding/json"
	"strings"
	"time"
)

type SeoDomainUsecase struct {
	seoDomainRepo   repo.SeoDomainRepo
	seoInternalRepo repo.SeoInternalLinkRepo
	seoRelativeRepo repo.SeoRelativeRepo
	transaction     provider.ITransaction
}

func NewSeoDomainUsecase(
	seoDomainRepo repo.SeoDomainRepo,
	seoInternalRepo repo.SeoInternalLinkRepo,
	seoRelativeRepo repo.SeoRelativeRepo,
	transaction provider.ITransaction,
) *SeoDomainUsecase {
	return &SeoDomainUsecase{
		seoDomainRepo:   seoDomainRepo,
		seoInternalRepo: seoInternalRepo,
		seoRelativeRepo: seoRelativeRepo,
		transaction:     transaction,
	}
}

func (u *SeoDomainUsecase) CreateSeoDomain(ctx context.Context, req *dto.SeoDomainCreateRequest) (*seo_domain.SeoDomain, error) {
	if strings.TrimSpace(req.OriginURL) == "" {
		return nil, _errors.ReturnError(400, "originUrl không được để trống")
	}
	if req.Classify != 0 && !enums.IsValidSeoClassify(req.Classify) {
		return nil, _errors.ReturnError(400, "classify không hợp lệ")
	}
	if req.Metadata != nil && strings.TrimSpace(*req.Metadata) != "" && !json.Valid([]byte(*req.Metadata)) {
		return nil, _errors.ReturnError(400, "metadata phải là JSON hợp lệ")
	}

	canonicalURL := strings.TrimSpace(req.OriginURL)
	if req.CanonicalURL != nil && strings.TrimSpace(*req.CanonicalURL) != "" {
		canonicalURL = strings.TrimSpace(*req.CanonicalURL)
	}

	exist, err := u.seoDomainRepo.ExistByCanonicalURL(ctx, canonicalURL, nil)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, _errors.ReturnError(400, "canonicalUrl đã tồn tại")
	}

	var siteMapLastedAt *time.Time
	if req.SiteMapLastedAt != nil && strings.TrimSpace(*req.SiteMapLastedAt) != "" {
		parsed, err := parseSeoTime(*req.SiteMapLastedAt)
		if err != nil {
			return nil, _errors.ReturnError(400, "siteMapLastedAt phải theo RFC3339 hoặc YYYY-MM-DD")
		}
		siteMapLastedAt = parsed
	}

	classify := enums.ESeoClassifyN
	if req.Classify != 0 {
		classify = enums.SeoClassifyFromUint32(req.Classify)
	}

	seoDomain := &seo_domain.SeoDomain{
		OriginURL:       strings.TrimSpace(req.OriginURL),
		CanonicalURL:    canonicalURL,
		Classify:        classify,
		SiteMapLastedAt: siteMapLastedAt,
		IsSiteMap:       req.IsSiteMap,
		IsIndex:         req.IsIndex,
		IsRobot:         req.IsRobot,
		Metadata:        req.Metadata,
		Note:            req.Note,
		RefType:         enums.ESEORefType(req.RefType),
		RefID:           req.RefID,
	}

	created, err := u.seoDomainRepo.Create(ctx, seoDomain)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (u *SeoDomainUsecase) GetSeoDomain(ctx context.Context, id uint64) (*seo_domain.SeoDomain, error) {
	seoDomain, err := u.seoDomainRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if seoDomain == nil {
		return nil, _errors.ReturnError(404, "seo domain không tồn tại")
	}
	return seoDomain, nil
}

func (u *SeoDomainUsecase) GetSeoDomainList(ctx context.Context, req *dto.SeoDomainListRequest) ([]seo_domain.SeoDomain, int64, error) {
	return u.seoDomainRepo.GetList(ctx, req)
}

func (u *SeoDomainUsecase) UpdateSeoDomain(ctx context.Context, id uint64, req *dto.SeoDomainUpdateRequest) (*seo_domain.SeoDomain, error) {
	existing, err := u.GetSeoDomain(ctx, id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}

	if req.OriginURL != nil {
		if strings.TrimSpace(*req.OriginURL) == "" {
			return nil, _errors.ReturnError(400, "originUrl không được để trống")
		}
		updates["origin_url"] = strings.TrimSpace(*req.OriginURL)
	}
	if req.CanonicalURL != nil {
		canonicalURL := strings.TrimSpace(*req.CanonicalURL)
		if canonicalURL == "" {
			canonicalURL = existing.OriginURL
		}
		exist, err := u.seoDomainRepo.ExistByCanonicalURL(ctx, canonicalURL, &id)
		if err != nil {
			return nil, err
		}
		if exist {
			return nil, _errors.ReturnError(400, "canonicalUrl đã tồn tại")
		}
		updates["canonical_url"] = canonicalURL
	}
	if req.Classify != nil {
		if !enums.IsValidSeoClassify(*req.Classify) {
			return nil, _errors.ReturnError(400, "classify không hợp lệ")
		}
		updates["classify"] = enums.SeoClassifyFromUint32(*req.Classify)
	}
	if req.SiteMapLastedAt != nil {
		if strings.TrimSpace(*req.SiteMapLastedAt) == "" {
			updates["site_map_lasted_at"] = nil
		} else {
			parsed, err := parseSeoTime(*req.SiteMapLastedAt)
			if err != nil {
				return nil, _errors.ReturnError(400, "siteMapLastedAt phải theo RFC3339 hoặc YYYY-MM-DD")
			}
			updates["site_map_lasted_at"] = parsed
		}
	}
	if req.IsSiteMap != nil {
		updates["is_site_map"] = *req.IsSiteMap
	}
	if req.IsIndex != nil {
		updates["is_index"] = *req.IsIndex
	}
	if req.IsRobot != nil {
		updates["is_robot"] = *req.IsRobot
	}
	if req.Metadata != nil {
		if strings.TrimSpace(*req.Metadata) != "" && !json.Valid([]byte(*req.Metadata)) {
			return nil, _errors.ReturnError(400, "metadata phải là JSON hợp lệ")
		}
		updates["metadata"] = req.Metadata
	}
	if req.Note != nil {
		updates["note"] = *req.Note
	}
	if req.RefType != nil {
		updates["ref_type"] = enums.ESEORefType(*req.RefType)
	}
	if req.RefID != nil {
		updates["ref_id"] = req.RefID
	}

	if len(updates) == 0 {
		return existing, nil
	}

	if err := u.seoDomainRepo.Update(ctx, id, updates); err != nil {
		return nil, err
	}
	return u.GetSeoDomain(ctx, id)
}

func (u *SeoDomainUsecase) DeleteSeoDomain(ctx context.Context, id uint64) error {
	if _, err := u.GetSeoDomain(ctx, id); err != nil {
		return err
	}
	return u.seoDomainRepo.Delete(ctx, id)
}

func (u *SeoDomainUsecase) CreateSeoInternalLink(ctx context.Context, req *dto.SeoInternalLinkCreateRequest) (*seo_domain.SeoInternalLink, error) {
	if req.ParentSeoID == 0 || req.ChildSeoID == 0 {
		return nil, _errors.ReturnError(400, "parentSeoId và childSeoId là bắt buộc")
	}
	if req.ParentSeoID == req.ChildSeoID {
		return nil, _errors.ReturnError(400, "parentSeoId không được trùng childSeoId")
	}
	if _, err := u.GetSeoDomain(ctx, req.ParentSeoID); err != nil {
		return nil, err
	}
	if _, err := u.GetSeoDomain(ctx, req.ChildSeoID); err != nil {
		return nil, err
	}
	exist, err := u.seoInternalRepo.Exist(ctx, req.ParentSeoID, req.ChildSeoID, nil)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, _errors.ReturnError(400, "internal link đã tồn tại")
	}

	link := &seo_domain.SeoInternalLink{
		Link:        req.Link,
		Title:       req.Title,
		ParentSeoID: req.ParentSeoID,
		ChildSeoID:  req.ChildSeoID,
	}
	return u.seoInternalRepo.Create(ctx, link)
}

func (u *SeoDomainUsecase) GetSeoInternalLinkList(ctx context.Context, req *dto.SeoInternalLinkListRequest) ([]seo_domain.SeoInternalLink, int64, error) {
	return u.seoInternalRepo.GetList(ctx, req)
}

func (u *SeoDomainUsecase) UpdateSeoInternalLink(ctx context.Context, id uint64, req *dto.SeoInternalLinkUpdateRequest) (*seo_domain.SeoInternalLink, error) {
	existing, err := u.seoInternalRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, _errors.ReturnError(404, "seo internal link không tồn tại")
	}

	parentID := existing.ParentSeoID
	childID := existing.ChildSeoID
	if req.ParentSeoID != nil {
		parentID = *req.ParentSeoID
	}
	if req.ChildSeoID != nil {
		childID = *req.ChildSeoID
	}
	if parentID == 0 || childID == 0 {
		return nil, _errors.ReturnError(400, "parentSeoId và childSeoId là bắt buộc")
	}
	if parentID == childID {
		return nil, _errors.ReturnError(400, "parentSeoId không được trùng childSeoId")
	}
	if _, err := u.GetSeoDomain(ctx, parentID); err != nil {
		return nil, err
	}
	if _, err := u.GetSeoDomain(ctx, childID); err != nil {
		return nil, err
	}
	exist, err := u.seoInternalRepo.Exist(ctx, parentID, childID, &id)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, _errors.ReturnError(400, "internal link đã tồn tại")
	}

	updates := map[string]interface{}{}
	if req.Link != nil {
		updates["link"] = *req.Link
	}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.ParentSeoID != nil {
		updates["parent_seo_id"] = *req.ParentSeoID
	}
	if req.ChildSeoID != nil {
		updates["child_seo_id"] = *req.ChildSeoID
	}
	if len(updates) == 0 {
		return existing, nil
	}
	if err := u.seoInternalRepo.Update(ctx, id, updates); err != nil {
		return nil, err
	}
	return u.seoInternalRepo.GetByID(ctx, id)
}

func (u *SeoDomainUsecase) DeleteSeoInternalLink(ctx context.Context, id uint64) error {
	existing, err := u.seoInternalRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return _errors.ReturnError(404, "seo internal link không tồn tại")
	}
	return u.seoInternalRepo.Delete(ctx, id)
}

func (u *SeoDomainUsecase) CreateSeoRelative(ctx context.Context, req *dto.SeoRelativeCreateRequest) (*seo_domain.SeoRelative, error) {
	if req.ParentSeoID == 0 || req.ChildSeoID == 0 {
		return nil, _errors.ReturnError(400, "parentSeoId và childSeoId là bắt buộc")
	}
	if req.ParentSeoID == req.ChildSeoID {
		return nil, _errors.ReturnError(400, "parentSeoId không được trùng childSeoId")
	}
	if _, err := u.GetSeoDomain(ctx, req.ParentSeoID); err != nil {
		return nil, err
	}
	if _, err := u.GetSeoDomain(ctx, req.ChildSeoID); err != nil {
		return nil, err
	}
	exist, err := u.seoRelativeRepo.Exist(ctx, req.ParentSeoID, req.ChildSeoID)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, _errors.ReturnError(400, "seo relative đã tồn tại")
	}

	relative := &seo_domain.SeoRelative{
		ParentSeoID: req.ParentSeoID,
		ChildSeoID:  req.ChildSeoID,
	}
	return u.seoRelativeRepo.Create(ctx, relative)
}

func (u *SeoDomainUsecase) GetSeoRelativeList(ctx context.Context, req *dto.SeoRelativeListRequest) ([]seo_domain.SeoRelative, int64, error) {
	return u.seoRelativeRepo.GetList(ctx, req)
}

func (u *SeoDomainUsecase) DeleteSeoRelative(ctx context.Context, id uint64) error {
	existing, err := u.seoRelativeRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return _errors.ReturnError(404, "seo relative không tồn tại")
	}
	return u.seoRelativeRepo.Delete(ctx, id)
}

func parseSeoTime(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	layouts := []string{time.RFC3339, "2006-01-02", "2006-01-02 15:04:05"}
	var lastErr error
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return &parsed, nil
		}
		lastErr = err
	}
	return nil, lastErr
}
