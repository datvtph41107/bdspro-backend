package usecases

import (
	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"strings"
	"user/enums"
	"user/internal"
	"user/internal/dto"
	"user/internal/interface/repo"
	models "user/internal/models"

	"gorm.io/gorm"
)

type TagUsecase struct {
	TagRepo repo.ITagRepo
}

func NewTagUsecase(tagRepo repo.ITagRepo) *TagUsecase {
	return &TagUsecase{
		TagRepo: tagRepo,
	}
}

func (u *TagUsecase) AdminCreateTag(ctx context.Context, req *dto.AdminTagRequest) (*models.TagEntity, error) {
	if err := validateAdminTagRequest(req); err != nil {
		return nil, err
	}

	tag := &models.TagEntity{
		Name:      strings.TrimSpace(req.Name),
		TagType:   req.TagType,
		IsDefault: true,
		IsActive:  req.IsActive,
		UserID:    nil,
	}

	if err := u.TagRepo.Create(ctx, tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (u *TagUsecase) AdminUpdateTag(ctx context.Context, req *dto.AdminTagRequest) (*models.TagEntity, error) {
	if req == nil || req.ID == 0 {
		return nil, _errors.ReturnError(service.TagPayloadRequired)
	}
	if err := validateAdminTagRequest(req); err != nil {
		return nil, err
	}

	existed, err := u.TagRepo.GetByID(ctx, req.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, _errors.ReturnError(service.TagNotFound)
		}
		return nil, err
	}

	existed.Name = strings.TrimSpace(req.Name)
	existed.TagType = req.TagType

	existed.IsDefault = req.IsDefault
	existed.IsActive = req.IsActive

	if err := u.TagRepo.Update(ctx, existed.ID, existed); err != nil {
		return nil, err
	}
	return existed, nil
}

func (u *TagUsecase) AdminDeleteTag(ctx context.Context, id uint64) error {
	if id == 0 {
		return _errors.ReturnError(service.TagPayloadRequired)
	}
	return u.TagRepo.Delete(ctx, id)
}

func (u *TagUsecase) ListTagsByType(ctx context.Context, req *dto.TagListRequest) ([]models.TagEntity, uint32, error) {
	if req == nil || !req.TagType.IsValid() {
		return nil, 0, _errors.ReturnError(service.TagTypeInvalid)
	}

	profileID := _utils.GetProfileIdWithContext(ctx)
	pagable := &_dto.Pagable{
		Page: req.Page,
		Size: req.Size,
		Sort: req.Sort,
	}
	return u.TagRepo.ListForUser(ctx, profileID, req.TagType, pagable, strings.TrimSpace(req.Text))
}

func (u *TagUsecase) AdminListTags(ctx context.Context, req *dto.AdminTagListRequest) ([]models.TagEntity, uint32, error) {
	var tagTypes []enums.TagType
	if req != nil && len(req.TagTypes) > 0 {
		tagTypes = req.TagTypes
	}
	pagable := &_dto.Pagable{
		Page: req.Page,
		Size: req.Size,
		Sort: req.Sort,
	}
	return u.TagRepo.ListAdmin(ctx, tagTypes, pagable, strings.TrimSpace(req.Text), req.IsDefault)
}

// GetProfileTagsByType lấy tags của profile cụ thể theo type (public API, không cần authentication)
func (u *TagUsecase) GetProfileTagsByType(ctx context.Context, profileID uint64, tagType enums.TagType) ([]models.TagEntity, uint32, error) {
	// if profileID == 0 {
	// 	return nil, 0, _errors.ReturnError(service.ProfileIDInvalidCaps)
	// }
	if !tagType.IsValid() {
		return nil, 0, _errors.ReturnError(service.TagTypeInvalid)
	}

	// Mặc định lấy tất cả, không phân trang cho public API
	pagable := &_dto.Pagable{
		Page: 0,
		Size: 1000,
	}
	return u.TagRepo.ListForUser(ctx, profileID, tagType, pagable, "")
}

func validateAdminTagRequest(req *dto.AdminTagRequest) error {
	if req == nil {
		return _errors.ReturnError(service.TagPayloadRequired)
	}
	if strings.TrimSpace(req.Name) == "" {
		return _errors.ReturnError(service.TagNameRequired)
	}
	if !req.TagType.IsValid() || req.TagType == 0 {
		return _errors.ReturnError(service.TagTypeInvalid)
	}
	return nil
}
